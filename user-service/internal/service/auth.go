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

type AuthService struct {
	emailLockRepository    repository.EmailLockRepository
	userRepository         repository.UserRepository
	mfaRepository          repository.MfaRepository
	refreshTokenRepository repository.RefreshTokenRepository
	emailSender            email.EmailSender
	tokenManager           auth.TokenManager
	passwordHasher         hash.PasswordHasher
}

func NewAuthService(
	emailLockRepository repository.EmailLockRepository,
	userRepository repository.UserRepository,
	mfaRepository repository.MfaRepository,
	refreshTokenRepository repository.RefreshTokenRepository,
	emailSender email.EmailSender,
	tokenManager auth.TokenManager,
	passwordHasher hash.PasswordHasher) AuthService {
	return AuthService{
		emailLockRepository:    emailLockRepository,
		userRepository:         userRepository,
		mfaRepository:          mfaRepository,
		refreshTokenRepository: refreshTokenRepository,
		emailSender:            emailSender,
		tokenManager:           tokenManager,
		passwordHasher:         passwordHasher,
	}
}

func (s AuthService) SignUp(ctx context.Context, request domain.SignUpRequest) (err error) {
	lock, err := s.emailLockRepository.ObtainEmailLock(ctx, request.Email)
	if err != nil {
		return err
	}

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
		err := s.userRepository.UpdatePasswordHashAndCreatedAtByEmail(ctx, request.Email, passwordHash, time.Now().UTC())

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

	err = s.emailSender.SendSignUpMfaEmail(request.Email, otp)
	if err != nil {
		return err
	}

	return lock.Release(ctx)
}

func (s AuthService) SubmitMfaOtp(ctx context.Context, request domain.SubmitMfaOtpRequest) (*domain.SessionResponse, error) {
	lock, err := s.emailLockRepository.ObtainEmailLock(ctx, request.Email)
	if err != nil {
		return nil, err
	}

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

	csrfToken, err := s.tokenManager.NewCsrfToken()
	if err != nil {
		return nil, err
	}

	accessTokenExpiry := time.Now().Add(15 * time.Minute)
	accessToken, err := s.tokenManager.NewAccessToken(user.Id, csrfToken, accessTokenExpiry)
	if err != nil {
		return nil, err
	}

	refreshTokenExpiry := time.Now().Add(20 * 24 * time.Hour)
	refreshToken, err := s.tokenManager.NewRefreshToken()
	if err != nil {
		return nil, err
	}

	err = s.refreshTokenRepository.AddRefreshToken(ctx, model.RefreshToken{
		UserId: user.Id,
		Token:  refreshToken,
		Expiry: refreshTokenExpiry,
	})
	if err != nil {
		return nil, err
	}

	err = lock.Release(ctx)
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
