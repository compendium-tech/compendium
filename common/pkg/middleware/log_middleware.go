package middleware

import (
	"context"
	"time"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/httphelp"
	"github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type LoggerMiddleware struct {
	LogProcessedRequests bool
	LogFinishedRequests  bool
}

func (l LoggerMiddleware) Handle(c *gin.Context) {
	start := time.Now()

	userUuid, _ := auth.GetUserId(c)

	var userId string

	if userUuid == uuid.Nil {
		userId = "unauthenticated"
	} else {
		userId = userUuid.String()
	}

	entry := logrus.WithFields(logrus.Fields{
		"clientIp":  httphelp.GetClientIP(c),
		"userId":    userId,
		"method":    c.Request.Method,
		"path":      c.Request.RequestURI,
		"status":    c.Writer.Status(),
		"referrer":  c.Request.Referer(),
		"requestId": c.Writer.Header().Get("Request-Id"),
	})

	if l.LogProcessedRequests {
		entry.Info("Request processing is started")
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
			entry.Info("Request processing is finished")
		}
	}
}

func GetDurationInMillseconds(start time.Time) float64 {
	end := time.Now()
	duration := end.Sub(start)
	milliseconds := float64(duration) / float64(time.Millisecond)
	rounded := float64(int(milliseconds*100+.5)) / 100
	return rounded
}
