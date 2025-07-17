package domain

import "time"

type SignUpRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=100"`
	Password string `json:"password" validate:"required,password,min=6,max=100"`
	Email    string `json:"email" validate:"required,email"`
}

type SignInRequest struct {
	Email     string
	Password  string
	IpAddress string
	UserAgent string
}

type SignInRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type SignInResponse struct {
	Session       *SessionResponse
	IsMfaRequired bool
}

func (s SignInResponse) IntoBody() SignInResponseBody {
	if s.Session != nil {
		return SignInResponseBody{
			AccessTokenExpiry:  &s.Session.AccessTokenExpiry,
			RefreshTokenExpiry: &s.Session.RefreshTokenExpiry,
			IsMfaRequired:      s.IsMfaRequired,
		}
	} else {
		return SignInResponseBody{
			AccessTokenExpiry:  nil,
			RefreshTokenExpiry: nil,
			IsMfaRequired:      s.IsMfaRequired,
		}
	}
}

type SignInResponseBody struct {
	AccessTokenExpiry  *time.Time `json:"accessTokenExpiry,omitempty"`
	RefreshTokenExpiry *time.Time `json:"refreshTokenExpiry,omitempty"`
	IsMfaRequired      bool       `json:"isMfaRequired"`
}

type SubmitMfaOtpRequest struct {
	Email     string
	Otp       string
	IpAddress string
	UserAgent string
}

type SubmitMfaOtpRequestBody struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required,len=6,number"`
}

type SessionResponse struct {
	CsrfToken          string
	AccessToken        string
	RefreshToken       string
	AccessTokenExpiry  time.Time
	RefreshTokenExpiry time.Time
}

func (r SessionResponse) IntoBody() SessionResponseBody {
	return SessionResponseBody{
		AccessTokenExpiry:  r.AccessTokenExpiry,
		RefreshTokenExpiry: r.RefreshTokenExpiry,
	}
}

type SessionResponseBody struct {
	AccessTokenExpiry  time.Time `json:"accessTokenExpiry"`
	RefreshTokenExpiry time.Time `json:"refreshTokenExpiry"`
}

type RefreshTokenRequest struct {
	RefreshToken string
	UserAgent    string
	IpAddress    string
}

type RefreshTokenResponse struct {
	CsrfToken          string `json:"csrfToken"`
	AccessToken        string `json:"accessToken"`
	RefreshToken       string `json:"refreshToken"`
	AccessTokenExpiry  int64  `json:"accessTokenExpiry"`
	RefreshTokenExpiry int64  `json:"refreshTokenExpiry"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type InitPasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type FinishPasswordResetRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Otp      string `json:"otp" validate:"required,len=6,number"`
	Password string `json:"password" validate:"required,password,min=6,max=100"`
}
