package auth

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/seacite-tech/compendium/common/pkg/auth"
	"github.com/seacite-tech/compendium/common/pkg/log"
	appErr "github.com/seacite-tech/compendium/user-service/internal/error"
)

type _isCsrfKey struct{}

var isCsrfKey = _isCsrfKey{}

const csrfTokenHeaderName = "X-Csrf-Token"

type AuthMiddleware struct {
	TokenManager TokenManager
}

func (a AuthMiddleware) Handle(c *gin.Context) {
	userId, isCsrfTokenValid := a.parseAccessTokenCookie(c)
	if userId == uuid.Nil {
		return
	}

	ctx := c.Request.Context()
	auth.SetUserId(&ctx, userId)
	ctx = context.WithValue(ctx, isCsrfKey, isCsrfTokenValid)

	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

var RequireAuth = appErr.HandleAppErr(requireAuth)
var RequireCsrf = appErr.HandleAppErr(requireCsrf)

func requireAuth(c *gin.Context) error {
	_, err := auth.GetUserId(c.Request.Context())

	if err != nil {
		log.L(c.Request.Context()).Warnf("Failed to require auth, check the previous logs to reveal the reason")
		return appErr.Errorf(appErr.InvalidSessionError, "Invalid session")
	}

	return nil
}

func requireCsrf(c *gin.Context) error {
	isCsrfPresent, ok := c.Request.Context().Value(isCsrfKey).(bool)

	if !ok || !isCsrfPresent {
		return appErr.Errorf(appErr.InvalidSessionError, "Invalid session")
	}

	return nil
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
