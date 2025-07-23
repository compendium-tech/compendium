package auth

import (
	"context"
	"fmt"

	"github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type _isCsrfKey struct{}

var isCsrfKey = _isCsrfKey{}

const csrfTokenHeaderName = "X-Csrf-Token"

type AuthMiddleware struct {
	TokenManager TokenManager
}

func (a AuthMiddleware) Handle(c *gin.Context) {
	userID, isCsrfTokenValid := a.parseAccessTokenCookie(c)
	if userID == uuid.Nil {
		return
	}

	ctx := c.Request.Context()
	SetUserID(&ctx, userID)
	ctx = context.WithValue(ctx, isCsrfKey, isCsrfTokenValid)

	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

func RequireAuth(c *gin.Context) {
	_, err := GetUserID(c.Request.Context())

	if err != nil {
		log.L(c.Request.Context()).Warnf("Failed to require auth, check the previous logs to reveal the reason")
		c.AbortWithError(401, fmt.Errorf("invalid session"))
	}
}

func RequireCsrf(c *gin.Context) {
	isCsrfPresent, ok := c.Request.Context().Value(isCsrfKey).(bool)

	if !ok || !isCsrfPresent {
		c.AbortWithError(401, fmt.Errorf("invalid session"))
	}
}

func (a AuthMiddleware) parseAccessTokenCookie(c *gin.Context) (uuid.UUID, bool) {
	accessToken, err := c.Cookie("accessToken")

	if err != nil {
		return uuid.Nil, false
	}

	claims, err := a.TokenManager.ParseAccessToken(accessToken)
	if err != nil {
		return uuid.Nil, false
	}

	csrfToken := c.GetHeader(csrfTokenHeaderName)
	fmt.Println(csrfToken)
	if csrfToken != claims.CsrfToken {
		fmt.Println(claims.CsrfToken)
		return claims.Subject, false
	}

	return claims.Subject, true
}
