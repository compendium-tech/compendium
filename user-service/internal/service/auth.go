package service

import (
	"context"
	"fmt"
	"time"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	log "github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/compendium-tech/compendium/common/pkg/random"
	"github.com/compendium-tech/compendium/common/pkg/ratelimit"
	emailDelivery "github.com/compendium-tech/compendium/email-delivery-service/pkg/email"
	"github.com/compendium-tech/compendium/user-service/internal/domain"
	"github.com/compendium-tech/compendium/user-service/internal/email"
	appErr "github.com/compendium-tech/compendium/user-service/internal/error"
	"github.com/compendium-tech/compendium/user-service/internal/geoip"
	"github.com/compendium-tech/compendium/user-service/internal/hash"
	"github.com/compendium-tech/compendium/user-service/internal/model"
	"github.com/compendium-tech/compendium/user-service/internal/repository"
	"github.com/compendium-tech/compendium/user-service/internal/ua"
	"github.com/google/uuid"
)

// AuthService defines the interface for user authentication and session management.
// It offers a comprehensive suite of features, including user registration,
// login with multi-factor authentication (MFA), password reset workflows,
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
	SignUp(ctx context.Context, request domain.SignUpRequest) error
	SubmitMfaOtp(ctx context.Context, request domain.SubmitMfaOtpRequest) (*domain.SessionResponse, error)
	SignIn(ctx context.Context, request domain.SignInRequest) (*domain.SignInResponse, error)
	InitPasswordReset(ctx context.Context, request domain.InitPasswordResetRequest) error
	FinishPasswordReset(ctx context.Context, request domain.FinishPasswordResetRequest) error
	Refresh(ctx context.Context, request domain.RefreshTokenRequest) (*domain.SessionResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	GetSessionsForAuthenticatedUser(ctx context.Context, refreshToken string) ([]domain.Session, error)
	RemoveSessionByID(ctx context.Context, sessionID uuid.UUID, refreshToken string) error
}

type authService struct {
	authLockRepository      repository.AuthLockRepository
	trustedDeviceRepository repository.TrustedDeviceRepository
	userRepository          repository.UserRepository
	mfaRepository           repository.MfaRepository
	refreshTokenRepository  repository.RefreshTokenRepository
	emailSender             emailDelivery.EmailSender
	emailMessageBuilder     email.EmailMessageBuilder
	geoIP                   geoip.GeoIP
	userAgentParser         ua.UserAgentParser
	tokenManager            auth.TokenManager
	passwordHasher          hash.PasswordHasher
	ratelimiter             ratelimit.RateLimiter
}

func NewAuthService(
	authEmailLockRepository repository.AuthLockRepository,
	trustedDeviceRepository repository.TrustedDeviceRepository,
	userRepository repository.UserRepository,
	mfaRepository repository.MfaRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	emailSender emailDelivery.EmailSender,
	emailMessageBuilder email.EmailMessageBuilder,
	geoIP geoip.GeoIP,
	userAgentParser ua.UserAgentParser,
	tokenManager auth.TokenManager,
	passwordHasher hash.PasswordHasher,
	ratelimiter ratelimit.RateLimiter) AuthService {
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
		ratelimiter:             ratelimiter,
	}
}

func (s *authService) SignUp(ctx context.Context, request domain.SignUpRequest) (finalErr error) {
	logger := log.L(ctx).WithField("email", request.Email)
	logger.Info("Signing up user")

	lock, err := s.authLockRepository.ObtainLock(ctx, request.Email)
	if err != nil {
		return err
	}

	defer lock.ReleaseAndHandleErr(ctx, &finalErr)

	user, err := s.userRepository.FindUserByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	if user != nil {
		if user.IsEmailVerified {
			logger.Error("User tried to sign up but email is already taken and verified")
			return appErr.New(appErr.EmailTakenError)
		}

		if time.Now().UTC().Sub(user.CreatedAt) < 60*time.Second {
			logger.Error("User tried to repeat sign up in less than 60 seconds")
			return appErr.New(appErr.TooManyRequestsError)
		}
	}

	passwordHash, err := s.passwordHasher.HashPassword(request.Password)
	if err != nil {
		return err
	}

	if user != nil && !user.IsEmailVerified {
		logger.Info("Updating password hash and creation time for unverified user")

		err := s.userRepository.UpdatePasswordHashAndCreatedAt(ctx, user.ID, passwordHash, time.Now().UTC())
		if err != nil {
			return err
		}
	} else {
		logger.Info("Creating new user")

		isEmailTaken := false
		err := s.userRepository.CreateUser(ctx, model.User{
			ID:              uuid.New(),
			Name:            request.Name,
			Email:           request.Email,
			IsEmailVerified: false,
			IsAdmin:         false,
			PasswordHash:    passwordHash,
			CreatedAt:       time.Now().UTC(),
		}, &isEmailTaken)
		if err != nil {
			return err
		}

		if isEmailTaken {
			logger.Info("Failed to create new user model, email is already taken")

			return appErr.New(appErr.EmailTakenError)
		}
	}

	logger.Info("Setting MFA OTP")

	otp, err := random.NewRandomDigitCode(6)
	if err != nil {
		return err
	}

	err = s.mfaRepository.SetMfaOtpByEmail(ctx, request.Email, otp)
	if err != nil {
		return err
	}

	logger.Info("Sending sign-up MFA email")

	emailMessage, err := s.emailMessageBuilder.BuildSignUpMfaEmailMessage(request.Email, otp)
	if err != nil {
		return err
	}

	err = s.emailSender.SendMessage(emailMessage)
	if err != nil {
		return err
	}

	logger.Info("User sign-up initiated successfully")

	return nil
}

