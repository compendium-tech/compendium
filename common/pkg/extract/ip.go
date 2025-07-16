package extract

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetClientIP(c *gin.Context) string {
	requester := c.GetHeader("X-Forwarded-For")

	if len(requester) == 0 {
		requester = c.GetHeader("X-Real-IP")
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

func GetUserAgent(c *gin.Context) string {
	return c.GetHeader("User-Agent")
}
