package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	log "github.com/compendium-tech/compendium/common/pkg/log"
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

type AuthService interface {
	SignUp(ctx context.Context, request domain.SignUpRequest) error
	SubmitMfaOtp(ctx context.Context, request domain.SubmitMfaOtpRequest) (*domain.SessionResponse, error)
	SignIn(ctx context.Context, request domain.SignInRequest) (*domain.SignInResponse, error)
	InitPasswordReset(ctx context.Context, request domain.InitPasswordResetRequest) error
	FinishPasswordReset(ctx context.Context, request domain.FinishPasswordResetRequest) error
	Refresh(ctx context.Context, request domain.RefreshTokenRequest) (*domain.SessionResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	GetSessionsForAuthenticatedUser(ctx context.Context) ([]domain.Session, error)
	RemoveSessionByID(ctx context.Context, sessionID uuid.UUID, currentRefreshToken string) error
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
	log.L(ctx).Infof("Signing up user with email: %s", request.Email)

	lock, err := s.authLockRepository.ObtainLock(ctx, request.Email)
	if err != nil {
		return err
	}

	defer lock.ReleaseAndHandleErr(ctx, &finalErr)

	user, err := s.userRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	if user != nil {
		if user.IsEmailVerified {
			log.L(ctx).Errorf("User with email %s tried to sign up but email is already taken and verified", request.Email)

			return appErr.Errorf(appErr.EmailTakenError, "Email is already taken")
		}

		if time.Now().UTC().Sub(user.CreatedAt) < 60*time.Second {
			log.L(ctx).Errorf("User with email %s tried to repeat sign up in less than 60 seconds", request.Email)

			return appErr.Errorf(appErr.TooManyRequestsError, "Too many sign up attempts")
		}
	}

	passwordHash, err := s.passwordHasher.HashPassword(request.Password)
	if err != nil {
		return err
	}

	if user != nil && !user.IsEmailVerified {
		log.L(ctx).Infof("Updating password hash and creation time for unverified user: %s", request.Email)

		err := s.userRepository.UpdatePasswordHashAndCreatedAt(ctx, user.ID, passwordHash, time.Now().UTC())
		if err != nil {
			return err
		}
	} else {
		log.L(ctx).Infof("Creating new user with email: %s", request.Email)

		err := s.userRepository.CreateUser(ctx, model.User{
			ID:              uuid.New(),
			Name:            request.Name,
			Email:           request.Email,
			IsEmailVerified: false,
			IsAdmin:         false,
			PasswordHash:    passwordHash,
			CreatedAt:       time.Now().UTC(),
		})
		if err != nil {
			return err
		}
	}

	log.L(ctx).Infof("Setting MFA OTP for email: %s", request.Email)

	otp := newSixDigitOtp()
	err = s.mfaRepository.SetMfaOtpByEmail(ctx, request.Email, otp)
	if err != nil {
		return err
	}

	log.L(ctx).Infof("Sending sign-up MFA email to: %s", request.Email)

	emailMessage, err := s.emailMessageBuilder.BuildSignUpMfaEmailMessage(request.Email, otp)
	if err != nil {
		return err
	}

	err = s.emailSender.SendMessage(emailMessage)
	if err != nil {
		return err
	}

	log.L(ctx).Info("User sign-up initiated successfully")

	return nil
}

func (s *authService) SubmitMfaOtp(ctx context.Context, request domain.SubmitMfaOtpRequest) (_ *domain.SessionResponse, finalErr error) {
	log.L(ctx).Infof("Submitting MFA OTP for email: %s", request.Email)

	lock, err := s.authLockRepository.ObtainLock(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	defer lock.ReleaseAndHandleErr(ctx, &finalErr)

	user, err := s.userRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		log.L(ctx).Errorf("MFA submission for non-existent user or MFA not requested for email: %s", request.Email)

		return nil, appErr.Errorf(appErr.MfaNotRequestedError, "2FA was not requested")
	}

	otp, err := s.mfaRepository.GetMfaOtpByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if otp == nil {
		log.L(ctx).Errorf("MFA submission but no OTP found for email: %s", request.Email)

		return nil, appErr.Errorf(appErr.MfaNotRequestedError, "2FA was not requested")
	}

	if *otp != request.Otp {
		log.L(ctx).Errorf("Invalid MFA OTP submitted for email: %s", request.Email)

		return nil, appErr.Errorf(appErr.InvalidMfaOtpError, "Invalid 2FA otp")
	}

	if !user.IsEmailVerified {
		log.L(ctx).Infof("Marking email as verified for user: %s", request.Email)

		err = s.userRepository.UpdateIsEmailVerifiedByEmail(ctx, request.Email, true)
		if err != nil {
			return nil, err
		}
	}

	log.L(ctx).Infof("Removing MFA OTP for email: %s", request.Email)

	err = s.mfaRepository.RemoveMfaOtpByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	session, err := s.createSession(ctx, user.ID, request.UserAgent, request.IPAddress)
	if err != nil {
		return nil, err
	}

	log.L(ctx).Info("MFA OTP submitted successfully and session created")

	return session, nil
}

func (s *authService) SignIn(ctx context.Context, request domain.SignInRequest) (_ *domain.SignInResponse, finalErr error) {
	log.L(ctx).Infof("Signing in user with email: %s", request.Email)

	lock, err := s.authLockRepository.ObtainLock(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	defer lock.ReleaseAndHandleErr(ctx, &finalErr)

	user, err := s.userRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		log.L(ctx).Errorf("Sign-in with invalid credentials for email: %s (user not found)", request.Email)

		return nil, appErr.Errorf(appErr.InvalidCredentialsError, "Invalid credentials")
	}

	isPasswordValid, err := s.passwordHasher.IsPasswordHashValid(user.PasswordHash, request.Password)
	if err != nil {
		return nil, err
	}

	if !isPasswordValid {
		log.L(ctx).Errorf("Sign-in with invalid credentials for email: %s (password mismatch)", request.Email)

		return nil, appErr.Errorf(appErr.InvalidCredentialsError, "Invalid credentials")
	}

	// Device check for MFA requirement
	deviceIsKnown, err := s.trustedDeviceRepository.DeviceExists(ctx, user.ID, request.UserAgent, request.IPAddress)
	if err != nil {
		return nil, err
	}

	session := (*domain.SessionResponse)(nil)

	if !deviceIsKnown {
		log.L(ctx).Infof("Device is not known for user %s. Initiating MFA", user.ID)

		otp := newSixDigitOtp()
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

		log.L(ctx).Info("Sign-in is successfully processed, MFA required")

		return &domain.SignInResponse{IsMfaRequired: true, Session: nil}, nil
	} else {
		log.L(ctx).Infof("Device is known for user %s. Creating session", user.ID)
		session, err = s.createSession(ctx, user.ID, request.UserAgent, request.IPAddress)
		if err != nil {
			return nil, err
		}

		log.L(ctx).Info("Sign-in successful, no MFA required")

		return &domain.SignInResponse{IsMfaRequired: false, Session: session}, nil
	}
}

func (s *authService) InitPasswordReset(ctx context.Context, request domain.InitPasswordResetRequest) (finalErr error) {
	log.L(ctx).Infof("Initiating password reset for email: %s", request.Email)

	lock, err := s.authLockRepository.ObtainLock(ctx, request.Email)
	if err != nil {
		return err
	}

	defer lock.ReleaseAndHandleErr(ctx, &finalErr)

	user, err := s.userRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	if user == nil || !user.IsEmailVerified {
		log.L(ctx).Errorf("Password reset initiated for non-existent or unverified user: %s", request.Email)
		return appErr.Errorf(appErr.UserNotFoundError, "User with given email address doesn't exist")
	}

	log.L(ctx).Infof("Setting MFA OTP for password reset for email: %s", request.Email)

	otp := newSixDigitOtp()
	err = s.mfaRepository.SetMfaOtpByEmail(ctx, request.Email, otp)
	if err != nil {
		return err
	}

	log.L(ctx).Infof("Sending password reset MFA email to: %s", request.Email)

	emailMessage, err := s.emailMessageBuilder.BuildPasswordResetMfaEmailMessage(request.Email, otp)
	if err != nil {
		return err
	}

	err = s.emailSender.SendMessage(emailMessage)
	if err != nil {
		return err
	}

	log.L(ctx).Info("Password reset initiated successfully")

	return nil
}

func (s *authService) FinishPasswordReset(ctx context.Context, request domain.FinishPasswordResetRequest) (finalErr error) {
	log.L(ctx).Infof("Finishing password reset for email: %s", request.Email)

	lock, err := s.authLockRepository.ObtainLock(ctx, request.Email)
	if err != nil {
		return err
	}

	defer lock.ReleaseAndHandleErr(ctx, &finalErr)

	user, err := s.userRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	if user == nil || !user.IsEmailVerified {
		log.L(ctx).Errorf("Password reset finish for non-existent or unverified user, or MFA not requested for email: %s", request.Email)

		return appErr.Errorf(appErr.MfaNotRequestedError, "2FA was not requested")
	}

	otp, err := s.mfaRepository.GetMfaOtpByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	if otp == nil {
		log.L(ctx).Errorf("Password reset finish but no OTP found for email: %s", request.Email)

		return appErr.Errorf(appErr.MfaNotRequestedError, "2FA was not requested")
	}

	if *otp != request.Otp {
		log.L(ctx).Errorf("Invalid MFA OTP submitted for password reset for email: %s", request.Email)

		return appErr.Errorf(appErr.InvalidMfaOtpError, "Invalid 2FA otp")
	}

	log.L(ctx).Infof("Removing MFA OTP for email %s after successful password reset OTP verification", request.Email)

	err = s.mfaRepository.RemoveMfaOtpByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	passwordHash, err := s.passwordHasher.HashPassword(request.Password)
	if err != nil {
		return err
	}

	log.L(ctx).Infof("Updating password hash for user: %s", user.ID)

	if err := s.userRepository.UpdatePasswordHash(ctx, user.ID, passwordHash); err != nil {
		return err
	}

	log.L(ctx).Info("Password reset finished successfully")

	return nil
}

func (s *authService) Refresh(ctx context.Context, request domain.RefreshTokenRequest) (*domain.SessionResponse, error) {
	log.L(ctx).Infof("Refreshing session %s", request.RefreshToken)

	token, isReused, err := s.refreshTokenRepository.GetRefreshToken(ctx, request.RefreshToken)
	if err != nil {
		return nil, err
	}

	if token == nil {
		log.L(ctx).Errorf("Invalid refresh token (%s) was used", request.RefreshToken)

		return nil, appErr.Errorf(appErr.InvalidSessionError, "Invalid session")
	}

	if isReused {
		log.L(ctx).Warnf("Refresh Token Reuse Detected! User %s attempted to use a previously invalidated token (%s). Forcing logout of all sessions for this user!", token.UserID, request.RefreshToken)

		if err := s.refreshTokenRepository.RemoveAllRefreshTokensForUser(ctx, token.UserID); err != nil {
			log.L(ctx).Errorf("Failed to remove all refresh tokens for user %s after reuse detection: %v", token.UserID, err)
			return nil, err
		}

		return nil, appErr.Errorf(appErr.InvalidSessionError, "Compromised session. All sessions terminated.")
	}

	isRateLimited, err := s.ratelimiter.IsRateLimited(ctx, fmt.Sprintf("refresh:%s", token.Session.ID), time.Minute, 1)

	if err != nil {
		return nil, err
	}

	if isRateLimited {
		log.L(ctx).Warnf("Refresh token rate limit exceeded for user %s", token.UserID)

		return nil, appErr.Errorf(appErr.TooManyRequestsError, "Too many refresh attempts. Please try again later.")
	}

	err = s.refreshTokenRepository.RemoveRefreshToken(ctx, request.RefreshToken, token.UserID)
	if err != nil {
		return nil, err
	}

	session, err := s.createSession(ctx, token.UserID, request.UserAgent, request.IPAddress)
	if err != nil {
		return nil, err
	}

	log.L(ctx).Infof("Session %s refreshed successfully (%s) for user %s", request.RefreshToken, session.RefreshToken, token.UserID)

	return session, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	log.L(ctx).Infof("Invalidating session %s", refreshToken)

	token, isReused, err := s.refreshTokenRepository.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	if token == nil {
		log.L(ctx).Errorf("Invalid refresh token (%s) was used for logout (not found)", refreshToken)

		return appErr.Errorf(appErr.InvalidSessionError, "Invalid session")
	}

	if isReused {
		log.L(ctx).Warnf("Refresh Token Reuse Detected during logout! User %s attempted to use a previously invalidated token (%s). Forcing logout of all sessions for this user!", token.UserID, refreshToken)

		err := s.refreshTokenRepository.RemoveAllRefreshTokensForUser(ctx, token.UserID)
		if err != nil {
			return err
		}

		return appErr.Errorf(appErr.InvalidSessionError, "Compromised session. All sessions terminated.")
	}

	err = s.refreshTokenRepository.RemoveRefreshToken(ctx, refreshToken, token.UserID)
	if err != nil {
		return err
	}

	log.L(ctx).Infof("Session %s was invalidated successfully for user %s", refreshToken, token.UserID)

	return nil
}

func (s *authService) GetSessionsForAuthenticatedUser(ctx context.Context) ([]domain.Session, error) {
	log.L(ctx).Info("Fetching sessions for authenticated user")

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
			Name:      deviceInfo.Name,
			Os:        deviceInfo.Os,
			Device:    deviceInfo.Device,
			Location:  location,
			UserAgent: token.Session.UserAgent,
			IPAddress: token.Session.IPAddress,
			CreatedAt: token.Session.CreatedAt,
		}
	}

	log.L(ctx).Infof("Found %d sessions for user %s", len(sessionsResponse), userID)

	return sessionsResponse, nil
}

