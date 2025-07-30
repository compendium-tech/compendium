package middleware

import (
	"time"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/http"
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

	userUuid, _ := auth.GetUserID(c)

	var userID string

	if userUuid == uuid.Nil {
		userID = "unauthenticated"
	} else {
		userID = userUuid.String()
	}

	entry := logrus.WithFields(logrus.Fields{
		"clientIp":  httputils.GetClientIP(c),
		"userId":    userID,
		"method":    c.Request.Method,
		"path":      c.Request.RequestURI,
		"status":    c.Writer.Status(),
		"referrer":  c.Request.Referer(),
		"requestID": c.Writer.Header().Get("Request-ID"),
	})

	if l.LogProcessedRequests {
		entry.Info("Request processing is started")
	}

	ctx := c.Request.Context()
	log.SetLogger(&ctx, entry)
	c.Request = c.Request.WithContext(ctx)

	c.Next()

	duration := GetDurationInMilliseconds(start)

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

func GetDurationInMilliseconds(start time.Time) float64 {
	end := time.Now()
	duration := end.Sub(start)
	milliseconds := float64(duration) / float64(time.Millisecond)
	rounded := float64(int(milliseconds*100+.5)) / 100
	return rounded
}
