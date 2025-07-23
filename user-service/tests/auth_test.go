package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	emailDelivery "github.com/compendium-tech/compendium/email-delivery-service/pkg/email"
	"github.com/compendium-tech/compendium/user-service/internal/model"
	"github.com/compendium-tech/compendium/user-service/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func (s *APITestSuite) Test_SignUp_WithValidCredentials_GeneratesAndStoresMfaCode() {
	var capturedMfaCode string

	r := s.Require()

	name, email, password := "John", "johndoe@test.com", "Qwerty123!!!"
	body := fmt.Sprintf(`{"name":"%s","email":"%s","password":"%s"}`, name, email, password)

	s.mockEmailMessageBuilder.On("BuildSignUpMfaEmailMessage", email, mock.Anything).Return(emailDelivery.EmailMessage{}, nil).Run(func(args mock.Arguments) {
		capturedMfaCode = args.Get(1).(string)
	}).Once()
	s.mockEmailSender.On("SendMessage", emailDelivery.EmailMessage{}).Return(nil).Once()

	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	code, err := repository.NewRedisMfaRepository(s.Dependencies.RedisClient).GetMfaOtpByEmail(s.ctx, email)
	r.NoError(err, "Failed to fetch MFA OTP")
	r.NotNil(code)
	r.Equal(*code, capturedMfaCode)
}

