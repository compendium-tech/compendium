package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	log "github.com/seacite-tech/compendium/common/pkg/log"
	"github.com/seacite-tech/compendium/user-service/internal/domain"
	"github.com/seacite-tech/compendium/user-service/internal/email"
	appErr "github.com/seacite-tech/compendium/user-service/internal/error"
	"github.com/seacite-tech/compendium/user-service/internal/hash"
	"github.com/seacite-tech/compendium/user-service/internal/model"
	"github.com/seacite-tech/compendium/user-service/internal/repository"
	"github.com/seacite-tech/compendium/user-service/pkg/auth"
)

type AuthService interface {
	SignUp(ctx context.Context, request domain.SignUpRequest) error
	SubmitMfaOtp(ctx context.Context, request domain.SubmitMfaOtpRequest) (*domain.SessionResponse, error)
	SignIn(ctx context.Context, request domain.SignInRequest) (*domain.SignInResponse, error)
	InitPasswordReset(ctx context.Context, request domain.InitPasswordResetRequest) error
	FinishPasswordReset(ctx context.Context, request domain.FinishPasswordResetRequest) error
	Refresh(ctx context.Context, request domain.RefreshTokenRequest) (*domain.SessionResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}

type authService struct {
	authLockRepository     repository.AuthLockRepository
	deviceRepository       repository.DeviceRepository
	userRepository         repository.UserRepository
	mfaRepository          repository.MfaRepository
	refreshTokenRepository repository.RefreshTokenRepository
	emailSender            email.EmailSender
	tokenManager           auth.TokenManager
	passwordHasher         hash.PasswordHasher
}

func NewAuthService(
	authEmailLockRepository repository.AuthLockRepository,
	deviceRepository repository.DeviceRepository,
	userRepository repository.UserRepository,
	mfaRepository repository.MfaRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	emailSender email.EmailSender,
	tokenManager auth.TokenManager,
	passwordHasher hash.PasswordHasher) AuthService {
	return &authService{
		authLockRepository:     authEmailLockRepository,
		deviceRepository:       deviceRepository,
		userRepository:         userRepository,
		mfaRepository:          mfaRepository,
		refreshTokenRepository: refreshTokenRepository,
		emailSender:            emailSender,
		tokenManager:           tokenManager,
		passwordHasher:         passwordHasher,
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

		err := s.userRepository.UpdatePasswordHashAndCreatedAt(ctx, user.Id, passwordHash, time.Now().UTC())
		if err != nil {
			return err
		}
	} else {
		log.L(ctx).Infof("Creating new user with email: %s", request.Email)

		err := s.userRepository.CreateUser(ctx, model.User{
			Id:              uuid.New(),
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

	if err := s.emailSender.SendSignUpMfaEmail(request.Email, otp); err != nil {
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

	session, err := s.createSession(ctx, user.Id, request.UserAgent, request.IpAddress)
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

	deviceIsKnown, err := s.deviceRepository.DeviceExists(user.Id, request.UserAgent, request.IpAddress)
	if err != nil {
		return nil, err
	}

	session := (*domain.SessionResponse)(nil)

	if !deviceIsKnown {
		log.L(ctx).Infof("Device is not known for user %s. Initiating MFA", user.Id)

		otp := newSixDigitOtp()
		err = s.mfaRepository.SetMfaOtpByEmail(ctx, request.Email, otp)
		if err != nil {
			return nil, err
		}

		err = s.emailSender.SendSignInMfaEmail(request.Email, otp)
		if err != nil {
			return nil, err
		}

		log.L(ctx).Info("Sign-in is successfully processed, MFA required")

		return &domain.SignInResponse{IsMfaRequired: true, Session: nil}, nil
	} else {
		log.L(ctx).Infof("Device is known for user %s. Creating session", user.Id)
		session, err = s.createSession(ctx, user.Id, request.UserAgent, request.IpAddress)
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

	if err := s.emailSender.SendSignInMfaEmail(request.Email, otp); err != nil {
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

	log.L(ctx).Infof("Updating password hash for user: %s", user.Id)

	if err := s.userRepository.UpdatePasswordHash(ctx, user.Id, passwordHash); err != nil {
		return err
	}

	log.L(ctx).Info("Password reset finished successfully")

	return nil
}

func (s *authService) Refresh(ctx context.Context, request domain.RefreshTokenRequest) (*domain.SessionResponse, error) {
	log.L(ctx).Infof("Refreshing session %s", request.RefreshToken)

	userId, isReused, err := s.refreshTokenRepository.TryRemoveRefreshTokenByToken(ctx, request.RefreshToken)
	if err != nil {
		return nil, err
	}

	if userId == uuid.Nil {
		log.L(ctx).Errorf("Invalid refresh token (%s) was used", request.RefreshToken)

		return nil, appErr.Errorf(appErr.InvalidSessionError, "Invalid session")
	}

	if isReused {
		log.L(ctx).Warnf("Refresh Token Reuse Detected! User %s attempted to use a previously invalidated token (%s). Forcing logout of all sessions for this user!", userId, request.RefreshToken)

		if err := s.refreshTokenRepository.RemoveAllRefreshTokensForUser(ctx, userId); err != nil {
			log.L(ctx).Errorf("Failed to remove all refresh tokens for user %s after reuse detection: %v", userId, err)
			return nil, err
		}

		return nil, appErr.Errorf(appErr.InvalidSessionError, "Compromised session. All sessions terminated.")
	}

	session, err := s.createSession(ctx, userId, request.UserAgent, request.IpAddress)
	if err != nil {
		return nil, err
	}

	log.L(ctx).Infof("Session %s refreshed successfully (%s) for user %s", request.RefreshToken, session.RefreshToken, userId)

	return session, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	log.L(ctx).Infof("Invalidating session %s", refreshToken)

	userId, isReused, err := s.refreshTokenRepository.TryRemoveRefreshTokenByToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	if userId == uuid.Nil {
		log.L(ctx).Errorf("Invalid refresh token (%s) was used", refreshToken)

		return appErr.Errorf(appErr.InvalidSessionError, "Invalid session")
	}

	if isReused {
		log.L(ctx).Warnf("Refresh Token Reuse Detected! User %s attempted to use a previously invalidated token (%s). Forcing logout of all sessions for this user!", userId, refreshToken)

		if err := s.refreshTokenRepository.RemoveAllRefreshTokensForUser(ctx, userId); err != nil {
			log.L(ctx).Errorf("Failed to remove all refresh tokens for user %s after reuse detection: %v", userId, err)
			return err
		}

		return appErr.Errorf(appErr.InvalidSessionError, "Compromised session. All sessions terminated.")
	}

	log.L(ctx).Infof("Session %s was invalidated successfully", refreshToken)

	return nil
}

func (s *authService) createSession(ctx context.Context, userId uuid.UUID, userAgent, ipAddress string) (*domain.SessionResponse, error) {
	log.L(ctx).Infof("Creating new session for user: %s", userId)

	csrfToken, err := s.tokenManager.NewCsrfToken()
	if err != nil {
		return nil, err
	}

	log.L(ctx).Infof("Creating new device entry for user %s", userId)

	err = s.deviceRepository.CreateDevice(ctx, model.Device{
		Id:        uuid.New(),
		UserId:    userId,
		UserAgent: userAgent,
		IpAddress: ipAddress,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}

	accessTokenExpiry := time.Now().Add(15 * time.Minute)
	accessToken, err := s.tokenManager.NewAccessToken(userId, csrfToken, accessTokenExpiry)
	if err != nil {
		return nil, err
	}

	refreshTokenExpiry := time.Now().Add(20 * 24 * time.Hour)
	refreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, err
	}

	log.L(ctx).Infof("Adding refresh token to repository for user %s", userId)

	err = s.refreshTokenRepository.AddRefreshToken(ctx, model.RefreshToken{
		UserId: userId,
		Token:  refreshToken,
		Expiry: refreshTokenExpiry,
	})
	if err != nil {
		return nil, err
	}

	log.L(ctx).Infof("Session created successfully for user %s", userId)

	return &domain.SessionResponse{
		CsrfToken:          csrfToken,
		AccessToken:        accessToken,
		AccessTokenExpiry:  accessTokenExpiry,
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: refreshTokenExpiry,
	}, nil
}

func newSixDigitOtp() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
