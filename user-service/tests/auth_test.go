package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/compendium-tech/compendium/user-service/internal/email"
	"github.com/compendium-tech/compendium/user-service/internal/model"
	"github.com/compendium-tech/compendium/user-service/internal/repository"
)

func (s *APITestSuite) Test_SignUp_WithValidCredentials_GeneratesAndStoresMfaCode() {
	var capturedMfaCode string

	r := s.Require()

	name, emailAddress, password := "John", "johndoe@test.com", "Qwerty123!!!"
	body := fmt.Sprintf(`{"name":"%s","email":"%s","password":"%s"}`, name, emailAddress, password)

	s.mockEmailMessageBuilder.On("SignUpEmail", emailAddress, mock.Anything).Return(email.Message{}, nil).Run(func(args mock.Arguments) {
		capturedMfaCode = args.Get(1).(string)
	}).Once()
	s.mockEmailSender.On("SendMessage", email.Message{}).Return(nil).Once()

	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	code := repository.NewRedisMfaRepository(s.GinAppDependencies.RedisClient).GetMfaOtpByEmail(s.ctx, emailAddress)
	r.NotNil(code)
	r.Equal(*code, capturedMfaCode)
}

func (s *APITestSuite) Test_SignUp_TooFrequently_ReturnsTooManySignupAttemptsError() {
	r := s.Require()

	name, emailAddress, password := "John", "johndoe@test.com", "Qwerty123!!!"
	body := fmt.Sprintf(`{"name":"%s","email":"%s","password":"%s"}`, name, emailAddress, password)

	s.mockEmailMessageBuilder.On("SignUpEmail", emailAddress, mock.Anything).Return(email.Message{}, nil).Once()
	s.mockEmailSender.On("SendMessage", email.Message{}).Return(nil).Once()

	req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	req, _ = http.NewRequest("POST", "/v1/users", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-type", "application/json")

	resp = httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	r.Equal(http.StatusTooManyRequests, resp.Result().StatusCode)
}

func (s *APITestSuite) Test_SignUp_WithInvalidBody_ReturnsError() {
	r := s.Require()
	bigString := strings.Repeat("a", 101)

	bodies := []string{
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

	for _, body := range bodies {
		req, _ := http.NewRequest("POST", "/v1/users", bytes.NewBuffer([]byte(body)))
		req.Header.Set("Content-type", "application/json")

		resp := httptest.NewRecorder()
		s.app.ServeHTTP(resp, req)

		r.Equal(http.StatusBadRequest, resp.Result().StatusCode)
	}
}

func (s *APITestSuite) Test_SubmitMfaOtp_WithValidOtp_CreatesNewSession() {
	r := s.Require()

	userID, email, otp := uuid.New(), "johndoe@test.com", "123456"
	ipAddress, userAgent := "1.0.0.0", "Test"

	repository.NewPgUserRepository(s.PgDB).CreateUser(s.ctx, model.User{
		ID:              userID,
		Name:            "John",
		Email:           email,
		IsEmailVerified: false,
		IsAdmin:         false,
		PasswordHash:    []byte{},
		CreatedAt:       time.Now().UTC(),
	}, nil)
	repository.NewRedisMfaRepository(s.GinAppDependencies.RedisClient).SetMfaOtpByEmail(s.ctx, email, otp)

	body := fmt.Sprintf(`{"email":"%s","otp":"%s"}`, email, otp)

	req, _ := http.NewRequest("POST", "/v1/sessions?flow=mfa", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-Real-IP", ipAddress)

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	respBody := resp.Result()

	var responseBody map[string]any
	err := json.NewDecoder(respBody.Body).Decode(&responseBody)
	r.NoError(err, "Failed to decode response body")

	fmt.Println(responseBody)

	_, ok := responseBody["accessTokenExpiresAt"]
	r.True(ok, "Response body should contain 'accessTokenExpiresAt'")

	_, ok = responseBody["refreshTokenExpiresAt"]
	r.True(ok, "Response body should contain 'refreshTokenExpiresAt'")

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	r.Nil(repository.NewRedisMfaRepository(s.GinAppDependencies.RedisClient).GetMfaOtpByEmail(s.ctx, email))
	r.True(repository.NewPgTrustedDeviceRepository(s.PgDB).DeviceExists(s.ctx, userID, userAgent, ipAddress))
}

func (s *APITestSuite) Test_SubmitMfaOtp_WithInvalidOtp_ReturnsUnauthorized() {
	r := s.Require()

	email, otp, otp2 := "johndoe@test.com", "123456", "234567"
	repository.NewPgUserRepository(s.PgDB).CreateUser(s.ctx, model.User{
		ID:              uuid.New(),
		Name:            "John",
		Email:           email,
		IsEmailVerified: false,
		IsAdmin:         false,
		PasswordHash:    []byte{},
		CreatedAt:       time.Now().UTC(),
	}, nil)
	repository.NewRedisMfaRepository(s.GinAppDependencies.RedisClient).SetMfaOtpByEmail(s.ctx, email, otp)

	body := fmt.Sprintf(`{"email":"%s","otp":"%s"}`, email, otp2)

	req, _ := http.NewRequest("POST", "/v1/sessions?flow=mfa", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	r.Equal(http.StatusUnauthorized, resp.Result().StatusCode)

	code := repository.NewRedisMfaRepository(s.GinAppDependencies.RedisClient).GetMfaOtpByEmail(s.ctx, email)
	r.NotEqual(code, nil)
	r.Equal(*code, otp)
}

func (s *APITestSuite) Test_SubmitMfaOtp_WithInvalidBody_ReturnsError() {
	r := s.Require()

	bodies := []string{
		`{"email":"","otp":"123456"}`,
		`{"email":"test","otp":"123456"}`,
		`{"email":"test@test","otp":"123456"}`,
		`{"email":"johndoe@test.com","otp":""}`,
		`{"email":"johndoe@test.com","otp":"1234"}`,
		`{"email":"johndoe@test.com","otp":"abcdef"}`,
	}

	for _, body := range bodies {
		req, _ := http.NewRequest("POST", "/v1/sessions?flow=mfa", bytes.NewBuffer([]byte(body)))
		req.Header.Set("Content-type", "application/json")

		resp := httptest.NewRecorder()
		s.app.ServeHTTP(resp, req)

		r.Equal(http.StatusBadRequest, resp.Result().StatusCode)
	}
}

func (s *APITestSuite) Test_SignIn_WithValidCredentialsOnKnownDevice_CreatesNewSession() {
	r := s.Require()

	userID, email, password := uuid.New(), "johndoe@test.com", "Qwerty12345!!"
	ipAddress, userAgent := "1.0.0.0", "Test"

	passwordHash := s.GinAppDependencies.PasswordHasher.HashPassword(password)
	repository.NewPgUserRepository(s.PgDB).CreateUser(s.ctx, model.User{
		ID:              userID,
		Name:            "John",
		Email:           email,
		IsEmailVerified: true,
		IsAdmin:         false,
		PasswordHash:    passwordHash,
		CreatedAt:       time.Now().UTC(),
	}, nil)
	repository.NewPgTrustedDeviceRepository(s.PgDB).UpsertDevice(s.ctx, model.TrustedDevice{
		ID:        uuid.New(),
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})

	s.mockEmailSender.AssertNotCalled(s.T(), "SendSignInMfaEmail")

	body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)

	req, _ := http.NewRequest("POST", "/v1/sessions?flow=password", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-Real-IP", ipAddress)

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	respBody := resp.Result()

	var responseBody map[string]any
	err := json.NewDecoder(respBody.Body).Decode(&responseBody)
	r.NoError(err, "Failed to decode response body")

	isMfaRequired, ok := responseBody["isMfaRequired"]
	r.True(ok, "Response body should contain 'isMfaRequired'")
	r.Equal(isMfaRequired, false)

	_, ok = responseBody["accessTokenExpiresAt"]
	log.Println(responseBody["accessTokenExpiresAt"])
	r.True(ok, "Response body should contain 'accessTokenExpiresAt'")

	_, ok = responseBody["refreshTokenExpiresAt"]
	r.True(ok, "Response body should contain 'refreshTokenExpiresAt'")

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	r.Nil(repository.NewRedisMfaRepository(s.GinAppDependencies.RedisClient).GetMfaOtpByEmail(s.ctx, email))
}

func (s *APITestSuite) Test_SignIn_WithValidCredentialsOnNewDevice_CreatesNewSession() {
	var capturedMfaCode string

	r := s.Require()

	userID, emailAddress, password := uuid.New(), "johndoe@test.com", "Qwerty12345!!"
	ipAddress, userAgent := "1.0.0.0", "Test"

	passwordHash := s.GinAppDependencies.PasswordHasher.HashPassword(password)

	repository.NewPgUserRepository(s.PgDB).CreateUser(s.ctx, model.User{
		ID:              userID,
		Name:            "John",
		Email:           emailAddress,
		IsEmailVerified: true,
		IsAdmin:         false,
		PasswordHash:    passwordHash,
		CreatedAt:       time.Now().UTC(),
	}, nil)

	s.mockEmailMessageBuilder.On("SignInEmail", emailAddress, mock.Anything).Return(email.Message{}, nil).Run(func(args mock.Arguments) {
		capturedMfaCode = args.Get(1).(string)
	}).Once()
	s.mockEmailSender.On("SendMessage", email.Message{}).Return(nil).Once()

	body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, emailAddress, password)

	req, _ := http.NewRequest("POST", "/v1/sessions?flow=password", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-Real-IP", ipAddress)

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	respBody := resp.Result()

	var responseBody map[string]any
	err := json.NewDecoder(respBody.Body).Decode(&responseBody)
	r.NoError(err, "Failed to decode response body")

	isMfaRequired, ok := responseBody["isMfaRequired"]
	r.True(ok, "Response body should contain 'isMfaRequired'")
	r.Equal(isMfaRequired, true)

	_, ok = responseBody["accessTokenExpiresAt"]
	r.False(ok, "Response body should not contain 'accessTokenExpiresAt'")

	_, ok = responseBody["refreshTokenExpiresAt"]
	r.False(ok, "Response body should not contain 'refreshTokenExpiresAt'")

	r.Equal(http.StatusAccepted, resp.Result().StatusCode)

	code := repository.NewRedisMfaRepository(s.GinAppDependencies.RedisClient).GetMfaOtpByEmail(s.ctx, emailAddress)
	r.NotNil(code)
	r.Equal(*code, capturedMfaCode)
}
