package httpv1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	httputils "github.com/compendium-tech/compendium/common/pkg/http"

	"github.com/compendium-tech/compendium/application-service/internal/middleware"
	"github.com/compendium-tech/compendium/application-service/internal/service"
)

type ApplicationEvaluationController struct {
	applicationService           service.ApplicationService
	applicationEvaluationService service.ApplicationEvaluationService
}

func NewApplicationEvaluationController(
	applicationService service.ApplicationService,
	applicationEvaluationService service.ApplicationEvaluationService) ApplicationEvaluationController {
	return ApplicationEvaluationController{
		applicationService:           applicationService,
		applicationEvaluationService: applicationEvaluationService,
	}
}

func (a ApplicationEvaluationController) MakeRoutes(e *gin.Engine) {
	var eh httputils.ErrorHandler

	v1 := e.Group("/v1")
	{
		authenticated := v1.Group("/")
		authenticated.Use(auth.RequireAuth)
		{
			application := authenticated.Group("/applications/:applicationId")
			application.Use(middleware.NewSetApplicationFromRequest(a.applicationService).Handle)
			{
				application.POST("/evaluations", eh.Handle(a.evaluateApplication))
			}
		}
	}
}

func (a ApplicationEvaluationController) evaluateApplication(c *gin.Context) {
	c.JSON(http.StatusOK, a.applicationEvaluationService.EvaluateCurrentApplication(c.Request.Context()))
}