func (s *authService) SubmitMfaOtp(ctx context.Context, request domain.SubmitMfaOtpRequest) (_ *domain.SessionResponse, finalErr error) {
	logger := log.L(ctx).WithField("email", request.Email)
	logger.Info("Submitting MFA OTP")

	lock, err := s.authLockRepository.ObtainLock(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	defer lock.ReleaseAndHandleErr(ctx, &finalErr)

	user, err := s.userRepository.FindUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		logger.Error("MFA submission for non-existent user or MFA not requested")
		return nil, appErr.New(appErr.MfaNotRequestedError)
	}

	otp, err := s.mfaRepository.GetMfaOtpByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if otp == nil {
		logger.Error("MFA submission but no OTP found")
		return nil, appErr.New(appErr.MfaNotRequestedError)
	}

	if *otp != request.Otp {
		logger.Error("Invalid MFA OTP submitted")
		return nil, appErr.New(appErr.InvalidMfaOtpError)
	}

	if !user.IsEmailVerified {
		logger.Info("Marking email as verified for user")

		err = s.userRepository.UpdateIsEmailVerifiedByEmail(ctx, request.Email, true)
		if err != nil {
			return nil, err
		}
	}

	logger.Info("Removing MFA OTP")

	err = s.mfaRepository.RemoveMfaOtpByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	session, err := s.createSession(ctx, user.ID, request.UserAgent, request.IPAddress)
	if err != nil {
		return nil, err
	}

	logger.Info("MFA OTP submitted successfully and session created")

	return session, nil
}

func (s *authService) SignIn(ctx context.Context, request domain.SignInRequest) (_ *domain.SignInResponse, finalErr error) {
	logger := log.L(ctx).WithField("email", request.Email)
	logger.Info("Signing in user")

	lock, err := s.authLockRepository.ObtainLock(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	defer lock.ReleaseAndHandleErr(ctx, &finalErr)

	user, err := s.userRepository.FindUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		logger.Error("Sign-in with invalid credentials (user not found)")
		return nil, appErr.New(appErr.InvalidCredentialsError)
	}

	isPasswordValid, err := s.passwordHasher.IsPasswordHashValid(user.PasswordHash, request.Password)
	if err != nil {
		return nil, err
	}

	if !isPasswordValid {
		logger.Error("Sign-in with invalid credentials (password mismatch)")
		return nil, appErr.New(appErr.InvalidCredentialsError)
	}

	deviceIsKnown, err := s.trustedDeviceRepository.DeviceExists(ctx, user.ID, request.UserAgent, request.IPAddress)
	if err != nil {
		return nil, err
	}

	session := (*domain.SessionResponse)(nil)

	if !deviceIsKnown {
		logger.Info("Device is not known. Initiating MFA")

		otp, err := random.NewRandomDigitCode(6)
		if err != nil {
			return nil, err
		}

		err = s.mfaRepository.SetMfaOtpByEmail(ctx, request.Email, otp)
		if err != nil {
			return nil, err
		}

		emailMessage, err := s.emailMessageBuilder.BuildSignInMfaEmailMessage(request.Email, otp)
		if err != nil {
			return nil, err
		}

		err = s.emailSender.SendMessage(emailMessage)
		if err != nil {
			return nil, err
		}

		logger.Info("Sign-in successfully processed, MFA required")

		return &domain.SignInResponse{IsMfaRequired: true, Session: nil}, nil
	} else {
		logger.Info("Device is known. Creating session")
		session, err = s.createSession(ctx, user.ID, request.UserAgent, request.IPAddress)
		if err != nil {
			return nil, err
		}

		logger.Info("Sign-in successful, no MFA required")

		return &domain.SignInResponse{IsMfaRequired: false, Session: session}, nil
	}
}

