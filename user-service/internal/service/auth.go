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
	apperr "github.com/seacite-tech/compendium/user-service/internal/error"
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
}

type authServiceImpl struct {
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
	passwordHasher hash.PasswordHasher) *authServiceImpl {
	return &authServiceImpl{
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

func (s *authServiceImpl) SignUp(ctx context.Context, request domain.SignUpRequest) (finalErr error) {
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
			log.L(ctx).Error("User tried to sign up even though email is taken")

			return apperr.Errorf(apperr.EmailTakenError, "Email is already taken")
		}

		if time.Now().UTC().Sub(user.CreatedAt) < 60*time.Second {
			log.L(ctx).Error("User tried to repeat sign up attempt in less than 60 seconds")

			return apperr.Errorf(apperr.TooManyRequestsError, "Too many sign up attempts")
		}
	}

	passwordHash, err := s.passwordHasher.HashPassword(request.Password)
	if err != nil {
		return err
	}

	if user != nil && !user.IsEmailVerified {
		err := s.userRepository.UpdatePasswordHashAndCreatedAt(ctx, user.Id, passwordHash, time.Now().UTC())

		if err != nil {
			return err
		}
	} else {
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

	otp := newSixDigitOtp()
	err = s.mfaRepository.SetMfaOtpByEmail(ctx, request.Email, otp)
	if err != nil {
		return err
	}

	return s.emailSender.SendSignUpMfaEmail(request.Email, otp)
}

func (s *authServiceImpl) SubmitMfaOtp(ctx context.Context, request domain.SubmitMfaOtpRequest) (_ *domain.SessionResponse, finalErr error) {
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
		return nil, apperr.Errorf(apperr.MfaNotRequestedError, "2FA was not requested")
	}

	otp, err := s.mfaRepository.GetMfaOtpByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	if otp == nil {
		return nil, apperr.Errorf(apperr.MfaNotRequestedError, "2FA was not requested")
	}

	if *otp != request.Otp {
		return nil, apperr.Errorf(apperr.InvalidMfaOtpError, "Invalid 2FA otp")
	}

	if !user.IsEmailVerified {
		err = s.userRepository.UpdateIsEmailVerifiedByEmail(ctx, request.Email, true)

		if err != nil {
			return nil, err
		}
	}

	err = s.mfaRepository.RemoveMfaOtpByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}

	session, err := s.createSession(ctx, user.Id, request.UserAgent, request.IpAddress)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *authServiceImpl) SignIn(ctx context.Context, request domain.SignInRequest) (_ *domain.SignInResponse, finalErr error) {
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
		return nil, apperr.Errorf(apperr.InvalidCredentialsError, "Invalid credentials")
	}

	isPasswordValid, err := s.passwordHasher.IsPasswordHashValid(user.PasswordHash, request.Password)
	if err != nil {
		return nil, err
	}

	if !isPasswordValid {
		return nil, apperr.Errorf(apperr.InvalidCredentialsError, "Invalid credentials")
	}

	deviceIsKnown, err := s.deviceRepository.DeviceExists(user.Id, request.UserAgent, request.IpAddress)
	if err != nil {
		return nil, err
	}

	session := (*domain.SessionResponse)(nil)

	if !deviceIsKnown {
		otp := newSixDigitOtp()
		err = s.mfaRepository.SetMfaOtpByEmail(ctx, request.Email, otp)
		if err != nil {
			return nil, err
		}

		err = s.emailSender.SendSignInMfaEmail(request.Email, otp)
		if err != nil {
			return nil, err
		}
	} else {
		session, err = s.createSession(ctx, user.Id, request.UserAgent, request.IpAddress)

		if err != nil {
			return nil, err
		}
	}

	if !deviceIsKnown {
		return &domain.SignInResponse{IsMfaRequired: true, Session: nil}, nil
	} else {
		return &domain.SignInResponse{IsMfaRequired: false, Session: session}, nil
	}
}

func (s *authServiceImpl) InitPasswordReset(ctx context.Context, request domain.InitPasswordResetRequest) (finalErr error) {
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
		return apperr.Errorf(apperr.UserNotFoundError, "User with given email address doesn't exist")
	}

	otp := newSixDigitOtp()
	err = s.mfaRepository.SetMfaOtpByEmail(ctx, request.Email, otp)
	if err != nil {
		return err
	}

	return s.emailSender.SendSignInMfaEmail(request.Email, otp)
}

func (s *authServiceImpl) FinishPasswordReset(ctx context.Context, request domain.FinishPasswordResetRequest) (finalErr error) {
	lock, err := s.authLockRepository.ObtainLock(ctx, request.Email)
	if err != nil {
		return err
	}

	defer func() {
		lock.ReleaseAndHandleErr(ctx, &finalErr)
	}()

	user, err := s.userRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	if user == nil || !user.IsEmailVerified {
		return apperr.Errorf(apperr.MfaNotRequestedError, "2FA was not requested")
	}

	otp, err := s.mfaRepository.GetMfaOtpByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	if otp == nil {
		return apperr.Errorf(apperr.MfaNotRequestedError, "2FA was not requested")
	}

	if *otp != request.Otp {
		return apperr.Errorf(apperr.InvalidMfaOtpError, "Invalid 2FA otp")
	}

	err = s.mfaRepository.RemoveMfaOtpByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	passwordHash, err := s.passwordHasher.HashPassword(request.Password)
	if err != nil {
		return err
	}

	return s.userRepository.UpdatePasswordHash(ctx, user.Id, passwordHash)
}

func (s *authServiceImpl) createSession(ctx context.Context, userId uuid.UUID, userAgent, ipAddress string) (*domain.SessionResponse, error) {
	csrfToken, err := s.tokenManager.NewCsrfToken()
	if err != nil {
		return nil, err
	}

	err = s.deviceRepository.CreateDevice(ctx, model.Device{
		Id:        uuid.New(),
		UserId:    userId,
		UserAgent: userAgent,
		IpAddress: ipAddress,
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

	err = s.refreshTokenRepository.AddRefreshToken(ctx, model.RefreshToken{
		UserId: userId,
		Token:  refreshToken,
		Expiry: refreshTokenExpiry,
	})
	if err != nil {
		return nil, err
	}

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
