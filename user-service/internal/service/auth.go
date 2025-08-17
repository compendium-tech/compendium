package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/compendium-tech/compendium/common/pkg/random"
	"github.com/compendium-tech/compendium/common/pkg/ratelimit"

	"github.com/compendium-tech/compendium/user-service/internal/domain"
	"github.com/compendium-tech/compendium/user-service/internal/email"
	myerror "github.com/compendium-tech/compendium/user-service/internal/error"
	"github.com/compendium-tech/compendium/user-service/internal/geoip"
	"github.com/compendium-tech/compendium/user-service/internal/hash"
	"github.com/compendium-tech/compendium/user-service/internal/model"
	"github.com/compendium-tech/compendium/user-service/internal/repository"
	"github.com/compendium-tech/compendium/user-service/internal/ua"
)

// AuthService defines the interface for user authentication and session management.
// It offers a comprehensive suite of features, including user registration,
// login with multifactor authentication (MFA), password reset workflows,
// session refreshing, and session invalidation.
//
// # Authentication
//
// For operations requiring an authenticated user, an authentication middleware
// typically handles the authentication process, populating user-specific data
// (userID) into the context, as in most cases.
//
// However, for operations involving refresh tokens - specifically Logout, RemoveSessionByID,
// and Refresh - the refreshToken is explicitly provided. This allows the service
// to check if the caller's refresh token has already been revoked or is compromised.
// In GetSessionsForAuthenticatedUser refreshToken is used to identify
// the current session and is not checked against reuse.
type AuthService interface {
	SignUp(ctx context.Context, request domain.SignUpRequest)
	SubmitMfaOtp(ctx context.Context, request domain.SubmitMfaOtpRequest) domain.SessionResponse
	SignIn(ctx context.Context, request domain.SignInRequest) domain.SignInResponse
	InitPasswordReset(ctx context.Context, request domain.InitPasswordResetRequest)
	FinishPasswordReset(ctx context.Context, request domain.FinishPasswordResetRequest)
	Refresh(ctx context.Context, request domain.RefreshTokenRequest) domain.SessionResponse
	Logout(ctx context.Context, refreshToken string)
	GetSessionsForAuthenticatedUser(ctx context.Context, refreshToken string) []domain.Session
	RemoveSessionByID(ctx context.Context, sessionID uuid.UUID, refreshToken string)
}

type authService struct {
	authLockRepository      repository.AuthLockRepository
	trustedDeviceRepository repository.TrustedDeviceRepository
	userRepository          repository.UserRepository
	mfaRepository           repository.MfaRepository
	refreshTokenRepository  repository.RefreshTokenRepository
	emailSender             email.Sender
	emailMessageBuilder     email.MessageBuilder
	geoIP                   geoip.GeoIP
	userAgentParser         ua.UserAgentParser
	tokenManager            auth.TokenManager
	passwordHasher          hash.PasswordHasher
	rateLimiter             ratelimit.RateLimiter
}

func NewAuthService(
	authEmailLockRepository repository.AuthLockRepository,
	trustedDeviceRepository repository.TrustedDeviceRepository,
	userRepository repository.UserRepository,
	mfaRepository repository.MfaRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	emailSender email.Sender,
	emailMessageBuilder email.MessageBuilder,
	geoIP geoip.GeoIP,
	userAgentParser ua.UserAgentParser,
	tokenManager auth.TokenManager,
	passwordHasher hash.PasswordHasher,
	rateLimiter ratelimit.RateLimiter) AuthService {
	return &authService{
		authLockRepository:      authEmailLockRepository,
		trustedDeviceRepository: trustedDeviceRepository,
		userRepository:          userRepository,
		mfaRepository:           mfaRepository,
		refreshTokenRepository:  refreshTokenRepository,
		emailSender:             emailSender,
		emailMessageBuilder:     emailMessageBuilder,
		geoIP:                   geoIP,
		userAgentParser:         userAgentParser,
		tokenManager:            tokenManager,
		passwordHasher:          passwordHasher,
		rateLimiter:             rateLimiter,
	}
}