func (s *authService) InitPasswordReset(ctx context.Context, request domain.InitPasswordResetRequest) (finalErr error) {
	logger := log.L(ctx).WithField("email", request.Email)
	logger.Info("Initiating password reset")

	lock, err := s.authLockRepository.ObtainLock(ctx, request.Email)
	if err != nil {
		return err
	}

	defer lock.ReleaseAndHandleErr(ctx, &finalErr)

	user, err := s.userRepository.FindUserByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	if user == nil || !user.IsEmailVerified {
		logger.Error("Password reset initiated for non-existent or unverified user")
		return appErr.New(appErr.UserNotFoundError)
	}

	logger.Info("Setting MFA OTP for password reset")

	otp, err := random.NewRandomDigitCode(6)
	if err != nil {
		return err
	}

	err = s.mfaRepository.SetMfaOtpByEmail(ctx, request.Email, otp)
	if err != nil {
		return err
	}

	logger.Info("Sending password reset MFA email")

	emailMessage, err := s.emailMessageBuilder.BuildPasswordResetMfaEmailMessage(request.Email, otp)
	if err != nil {
		return err
	}

	err = s.emailSender.SendMessage(emailMessage)
	if err != nil {
		return err
	}

	logger.Info("Password reset initiated successfully")

	return nil
}

func (s *authService) FinishPasswordReset(ctx context.Context, request domain.FinishPasswordResetRequest) (finalErr error) {
	logger := log.L(ctx).WithField("email", request.Email)
	logger.Info("Finishing password reset")

	lock, err := s.authLockRepository.ObtainLock(ctx, request.Email)
	if err != nil {
		return err
	}

	defer lock.ReleaseAndHandleErr(ctx, &finalErr)

	user, err := s.userRepository.FindUserByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	if user == nil || !user.IsEmailVerified {
		logger.Error("Password reset finish for non-existent or unverified user, or MFA not requested")
		return appErr.New(appErr.MfaNotRequestedError)
	}

	otp, err := s.mfaRepository.GetMfaOtpByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	if otp == nil {
		logger.Error("Password reset finish but no OTP found")
		return appErr.New(appErr.MfaNotRequestedError)
	}

	if *otp != request.Otp {
		logger.Error("Invalid MFA OTP submitted for password reset")
		return appErr.New(appErr.InvalidMfaOtpError)
	}

	logger.Info("Removing MFA OTP after successful password reset OTP verification")

	err = s.mfaRepository.RemoveMfaOtpByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	passwordHash, err := s.passwordHasher.HashPassword(request.Password)
	if err != nil {
		return err
	}

	logger.Info("Updating password hash for user")

	if err := s.userRepository.UpdatePasswordHash(ctx, user.ID, passwordHash); err != nil {
		return err
	}

	logger.Info("Password reset finished successfully")

	return nil
}

func (s *authService) Refresh(ctx context.Context, request domain.RefreshTokenRequest) (*domain.SessionResponse, error) {
	logger := log.L(ctx).WithField("refreshTokenPrefix", request.RefreshToken[:8])
	logger.Info("Refreshing session")

	token, isReused, err := s.refreshTokenRepository.GetRefreshToken(ctx, request.RefreshToken)
	if err != nil {
		return nil, err
	}

	if token == nil {
		logger.Error("Invalid refresh token was used (not found)")
		return nil, appErr.New(appErr.InvalidSessionError)
	}

	if isReused {
		logger.Warn("Refresh Token Reuse Detected! Forcing logout of all sessions for this user!")
		if err := s.refreshTokenRepository.RemoveAllRefreshTokensForUser(ctx, token.UserID); err != nil {
			logger.Errorf("Failed to remove all refresh tokens after reuse detection: %v", err)
			return nil, err
		}
		return nil, appErr.New(appErr.InvalidSessionError)
	}

	isRateLimited, err := s.ratelimiter.IsRateLimited(ctx, fmt.Sprintf("refresh:%s", token.Session.ID), time.Minute, 1)
	if err != nil {
		return nil, err
	}

	if isRateLimited {
		logger.Warn("Refresh token rate limit exceeded")
		return nil, appErr.New(appErr.TooManyRequestsError)
	}

	err = s.refreshTokenRepository.RemoveRefreshToken(ctx, request.RefreshToken, token.UserID)
	if err != nil {
		return nil, err
	}

	session, err := s.createSession(ctx, token.UserID, request.UserAgent, request.IPAddress)
	if err != nil {
		return nil, err
	}

	logger.Infof("Session refreshed successfully (new refresh token prefix: %s)", session.RefreshToken[:8])

	return session, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	logger := log.L(ctx).WithField("refreshTokenPrefix", refreshToken[:8])
	logger.Info("Invalidating session")

	token, isReused, err := s.refreshTokenRepository.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	if token == nil {
		logger.Error("Invalid refresh token was used for logout (not found)")
		return appErr.New(appErr.InvalidSessionError)
	}

	if isReused {
		logger.Warn("Refresh Token Reuse Detected during logout! Forcing logout of all sessions for this user!")
		err := s.refreshTokenRepository.RemoveAllRefreshTokensForUser(ctx, token.UserID)
		if err != nil {
			return err
		}

		return appErr.New(appErr.InvalidSessionError)
	}

	err = s.refreshTokenRepository.RemoveRefreshToken(ctx, refreshToken, token.UserID)
	if err != nil {
		return err
	}

	logger.Info("Session invalidated successfully")

	return nil
}

