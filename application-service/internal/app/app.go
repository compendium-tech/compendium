package app

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/compendium-tech/compendium/common/pkg/middleware"

	"github.com/compendium-tech/compendium/application-service/internal/config"
	httpv1 "github.com/compendium-tech/compendium/application-service/internal/controller/http/v1"
	"github.com/compendium-tech/compendium/application-service/internal/interop"
	"github.com/compendium-tech/compendium/application-service/internal/repository"
	"github.com/compendium-tech/compendium/application-service/internal/service"
)

type Dependencies struct {
	Config       *config.AppConfig
	PgDB         *sql.DB
	TokenManager auth.TokenManager
	LLMService   interop.LLMService
}

func NewApp(deps Dependencies) *gin.Engine {
	logrus.SetFormatter(&log.LogFormatter{
		Program:     "application-service",
		Environment: deps.Config.Environment,
	})
	logrus.SetReportCaller(true)

	applicationRepository := repository.NewPgApplicationRepository(deps.PgDB)
	applicationService := service.NewApplicationService(applicationRepository)
	applicationEvaluationService := service.NewApplicationEvaluateService(applicationRepository, deps.LLMService)

	r := gin.Default()
	r.Use(middleware.RequestIDMiddleware{AllowToSet: false}.Handle)
	r.Use(auth.Middleware{TokenManager: deps.TokenManager}.Handle)
	r.Use(middleware.LoggerMiddleware{LogProcessedRequests: true, LogFinishedRequests: true}.Handle)

	httpv1.NewApplicationController(applicationService).MakeRoutes(r)
	httpv1.NewApplicationEvaluationController(applicationService, applicationEvaluationService).MakeRoutes(r)

	return r
}