func (s *authService) SignUp(ctx context.Context, request domain.SignUpRequest) {
	logger := log.L(ctx).WithField("email", request.Email)
	logger.Info("Signing up user")

	lock := s.authLockRepository.ObtainLock(ctx, request.Email)
	defer lock.Release(ctx)

	user := s.userRepository.FindUserByEmail(ctx, request.Email)
	if user != nil {
		if user.IsEmailVerified {
			logger.Error("User tried to sign up but email is already taken and verified")
			myerror.New(myerror.EmailTakenError).Throw()
		}

		if time.Now().UTC().Sub(user.CreatedAt) < 60*time.Second {
			logger.Error("User tried to repeat sign up in less than 60 seconds")
			myerror.New(myerror.TooManyRequestsError).Throw()
		}
	}

	passwordHash := s.passwordHasher.HashPassword(request.Password)

	if user != nil && !user.IsEmailVerified {
		logger.Info("Updating password hash and creation time for unverified user")

		s.userRepository.UpdatePasswordHashAndCreatedAt(ctx, user.ID, passwordHash, time.Now().UTC())
	} else {
		logger.Info("Creating new user")

		isEmailTaken := false
		s.userRepository.CreateUser(ctx, model.User{
			ID:              uuid.New(),
			Name:            request.Name,
			Email:           request.Email,
			IsEmailVerified: false,
			IsAdmin:         false,
			PasswordHash:    passwordHash,
			CreatedAt:       time.Now().UTC(),
		}, &isEmailTaken)

		if isEmailTaken {
			logger.Info("Failed to create new user model, email is already taken")
			myerror.New(myerror.EmailTakenError).Throw()
		}
	}

	logger.Info("Sending MFA OTP")

	otp := random.NewRandomDigitCode(6)
	s.mfaRepository.SetMfaOtpByEmail(ctx, request.Email, otp)
	s.emailSender.SendMessage(
		s.emailMessageBuilder.SignUpEmail(request.Email, otp))

	logger.Info("User sign-up initiated successfully")
}

func (s *authService) SubmitMfaOtp(ctx context.Context, request domain.SubmitMfaOtpRequest) domain.SessionResponse {
	logger := log.L(ctx).WithField("email", request.Email)
	logger.Info("Submitting MFA OTP")

	lock := s.authLockRepository.ObtainLock(ctx, request.Email)
	defer lock.Release(ctx)

	user := s.userRepository.FindUserByEmail(ctx, request.Email)
	if user == nil {
		logger.Error("MFA submission for non-existent user or MFA not requested")
		myerror.New(myerror.MfaNotRequestedError).Throw()
	}

	otp := s.mfaRepository.GetMfaOtpByEmail(ctx, request.Email)
	if otp == nil {
		logger.Error("MFA submission but no OTP found")
		myerror.New(myerror.MfaNotRequestedError).Throw()
	}

	if *otp != request.Otp {
		logger.Error("Invalid MFA OTP submitted")
		myerror.New(myerror.InvalidMfaOtpError).Throw()
	}

	if !user.IsEmailVerified {
		logger.Info("Marking email as verified for user")
		s.userRepository.UpdateIsEmailVerifiedByEmail(ctx, request.Email, true)
	}

	logger.Info("Removing MFA OTP")

	s.mfaRepository.RemoveMfaOtpByEmail(ctx, request.Email)
	session := s.createSession(ctx, user.ID, request.UserAgent, request.IPAddress)

	logger.Info("MFA OTP submitted successfully and session created")

	return session
}

func (s *authService) SignIn(ctx context.Context, request domain.SignInRequest) domain.SignInResponse {
	logger := log.L(ctx).WithField("email", request.Email)
	logger.Info("Signing in user")

	lock := s.authLockRepository.ObtainLock(ctx, request.Email)
	defer lock.Release(ctx)

	user := s.userRepository.FindUserByEmail(ctx, request.Email)
	if user == nil {
		logger.Error("Sign-in with invalid credentials (user not found)")
		myerror.New(myerror.InvalidCredentialsError).Throw()
	}

	if !s.passwordHasher.IsPasswordHashValid(user.PasswordHash, request.Password) {
		logger.Error("Sign-in with invalid credentials (password mismatch)")
		myerror.New(myerror.InvalidCredentialsError).Throw()
	}

	if !s.trustedDeviceRepository.DeviceExists(ctx, user.ID, request.UserAgent, request.IPAddress) {
		logger.Info("Device is not known. Initiating MFA")

		otp := random.NewRandomDigitCode(6)
		s.mfaRepository.SetMfaOtpByEmail(ctx, request.Email, otp)
		s.emailSender.SendMessage(
			s.emailMessageBuilder.SignInEmail(request.Email, otp))

		logger.Info("Sign-in successfully processed, MFA required")

		return domain.SignInResponse{IsMfaRequired: true, Session: nil}
	} else {
		logger.Info("Device is known. Creating session")
		session := s.createSession(ctx, user.ID, request.UserAgent, request.IPAddress)

		logger.Info("Sign-in successful, no MFA required")

		return domain.SignInResponse{IsMfaRequired: false, Session: &session}
	}
}

