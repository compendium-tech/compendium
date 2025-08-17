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
	NewAccessToken(userID uuid.UUID, csrfToken string, expiresAt time.Time) string
	NewCsrfToken() string
	NewRefreshToken() string

	ParseAccessToken(token string) (*JwtTokenClaims, error)
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
	Issuer    string    `json:"iss,omitempty"`
	ExpiresAt int64     `json:"exp,omitempty"`
	Subject   uuid.UUID `json:"sub,omitempty"`
	CsrfToken string    `json:"csrf,omitempty"`
}

func (c JwtTokenClaims) Valid() error {
	return jwt.StandardClaims{
		Issuer:    c.Issuer,
		ExpiresAt: c.ExpiresAt,
		Subject:   c.Subject.String(),
	}.Valid()
}

func (m *JwtBasedTokenManager) NewAccessToken(userID uuid.UUID, csrfToken string, expiresAt time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtTokenClaims{
		Issuer:    "user-service",
		ExpiresAt: expiresAt.Unix(),
		Subject:   userID,
		CsrfToken: csrfToken,
	})

	jwt, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		panic(err)
	}

	return jwt
}

func (m *JwtBasedTokenManager) ParseAccessToken(accessToken string) (*JwtTokenClaims, error) {
	claims := JwtTokenClaims{}
	_, err := jwt.ParseWithClaims(accessToken, &claims, func(token *jwt.Token) (i any, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, tracerr.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})

	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	return &claims, nil
}

func (m *JwtBasedTokenManager) NewCsrfToken() string {
	return newRandomString(16)
}

func (m *JwtBasedTokenManager) NewRefreshToken() string {
	return newRandomString(32)
}

func newRandomString(size int) string {
	b := make([]byte, size)

	s := rand.NewSource(time.Now().UnixMicro())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", b)
}
