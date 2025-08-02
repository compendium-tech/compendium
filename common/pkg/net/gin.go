package netapp

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GinApp struct {
	engine *gin.Engine
}

func NewGinApp(engine *gin.Engine) GinApp {
	return GinApp{
		engine: engine,
	}
}

func (a GinApp) Run() error {
	port, err := getPortEnv()
	if err != nil {
		return err
	}

	logrus.Infof("Starting HTTP server on :%d", port)
	return a.engine.Run(fmt.Sprintf(":%d", port))
}

func (a GinApp) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.engine.ServeHTTP(w, req)
}
