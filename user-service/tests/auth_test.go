package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/seacite-tech/compendium/user-service/internal/model"
	"github.com/seacite-tech/compendium/user-service/internal/repository"
	"github.com/stretchr/testify/mock"
)

func (s *APITestSuite) Test_SignUp_WithValidCredentials_GeneratesAndStoresMfaCode() {
	var capturedMfaCode string

	r := s.Require()

	name, email, password := "John", "johndoe@test.com", "Qwerty123!!!"
	input := fmt.Sprintf(`{"name":"%s","email":"%s","password":"%s"}`, name, email, password)

	s.mockEmailSender.On("SendSignUpMfaEmail", email, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		capturedMfaCode = args.Get(1).(string)
	})

	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer([]byte(input)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	code, err := repository.NewRedisMfaRepository(s.Dependencies.RedisClient).GetMfaOtpByEmail(s.ctx, email)
	r.NoError(err, "Failed to fetch MFA OTP")
	r.NotEqual(code, nil)
	r.Equal(*code, capturedMfaCode)
}

func (s *APITestSuite) Test_SignUp_TooFrequently_ReturnsTooManySignupAttemptsError() {
	r := s.Require()

	name, email, password := "John", "johndoe@test.com", "Qwerty123!!!"
	input := fmt.Sprintf(`{"name":"%s","email":"%s","password":"%s"}`, name, email, password)

	s.mockEmailSender.On("SendSignUpMfaEmail", email, mock.Anything).Return(nil)

	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer([]byte(input)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	req, _ = http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer([]byte(input)))
	req.Header.Set("Content-type", "application/json")

	resp = httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	r.Equal(http.StatusTooManyRequests, resp.Result().StatusCode)
}

func (s *APITestSuite) Test_SignUp_WithInvalidBody_ReturnsError() {
	r := s.Require()
	bigString := strings.Repeat("a", 101)

	inputs := []string{
		`{"name":"","email":"johndoe@test.com","password":"Qwerty123!!!"}`,
		`{"name":"John","email":"","password":"Qwerty123!!!"}`,
		`{"name":"John","email":"test","password":"Qwerty123!!!"}`,
		`{"name":"John","email":"test@test","password":"Qwerty123!!!"}`,
		`{"name":"John","email":"johndoe@test.com","password":""}`,
		`{"name":"John","email":"johndoe@test.com","password":"ABCD"}`,
		`{"name":"John","email":"johndoe@test.com","password":"ABCDEF"}`,
		`{"name":"John","email":"johndoe@test.com","password":"ABCDEF!!"}`,
		`{"name":"John","email":"johndoe@test.com","password":"abcdef"}`,
		`{"name":"John","email":"johndoe@test.com","password":"abcdef!!"}`,
		`{"name":"John","email":"johndoe@test.com","password":"abcdef12345"}`,
		`{"name":"John","email":"johndoe@test.com","password":"ABCdef12345"}`,
		fmt.Sprintf(`{"name":"%s","email":"johndoe@test.com","password":"Qwerty12345!!"}`, bigString),
		fmt.Sprintf(`{"name":"John","email":"johndoe@test.com","password":"%s"}`, bigString),
	}

	for _, input := range inputs {
		req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer([]byte(input)))
		req.Header.Set("Content-type", "application/json")

		resp := httptest.NewRecorder()
		s.app.ServeHTTP(resp, req)

		r.Equal(http.StatusBadRequest, resp.Result().StatusCode)
	}
}

func (s *APITestSuite) Test_SubmitMfaOtp_WithValidOtp_CreatesNewSession() {
	r := s.Require()

	email, otp := "johndoe@test.com", "123456"
	err := repository.NewPgUserRepository(s.PgDb).CreateUser(s.ctx, model.User{
		Id:              uuid.New(),
		Name:            "John",
		Email:           email,
		IsEmailVerified: false,
		IsAdmin:         false,
		PasswordHash:    []byte{},
		CreatedAt:       time.Now().UTC(),
	})
	r.NoError(err, "Failed to create new user")

	err = repository.NewRedisMfaRepository(s.Dependencies.RedisClient).SetMfaOtpByEmail(s.ctx, email, otp)
	r.NoError(err, "Failed to set MFA OTP")

	input := fmt.Sprintf(`{"email":"%s","otp":"%s"}`, email, otp)

	req, _ := http.NewRequest("POST", "/api/v1/sessions?flow=mfa", bytes.NewBuffer([]byte(input)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)
	body := resp.Result()
	var responseBody map[string]interface{}
	err = json.NewDecoder(body.Body).Decode(&responseBody)
	r.NoError(err, "Failed to decode response body")

	_, ok := responseBody["accessTokenExpiry"]
	r.True(ok, "Response body should contain 'accessTokenExpiry'")

	_, ok = responseBody["refreshTokenExpiry"]
	r.True(ok, "Response body should contain 'refreshTokenExpiry'")

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	code, err := repository.NewRedisMfaRepository(s.Dependencies.RedisClient).GetMfaOtpByEmail(s.ctx, email)
	r.NoError(err, "Failed to fetch MFA OTP")
	r.Nil(code)
}

func (s *APITestSuite) Test_SubmitMfaOtp_WithInvalidOtp_ReturnsUnauthorized() {
	r := s.Require()

	email, otp, otp2 := "johndoe@test.com", "123456", "234567"
	err := repository.NewPgUserRepository(s.PgDb).CreateUser(s.ctx, model.User{
		Id:              uuid.New(),
		Name:            "John",
		Email:           email,
		IsEmailVerified: false,
		IsAdmin:         false,
		PasswordHash:    []byte{},
		CreatedAt:       time.Now().UTC(),
	})
	r.NoError(err, "Failed to create new user")

	err = repository.NewRedisMfaRepository(s.Dependencies.RedisClient).SetMfaOtpByEmail(s.ctx, email, otp)
	r.NoError(err, "Failed to set MFA OTP")

	input := fmt.Sprintf(`{"email":"%s","otp":"%s"}`, email, otp2)

	req, _ := http.NewRequest("POST", "/api/v1/sessions?flow=mfa", bytes.NewBuffer([]byte(input)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	r.Equal(http.StatusUnauthorized, resp.Result().StatusCode)

	code, err := repository.NewRedisMfaRepository(s.Dependencies.RedisClient).GetMfaOtpByEmail(s.ctx, email)
	r.NoError(err, "Failed to fetch MFA OTP")
	r.NotEqual(code, nil)
	r.Equal(*code, otp)
}

func (s *APITestSuite) Test_SubmitMfaOtp_WithInvalidBody_ReturnsError() {
	r := s.Require()

	inputs := []string{
		`{"email":"","otp":"123456"}`,
		`{"email":"test","otp":"123456"}`,
		`{"email":"test@test","otp":"123456"}`,
		`{"email":"johndoe@test.com","otp":""}`,
		`{"email":"johndoe@test.com","otp":"1234"}`,
		`{"email":"johndoe@test.com","otp":"abcdef"}`,
	}

	for _, input := range inputs {
		req, _ := http.NewRequest("POST", "/api/v1/sessions?flow=mfa", bytes.NewBuffer([]byte(input)))
		req.Header.Set("Content-type", "application/json")

		resp := httptest.NewRecorder()
		s.app.ServeHTTP(resp, req)

		r.Equal(http.StatusBadRequest, resp.Result().StatusCode)
	}
}
