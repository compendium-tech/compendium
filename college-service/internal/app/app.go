package app

import (
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/compendium-tech/compendium/common/pkg/middleware"
	netapp "github.com/compendium-tech/compendium/common/pkg/net"

	"github.com/compendium-tech/compendium/college-service/internal/config"
	httpv1 "github.com/compendium-tech/compendium/college-service/internal/delivery/http/v1"
	"github.com/compendium-tech/compendium/college-service/internal/repository"
	"github.com/compendium-tech/compendium/college-service/internal/service"
)

type Dependencies struct {
	Config              *config.AppConfig
	ElasticsearchClient *elasticsearch.Client
	TokenManager        auth.TokenManager
}

func NewApp(deps Dependencies) netapp.GinApp {
	logrus.SetFormatter(&log.LogFormatter{
		Program:     "application-service",
		Environment: deps.Config.Environment,
	})
	logrus.SetReportCaller(true)

	collegeRepository := repository.NewElasticsearchCollegeRepository(deps.ElasticsearchClient)
	collegeService := service.NewCollegeService(collegeRepository)

	r := gin.Default()
	r.Use(middleware.RequestIDMiddleware{AllowToSet: false}.Handle)
	r.Use(auth.Middleware{TokenManager: deps.TokenManager}.Handle)
	r.Use(middleware.LoggerMiddleware{LogProcessedRequests: true, LogFinishedRequests: true}.Handle)

	httpv1.NewCollegeController(collegeService).MakeRoutes(r)

	return netapp.NewGinApp(r)
}