func (s *authService) RemoveSessionByID(ctx context.Context, sessionID uuid.UUID, currentRefreshToken string) error {
	log.L(ctx).Infof("Removing session %s for authenticated user", sessionID)

	userID, err := auth.GetUserID(ctx)
	if err != nil {
		return err
	}

	// Get the refresh token details for the session to be removed
	tokenToRemove, _, err := s.refreshTokenRepository.GetRefreshTokenBySessionID(ctx, sessionID)
	if err != nil {
		return err
	}

	if tokenToRemove == nil {
		log.L(ctx).Errorf("Attempted to remove non-existent session ID: %s", sessionID)
		return appErr.Errorf(appErr.InvalidSessionError, "Session not found")
	}

	if tokenToRemove.UserID != userID {
		log.L(ctx).Errorf("User %s attempted to remove session %s which does not belong to them", userID, sessionID)
		return appErr.Errorf(appErr.InvalidSessionError, "You are not authorized to remove this session")
	}

	// Prevent user from removing their current session
	if tokenToRemove.Token == currentRefreshToken {
		log.L(ctx).Warnf("User %s attempted to remove their currently active session %s", userID, sessionID)
		return appErr.Errorf(appErr.SessionNotFoundError, "Cannot remove your current active session")
	}

	err = s.refreshTokenRepository.RemoveRefreshToken(ctx, tokenToRemove.Token, userID)
	if err != nil {
		return err
	}

	log.L(ctx).Infof("Session %s removed successfully for user %s", sessionID, userID)

	return nil
}

func (s *authService) createSession(ctx context.Context, userID uuid.UUID, userAgent, ipAddress string) (*domain.SessionResponse, error) {
	log.L(ctx).Infof("Creating new session for user: %s", userID)

	csrfToken, err := s.tokenManager.NewCsrfToken()
	if err != nil {
		return nil, err
	}

	accessTokenExpiresAt := time.Now().Add(15 * time.Minute)
	accessToken, err := s.tokenManager.NewAccessToken(userID, csrfToken, accessTokenExpiresAt)
	if err != nil {
		return nil, err
	}

	log.L(ctx).Infof("Creating new device entry for user %s", userID)

	err = s.trustedDeviceRepository.CreateDevice(ctx, model.TrustedDevice{
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

	log.L(ctx).Infof("Adding refresh token to repository for user %s with session ID %s", userID, sessionID)

	err = s.refreshTokenRepository.AddRefreshToken(ctx, model.RefreshToken{
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

	log.L(ctx).Infof("Session created successfully for user %s with refresh token %s", userID, refreshToken)

	return &domain.SessionResponse{
		CsrfToken:             csrfToken,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessTokenExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	}, nil
}

func newSixDigitOtp() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