func (s *authService) GetSessionsForAuthenticatedUser(ctx context.Context, refreshToken string) ([]domain.Session, error) {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	refreshTokens, err := s.refreshTokenRepository.GetAllRefreshTokensForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	sessionsResponse := make([]domain.Session, len(refreshTokens))
	for i, token := range refreshTokens {
		location := ""
		location, err := s.geoIP.GetLocation(token.Session.IPAddress)
		if err != nil {
			log.L(ctx).Warnf("Failed to get location for session %s by ip %s: %v", token.Session.ID, token.Session.IPAddress, err)
		}

		deviceInfo := s.userAgentParser.ParseUserAgent(token.Session.UserAgent)

		sessionsResponse[i] = domain.Session{
			ID:        token.Session.ID,
			IsCurrent: token.Token == refreshToken,
			Name:      deviceInfo.Name,
			Os:        deviceInfo.Os,
			Device:    deviceInfo.Device,
			Location:  location,
			UserAgent: token.Session.UserAgent,
			IPAddress: token.Session.IPAddress,
			CreatedAt: token.Session.CreatedAt,
		}
	}

	log.L(ctx).Infof("Found %d sessions", len(sessionsResponse))

	return sessionsResponse, nil
}

func (s *authService) RemoveSessionByID(ctx context.Context, sessionID uuid.UUID, refreshToken string) error {
	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return err
	}
	logger := log.L(ctx).
		WithField("sessionId", sessionID.String()).
		WithField("refreshTokenPrefix", refreshToken[:8])
	logger.Info("Removing session for authenticated user")

	token, isReused, err := s.refreshTokenRepository.GetRefreshTokenBySessionID(ctx, sessionID)
	if err != nil {
		return err
	}

	if token == nil {
		logger.Error("Attempted to remove non-existent session")
		return appErr.New(appErr.SessionNotFoundError)
	}

	if token.UserID != userID {
		logger.Error("User attempted to remove session which does not belong to them")
		return appErr.New(appErr.SessionNotFoundError)
	}

	if isReused {
		logger.Warn("Refresh Token Reuse Detected during attempt to remove session! Forcing logout of all sessions for this user!")
		err := s.refreshTokenRepository.RemoveAllRefreshTokensForUser(ctx, token.UserID)
		if err != nil {
			return err
		}
		return appErr.New(appErr.InvalidSessionError)
	}

	err = s.refreshTokenRepository.RemoveRefreshToken(ctx, token.Token, userID)
	if err != nil {
		return err
	}

	logger.Info("Session removed successfully")

	return nil
}

func (s *authService) createSession(ctx context.Context, userID uuid.UUID, userAgent, ipAddress string) (*domain.SessionResponse, error) {
	logger := log.L(ctx).WithField("userId", userID.String())
	logger.Info("Creating new session for user")

	csrfToken, err := s.tokenManager.NewCsrfToken()
	if err != nil {
		return nil, err
	}

	accessTokenExpiresAt := time.Now().Add(20 * time.Minute)
	accessToken, err := s.tokenManager.NewAccessToken(userID, csrfToken, accessTokenExpiresAt)
	if err != nil {
		return nil, err
	}

	logger.Info("Creating new device entry if not already present for user")

	err = s.trustedDeviceRepository.ExistsOrCreateDevice(ctx, model.TrustedDevice{
		ID:        uuid.New(),
		UserID:    userID,
		UserAgent: userAgent,
		IPAddress: ipAddress,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}

	refreshTokenExpiresAt := time.Now().Add(20 * 24 * time.Hour)
	refreshToken := uuid.New().String()
	sessionID := uuid.New()

	logger.WithField("sessionId", sessionID.String()).Infof("Adding refresh token to repository")

	err = s.refreshTokenRepository.CreateRefreshToken(ctx, model.RefreshToken{
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
	if err != nil {
		return nil, err
	}

	logger.Infof("Session created successfully (refresh token prefix: %s)", refreshToken[:8])

	return &domain.SessionResponse{
		CsrfToken:             csrfToken,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}