func (s *authService) InitPasswordReset(ctx context.Context, request domain.InitPasswordResetRequest) {
	logger := log.L(ctx).WithField("email", request.Email)
	logger.Info("Initiating password reset")

	lock := s.authLockRepository.ObtainLock(ctx, request.Email)
	defer lock.Release(ctx)

	user := s.userRepository.FindUserByEmail(ctx, request.Email)
	if user == nil || !user.IsEmailVerified {
		logger.Error("Password reset initiated for non-existent or unverified user")
		myerror.New(myerror.UserNotFoundError).Throw()
	}

	logger.Info("Sending MFA OTP for password reset")

	otp := random.NewRandomDigitCode(6)
	s.mfaRepository.SetMfaOtpByEmail(ctx, request.Email, otp)
	s.emailSender.SendMessage(
		s.emailMessageBuilder.PasswordResetEmail(request.Email, otp))

	logger.Info("Password reset initiated successfully")
}

func (s *authService) FinishPasswordReset(ctx context.Context, request domain.FinishPasswordResetRequest) {
	logger := log.L(ctx).WithField("email", request.Email)
	logger.Info("Finishing password reset")

	lock := s.authLockRepository.ObtainLock(ctx, request.Email)
	defer lock.Release(ctx)

	user := s.userRepository.FindUserByEmail(ctx, request.Email)
	if user == nil || !user.IsEmailVerified {
		logger.Error("Password reset finish for non-existent or unverified user, or MFA not requested")
		myerror.New(myerror.MfaNotRequestedError).Throw()
	}

	otp := s.mfaRepository.GetMfaOtpByEmail(ctx, request.Email)
	if otp == nil {
		logger.Error("Password reset finish but no OTP found")
		myerror.New(myerror.MfaNotRequestedError).Throw()
	}

	if *otp != request.Otp {
		logger.Error("Invalid MFA OTP submitted for password reset")
		myerror.New(myerror.InvalidMfaOtpError).Throw()
	}

	logger.Info("Removing MFA OTP after successful password reset OTP verification")

	s.mfaRepository.RemoveMfaOtpByEmail(ctx, request.Email)
	passwordHash := s.passwordHasher.HashPassword(request.Password)

	logger.Info("Updating password hash for user")
	s.userRepository.UpdatePasswordHash(ctx, user.ID, passwordHash)

	logger.Info("Password reset finished successfully")
}

func (s *authService) Refresh(ctx context.Context, request domain.RefreshTokenRequest) domain.SessionResponse {
	logger := log.L(ctx).WithField("refreshTokenPrefix", request.RefreshToken[:8])
	logger.Info("Refreshing session")

	token, isReused := s.refreshTokenRepository.GetRefreshToken(ctx, request.RefreshToken)
	if token == nil {
		logger.Error("Invalid refresh token was used (not found)")
		myerror.New(myerror.InvalidSessionError).Throw()
	}

	if isReused {
		logger.Warn("Refresh Token Reuse Detected! Forcing logout of all sessions for this user!")

		s.refreshTokenRepository.RemoveAllRefreshTokensForUser(ctx, token.UserID)
		myerror.New(myerror.InvalidSessionError).Throw()
	}

	if s.rateLimiter.IsRateLimited(ctx, fmt.Sprintf("refresh:%s", token.Session.ID), time.Minute, 1) {
		logger.Warn("Refresh token rate limit exceeded")
		myerror.New(myerror.TooManyRequestsError).Throw()
	}

	s.refreshTokenRepository.RemoveRefreshToken(ctx, request.RefreshToken, token.UserID)
	session := s.createSession(ctx, token.UserID, request.UserAgent, request.IPAddress)

	logger.Infof("Session refreshed successfully (new refresh token prefix: %s)", session.RefreshToken[:8])

	return session
}

