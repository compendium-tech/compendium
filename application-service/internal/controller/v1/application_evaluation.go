package v1

import (
	"github.com/compendium-tech/compendium/application-service/internal/middleware"
	"github.com/compendium-tech/compendium/application-service/internal/service"
	"github.com/compendium-tech/compendium/common/pkg/auth"
	httputils "github.com/compendium-tech/compendium/common/pkg/http"
	"github.com/gin-gonic/gin"
	"net/http"
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

func (a ApplicationEvaluationController) evaluateApplication(c *gin.Context) error {
	evaluation, err := a.applicationEvaluationService.EvaluateCurrentApplication(c.Request.Context())
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, evaluation)
	return nil
}
