package auth

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/ztrue/tracerr"
)

type TokenManager interface {
	NewAccessToken(userId uuid.UUID, csrfToken string, expiresAt time.Time) (string, error)
	ParseAccessToken(token string) (*JwtTokenClaims, error)
	NewCsrfToken() (string, error)
	NewRefreshToken() (string, error)
}

type JwtBasedTokenManager struct {
	signingKey string
}

func NewJwtBasedTokenManager(signingKey string) (*JwtBasedTokenManager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &JwtBasedTokenManager{signingKey: signingKey}, nil
}

type JwtTokenClaims struct {
	Id        string `json:"jti,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	Subject   string `json:"sub,omitempty"`
	CsrfToken string `json:"csrf,omitempty"`
}

func (c JwtTokenClaims) Valid() error {
	return jwt.StandardClaims{
		Issuer:    c.Issuer,
		ExpiresAt: c.ExpiresAt,
		Subject:   c.Subject,
	}.Valid()
}

func (m *JwtBasedTokenManager) NewAccessToken(userId uuid.UUID, csrfToken string, expiresAt time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtTokenClaims{
		Issuer:    "user-service",
		ExpiresAt: expiresAt.Unix(),
		Subject:   userId.String(),
		CsrfToken: csrfToken,
	})

	jwt, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	return jwt, nil
}

func (m *JwtBasedTokenManager) ParseAccessToken(accessToken string) (*JwtTokenClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, tracerr.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	claims, ok := token.Claims.(JwtTokenClaims)
	if !ok {
		return nil, tracerr.Errorf("failed to read user claims from token")
	}

	return &claims, nil
}

func (m *JwtBasedTokenManager) NewCsrfToken() (string, error) {
	return newRandomString(16)
}

func (m *JwtBasedTokenManager) newCsrfTokenHashSalt() (string, error) {
	return newRandomString(10)
}

func (m *JwtBasedTokenManager) NewRefreshToken() (string, error) {
	return newRandomString(32)
}

func newRandomString(size int) (string, error) {
	b := make([]byte, size)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", tracerr.Wrap(err)
	}

	return fmt.Sprintf("%x", b), nil
}
