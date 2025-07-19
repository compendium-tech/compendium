package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Cors struct {
	AllowOrigin      string
	MaxAge           uint64
	AllowMethods     string
	AllowHeaders     string
	ExposeHeaders    string
	AllowCredentials bool
}

func DefaultCors() Cors {
	return Cors{
		AllowOrigin:      "http://localhost:5173",
		MaxAge:           1728000,
		AllowMethods:     "GET, POST, OPTIONS",
		AllowHeaders:     "DNT,X-Csrf-Token,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Cookie",
		ExposeHeaders:    "*",
		AllowCredentials: true,
	}
}

func (c Cors) Handle(gc *gin.Context) {
	gc.Writer.Header().Set("Access-Control-Allow-Origin", "*") // allow any origin domain
	if c.AllowOrigin != "" {
		gc.Writer.Header().Set("Access-Control-Allow-Origin", c.AllowOrigin)
	}

	gc.Writer.Header().Set("Access-Control-Allow-Methods", c.AllowMethods)
	gc.Writer.Header().Set("Access-Control-Allow-Headers", c.AllowHeaders)
	gc.Writer.Header().Set("Access-Control-Expose-Headers", c.ExposeHeaders)
	gc.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	gc.Writer.Header().Set("Access-Control-Max-Age", strconv.FormatUint(c.MaxAge, 10))

	if gc.Request.Method == "OPTIONS" {
		gc.AbortWithStatus(200)
	} else {
		gc.Next()
	}
}