func (s *APITestSuite) Test_SignUp_TooFrequently_ReturnsTooManySignupAttemptsError() {
	r := s.Require()

	name, email, password := "John", "johndoe@test.com", "Qwerty123!!!"
	body := fmt.Sprintf(`{"name":"%s","email":"%s","password":"%s"}`, name, email, password)

	s.mockEmailMessageBuilder.On("BuildSignUpMfaEmailMessage", email, mock.Anything).Return(emailDelivery.EmailMessage{}, nil).Once()
	s.mockEmailSender.On("SendMessage", emailDelivery.EmailMessage{}).Return(nil).Once()

	req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-type", "application/json")

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	req, _ = http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer([]byte(body)))
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
		req, _ := http.NewRequest("POST", "/api/v1/users", bytes.NewBuffer([]byte(body)))
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

	err := repository.NewPgUserRepository(s.PgDB).CreateUser(s.ctx, model.User{
		ID:              userID,
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

	body := fmt.Sprintf(`{"email":"%s","otp":"%s"}`, email, otp)

	req, _ := http.NewRequest("POST", "/api/v1/sessions?flow=mfa", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-Real-IP", ipAddress)

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	respBody := resp.Result()

	var responseBody map[string]any
	err = json.NewDecoder(respBody.Body).Decode(&responseBody)
	r.NoError(err, "Failed to decode response body")

	_, ok := responseBody["accessTokenExpiry"]
	r.True(ok, "Response body should contain 'accessTokenExpiry'")

	_, ok = responseBody["refreshTokenExpiry"]
	r.True(ok, "Response body should contain 'refreshTokenExpiry'")

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	code, err := repository.NewRedisMfaRepository(s.Dependencies.RedisClient).GetMfaOtpByEmail(s.ctx, email)
	r.NoError(err, "Failed to fetch MFA OTP")
	r.Nil(code)

	isDeviceKnownNow, err := repository.NewPgTrustedDeviceRepository(s.PgDB).DeviceExists(s.ctx, userID, userAgent, ipAddress)
	r.NoError(err, "Failed to check if device was created in db")
	r.True(isDeviceKnownNow)
}

func (s *APITestSuite) Test_SubmitMfaOtp_WithInvalidOtp_ReturnsUnauthorized() {
	r := s.Require()

	email, otp, otp2 := "johndoe@test.com", "123456", "234567"
	err := repository.NewPgUserRepository(s.PgDB).CreateUser(s.ctx, model.User{
		ID:              uuid.New(),
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

	body := fmt.Sprintf(`{"email":"%s","otp":"%s"}`, email, otp2)

	req, _ := http.NewRequest("POST", "/api/v1/sessions?flow=mfa", bytes.NewBuffer([]byte(body)))
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

	bodies := []string{
		`{"email":"","otp":"123456"}`,
		`{"email":"test","otp":"123456"}`,
		`{"email":"test@test","otp":"123456"}`,
		`{"email":"johndoe@test.com","otp":""}`,
		`{"email":"johndoe@test.com","otp":"1234"}`,
		`{"email":"johndoe@test.com","otp":"abcdef"}`,
	}

	for _, body := range bodies {
		req, _ := http.NewRequest("POST", "/api/v1/sessions?flow=mfa", bytes.NewBuffer([]byte(body)))
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

	passwordHash, err := s.Dependencies.PasswordHasher.HashPassword(password)
	r.NoError(err, "Failed to obtain password hash")

	err = repository.NewPgUserRepository(s.PgDB).CreateUser(s.ctx, model.User{
		ID:              userID,
		Name:            "John",
		Email:           email,
		IsEmailVerified: true,
		IsAdmin:         false,
		PasswordHash:    passwordHash,
		CreatedAt:       time.Now().UTC(),
	})
	r.NoError(err, "Failed to create new user")

	err = repository.NewPgTrustedDeviceRepository(s.PgDB).CreateDevice(s.ctx, model.TrustedDevice{
		ID:        uuid.New(),
		UserID:    userID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
	})
	r.NoError(err, "Failed to save new device in db")

	s.mockEmailSender.AssertNotCalled(s.T(), "SendSignInMfaEmail")

	body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)

	req, _ := http.NewRequest("POST", "/api/v1/sessions?flow=password", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-Real-IP", ipAddress)

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	respBody := resp.Result()

	var responseBody map[string]any
	err = json.NewDecoder(respBody.Body).Decode(&responseBody)
	r.NoError(err, "Failed to decode response body")

	isMfaRequired, ok := responseBody["isMfaRequired"]
	r.True(ok, "Response body should contain 'isMfaRequired'")
	r.Equal(isMfaRequired, false)

	_, ok = responseBody["accessTokenExpiry"]
	r.True(ok, "Response body should contain 'accessTokenExpiry'")

	_, ok = responseBody["refreshTokenExpiry"]
	r.True(ok, "Response body should contain 'refreshTokenExpiry'")

	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	code, err := repository.NewRedisMfaRepository(s.Dependencies.RedisClient).GetMfaOtpByEmail(s.ctx, email)
	r.NoError(err, "Failed to fetch MFA OTP")
	r.Nil(code)
}

func (s *APITestSuite) Test_SignIn_WithValidCredentialsOnNewDevice_CreatesNewSession() {
	var capturedMfaCode string

	r := s.Require()

	userID, email, password := uuid.New(), "johndoe@test.com", "Qwerty12345!!"
	ipAddress, userAgent := "1.0.0.0", "Test"

	passwordHash, err := s.Dependencies.PasswordHasher.HashPassword(password)
	r.NoError(err, "Failed to obtain password hash")

	err = repository.NewPgUserRepository(s.PgDB).CreateUser(s.ctx, model.User{
		ID:              userID,
		Name:            "John",
		Email:           email,
		IsEmailVerified: true,
		IsAdmin:         false,
		PasswordHash:    passwordHash,
		CreatedAt:       time.Now().UTC(),
	})
	r.NoError(err, "Failed to create new user")

	s.mockEmailMessageBuilder.On("BuildSignInMfaEmailMessage", email, mock.Anything).Return(emailDelivery.EmailMessage{}, nil).Run(func(args mock.Arguments) {
		capturedMfaCode = args.Get(1).(string)
	}).Once()
	s.mockEmailSender.On("SendMessage", emailDelivery.EmailMessage{}).Return(nil).Once()

	body := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)

	req, _ := http.NewRequest("POST", "/api/v1/sessions?flow=password", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("X-Real-IP", ipAddress)

	resp := httptest.NewRecorder()
	s.app.ServeHTTP(resp, req)

	respBody := resp.Result()

	var responseBody map[string]any
	err = json.NewDecoder(respBody.Body).Decode(&responseBody)
	r.NoError(err, "Failed to decode response body")

	isMfaRequired, ok := responseBody["isMfaRequired"]
	r.True(ok, "Response body should contain 'isMfaRequired'")
	r.Equal(isMfaRequired, true)

	_, ok = responseBody["accessTokenExpiry"]
	r.False(ok, "Response body should not contain 'accessTokenExpiry'")

	_, ok = responseBody["refreshTokenExpiry"]
	r.False(ok, "Response body should not contain 'refreshTokenExpiry'")

	r.Equal(http.StatusAccepted, resp.Result().StatusCode)

	code, err := repository.NewRedisMfaRepository(s.Dependencies.RedisClient).GetMfaOtpByEmail(s.ctx, email)
	r.NoError(err, "Failed to fetch MFA OTP")
	r.NotNil(code)
	r.Equal(*code, capturedMfaCode)
}
