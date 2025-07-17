package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
)

type RequestIdMiddleware struct {
	AllowToSet bool
}

func (r RequestIdMiddleware) Handle(c *gin.Context) {
	requestId := uuid.New()

	if r.AllowToSet {
		requestId = c.Request.Header.Get("Set-Request-Id")
	}

	c.Writer.Header().Set("Request-Id", requestId)
	c.Next()
}
