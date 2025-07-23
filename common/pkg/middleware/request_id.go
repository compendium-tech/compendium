package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
)

type RequestIDMiddleware struct {
	AllowToSet bool
}

func (r RequestIDMiddleware) Handle(c *gin.Context) {
	requestID := uuid.New()

	if r.AllowToSet {
		requestID = c.Request.Header.Get("Set-Request-ID")
	}

	c.Writer.Header().Set("Request-ID", requestID)
	c.Next()
}
