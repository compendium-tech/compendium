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
	IPAddress string
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
			AccessTokenExpiresAt:  &s.Session.AccessTokenExpiresAt,
			RefreshTokenExpiresAt: &s.Session.RefreshTokenExpiresAt,
			IsMfaRequired:         s.IsMfaRequired,
		}
	} else {
		return SignInResponseBody{
			AccessTokenExpiresAt:  nil,
			RefreshTokenExpiresAt: nil,
			IsMfaRequired:         s.IsMfaRequired,
		}
	}
}

type SignInResponseBody struct {
	AccessTokenExpiresAt  *time.Time `json:"accessTokenExpiresAt,omitempty"`
	RefreshTokenExpiresAt *time.Time `json:"refreshTokenExpiresAt,omitempty"`
	IsMfaRequired         bool       `json:"isMfaRequired"`
}

type SubmitMfaOtpRequest struct {
	Email     string
	Otp       string
	IPAddress string
	UserAgent string
}

type SubmitMfaOtpRequestBody struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required,len=6,number"`
}

type SessionResponse struct {
	CsrfToken             string
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

func (r SessionResponse) IntoBody() SessionResponseBody {
	return SessionResponseBody{
		AccessTokenExpiresAt:  r.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: r.RefreshTokenExpiresAt,
	}
}

type SessionResponseBody struct {
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
}

type RefreshTokenRequest struct {
	RefreshToken string
	UserAgent    string
	IPAddress    string
}

type RefreshTokenResponse struct {
	CsrfToken             string `json:"csrfToken"`
	AccessToken           string `json:"accessToken"`
	RefreshToken          string `json:"refreshToken"`
	AccessTokenExpiresAt  int64  `json:"accessTokenExpiresAt"`
	RefreshTokenExpiresAt int64  `json:"refreshTokenExpiresAt"`
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
