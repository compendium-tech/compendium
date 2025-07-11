package middleware

import (
	"context"
	"time"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/seacite-tech/compendium/common/pkg/log"
	"github.com/sirupsen/logrus"
)

const UserIdKey = "userId"

type LoggerMiddleware struct {
	LogProcessedRequests bool
	LogFinishedRequests  bool
}

func (l LoggerMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		entry := logrus.WithFields(logrus.Fields{
			"clientIp":  GetClientIP(c),
			"userId":    GetUserID(c),
			"method":    c.Request.Method,
			"path":      c.Request.RequestURI,
			"status":    c.Writer.Status(),
			"referrer":  c.Request.Referer(),
			"requestId": c.Writer.Header().Get("Request-Id"),
		})

		if l.LogProcessedRequests {
			entry.Info("")
		}

		ctx := context.WithValue(c.Request.Context(), log.LoggerKey, entry)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		duration := GetDurationInMillseconds(start)

		entry = entry.WithFields(logrus.Fields{
			"duration": duration,
		})

		if l.LogFinishedRequests {
			if c.Writer.Status() >= 500 {
				entry.Error(c.Errors.String())
			} else {
				entry.Info("")
			}
		}
	}
}

func GetClientIP(c *gin.Context) string {
	requester := c.Request.Header.Get("X-Forwarded-For")

	if len(requester) == 0 {
		requester = c.Request.Header.Get("X-Real-IP")
	}

	if len(requester) == 0 {
		requester = c.Request.RemoteAddr
	}

	// if requester is a comma delimited list, take the first one
	// (this happens when proxied via elastic load balancer then again through nginx)
	if strings.Contains(requester, ",") {
		requester = strings.Split(requester, ",")[0]
	}

	return requester
}

func GetUserID(c *gin.Context) string {
	userID, exists := c.Get(UserIdKey)

	if exists {
		return userID.(string)
	}

	return ""
}

func GetDurationInMillseconds(start time.Time) float64 {
	end := time.Now()
	duration := end.Sub(start)
	milliseconds := float64(duration) / float64(time.Millisecond)
	rounded := float64(int(milliseconds*100+.5)) / 100
	return rounded
}
