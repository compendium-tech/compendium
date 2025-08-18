package auth

import (
	"context"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	common "github.com/compendium-tech/compendium/common/pkg"
	"github.com/compendium-tech/compendium/common/pkg/log"
)

type _isCsrfKey struct{}

var isCsrfKey = _isCsrfKey{}

const csrfTokenHeaderName = "X-Csrf-Token"

type Middleware struct {
	TokenManager TokenManager
}

func (a Middleware) Handle(c *gin.Context) {
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
	defer func() {
		if r := recover(); r != nil {
			log.L(c.Request.Context()).
				WithField("stack", debug.Stack()).
				Warnf("Failed to require auth, check the previous logs to reveal the reason")

			c.AbortWithStatusJSON(http.StatusUnauthorized, common.H{
				"errorType": 8,
			})
		}
	}()

	_ = GetUserID(c.Request.Context())
}

func RequireCsrf(c *gin.Context) {
	isCsrfPresent, ok := c.Request.Context().Value(isCsrfKey).(bool)

	if !ok || !isCsrfPresent {
		log.L(c.Request.Context()).Warnf("Failed to require csrf token, check the previous logs to reveal the reason")
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.H{
			"errorType": 8,
		})
	}
}

func (a Middleware) parseAccessTokenCookie(c *gin.Context) (uuid.UUID, bool) {
	accessToken, err := c.Cookie("accessToken")

	if err != nil {
		return uuid.Nil, false
	}

	claims, err := a.TokenManager.ParseAccessToken(accessToken)
	if err != nil {
		return uuid.Nil, false
	}

	csrfToken := c.GetHeader(csrfTokenHeaderName)
	if csrfToken != claims.CsrfToken {
		return claims.Subject, false
	}

	return claims.Subject, true
}
