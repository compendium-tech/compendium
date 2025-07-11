package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seacite-tech/compendium/common/pkg/middleware"
)

const csrfTokenHeaderName = "X-Csrf-Token"

func NoCsrfAuth(tokenManager TokenManager) noCsrfAuthMiddleware {
	return noCsrfAuthMiddleware{tokenManager: tokenManager}
}

func CsrfAuth(tokenManager TokenManager) csrfAuthMiddleware {
	return csrfAuthMiddleware{tokenManager: tokenManager}
}

type noCsrfAuthMiddleware struct {
	tokenManager TokenManager
}

type csrfAuthMiddleware struct {
	tokenManager TokenManager
}

func (a noCsrfAuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := a.parseAuthHeader(c)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.Set(middleware.UserIdKey, userId)
		c.Next()
	}
}

func (a noCsrfAuthMiddleware) parseAuthHeader(c *gin.Context) (string, error) {
	accessToken, err := c.Cookie("accessToken")
	if err != nil {
		return "", fmt.Errorf("access token cookie is not found")
	}

	claims, err := a.tokenManager.ParseAccessToken(accessToken)
	if err != nil {
		return "", fmt.Errorf("failed to parse access token")
	}

	return claims.Subject, nil
}

func (a csrfAuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := a.parseAuthHeader(c)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.Set(middleware.UserIdKey, userId)
		c.Next()
	}
}

func (a csrfAuthMiddleware) parseAuthHeader(c *gin.Context) (string, error) {
	accessToken, err := c.Cookie("accessToken")
	if err != nil {
		return "", fmt.Errorf("access token cookie is not found")
	}

	claims, err := a.tokenManager.ParseAccessToken(accessToken)
	if err != nil {
		return "", fmt.Errorf("failed to parse access token")
	}

	csrfToken := c.GetHeader(csrfTokenHeaderName)
	if csrfToken != claims.CsrfToken {
		return "", fmt.Errorf("invalid csrf token")
	}

	return claims.Subject, nil
}
