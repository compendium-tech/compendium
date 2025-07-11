package domain

import "time"

type SignUpRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=100"`
	Password string `json:"password" validate:"required,password,min=6,max=100"`
	Email    string `json:"email" validate:"required,email"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type SubmitMfaOtpRequest struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required,len=6,number"`
}

type SessionResponse struct {
	CsrfToken          string    `json:"csrfToken"`
	AccessToken        string    `json:"accessToken"`
	RefreshToken       string    `json:"refreshToken"`
	AccessTokenExpiry  time.Time `json:"accessTokenExpiry"`
	RefreshTokenExpiry time.Time `json:"refreshTokenExpiry"`
}

func (r SessionResponse) JsonResponse() SessionJsonResponse {
	return SessionJsonResponse{
		AccessTokenExpiry:  r.AccessTokenExpiry,
		RefreshTokenExpiry: r.RefreshTokenExpiry,
	}
}

type SessionJsonResponse struct {
	AccessTokenExpiry  time.Time `json:"accessTokenExpiry"`
	RefreshTokenExpiry time.Time `json:"refreshTokenExpiry"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken        string `json:"accessToken"`
	RefreshToken       string `json:"refreshToken"`
	AccessTokenExpiry  int64  `json:"accessTokenExpiry"`
	RefreshTokenExpiry int64  `json:"refreshTokenExpiry"`
}