func (s *authService) Logout(ctx context.Context, refreshToken string) {
	logger := log.L(ctx).WithField("refreshTokenPrefix", refreshToken[:8])
	logger.Info("Invalidating session")

	token, isReused := s.refreshTokenRepository.GetRefreshToken(ctx, refreshToken)
	if token == nil {
		logger.Error("Invalid refresh token was used for logout (not found)")
		myerror.New(myerror.InvalidSessionError).Throw()
	}

	if isReused {
		logger.Warn("Refresh Token Reuse Detected during logout! Forcing logout of all sessions for this user!")
		s.refreshTokenRepository.RemoveAllRefreshTokensForUser(ctx, token.UserID)

		myerror.New(myerror.InvalidSessionError).Throw()
	}

	s.refreshTokenRepository.RemoveRefreshToken(ctx, refreshToken, token.UserID)

	logger.Info("Session invalidated successfully")
}

func (s *authService) GetSessionsForAuthenticatedUser(ctx context.Context, refreshToken string) []domain.Session {
	userID := auth.GetUserID(ctx)
	refreshTokens := s.refreshTokenRepository.GetAllRefreshTokensForUser(ctx, userID)

	sessionsResponse := make([]domain.Session, len(refreshTokens))
	for i, token := range refreshTokens {
		deviceInfo := s.userAgentParser.ParseUserAgent(token.Session.UserAgent)

		sessionsResponse[i] = domain.Session{
			ID:        token.Session.ID,
			IsCurrent: token.Token == refreshToken,
			Name:      deviceInfo.Name,
			Os:        deviceInfo.Os,
			Device:    deviceInfo.Device,
			Location:  s.geoIP.GetLocation(token.Session.IPAddress),
			UserAgent: token.Session.UserAgent,
			IPAddress: token.Session.IPAddress,
			CreatedAt: token.Session.CreatedAt,
		}
	}

	log.L(ctx).Infof("Found %d sessions", len(sessionsResponse))

	return sessionsResponse
}

func (s *authService) RemoveSessionByID(ctx context.Context, sessionID uuid.UUID, refreshToken string) {
	userID := auth.GetUserID(ctx)
	logger := log.L(ctx).
		WithField("sessionId", sessionID.String()).
		WithField("refreshTokenPrefix", refreshToken[:8])
	logger.Info("Removing session for authenticated user")

	token, isReused := s.refreshTokenRepository.GetRefreshTokenBySessionID(ctx, sessionID)
	if token == nil {
		logger.Error("Attempted to remove non-existent session")
		myerror.New(myerror.SessionNotFoundError).Throw()
	}

	if token.UserID != userID {
		logger.Error("User attempted to remove session which does not belong to them")
		myerror.New(myerror.SessionNotFoundError).Throw()
	}

	if isReused {
		logger.Warn("Refresh Token Reuse Detected during attempt to remove session! Forcing logout of all sessions for this user!")
		s.refreshTokenRepository.RemoveAllRefreshTokensForUser(ctx, token.UserID)

		myerror.New(myerror.InvalidSessionError).Throw()
	}

	s.refreshTokenRepository.RemoveRefreshToken(ctx, token.Token, userID)

	logger.Info("Session removed successfully")
}

func (s *authService) createSession(ctx context.Context, userID uuid.UUID, userAgent, ipAddress string) domain.SessionResponse {
	logger := log.L(ctx).WithField("userId", userID.String())
	logger.Info("Creating new session for user")

	csrfToken := s.tokenManager.NewCsrfToken()

	accessTokenExpiresAt := time.Now().Add(20 * time.Minute)
	accessToken := s.tokenManager.NewAccessToken(userID, csrfToken, accessTokenExpiresAt)

	logger.Info("Creating new device entry if not already present for user")

	s.trustedDeviceRepository.UpsertDevice(ctx, model.TrustedDevice{
		ID:        uuid.New(),
		UserID:    userID,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		CreatedAt: time.Now().UTC(),
	})

	refreshTokenExpiresAt := time.Now().Add(20 * 24 * time.Hour)
	refreshToken := uuid.New().String()
	sessionID := uuid.New()

	logger.WithField("sessionId", sessionID.String()).Infof("Adding refresh token to repository")

	s.refreshTokenRepository.CreateRefreshToken(ctx, model.RefreshToken{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: refreshTokenExpiresAt,
		Session: model.Session{
			ID:        sessionID,
			UserAgent: userAgent,
			IPAddress: ipAddress,
			CreatedAt: time.Now().UTC(),
		},
	})

	logger.Infof("Session created successfully (refresh token prefix: %s)", refreshToken[:8])

	return domain.SessionResponse{
		CsrfToken:             csrfToken,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}
}
