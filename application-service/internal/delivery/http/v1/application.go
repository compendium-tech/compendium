package httpv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	httputils "github.com/compendium-tech/compendium/common/pkg/http"

	"github.com/compendium-tech/compendium/application-service/internal/domain"
	"github.com/compendium-tech/compendium/application-service/internal/middleware"
	"github.com/compendium-tech/compendium/application-service/internal/service"
)

type ApplicationController struct {
	applicationService service.ApplicationService
}

func NewApplicationController(applicationService service.ApplicationService) ApplicationController {
	return ApplicationController{
		applicationService: applicationService,
	}
}

func (a ApplicationController) MakeRoutes(e *gin.Engine) {
	var eh httputils.ErrorHandler

	v1 := e.Group("/v1")
	{
		authenticated := v1.Group("/")
		authenticated.Use(auth.RequireAuth)
		{
			authenticated.GET("/applications", eh.Handle(a.getApplications))
			authenticated.POST("/applications", eh.Handle(a.createApplication))

			application := authenticated.Group("/applications/:applicationId")
			application.Use(middleware.NewSetApplicationFromRequest(a.applicationService).Handle)
			{
				application.PUT("/", auth.RequireCsrf, eh.Handle(a.updateApplicationName))
				application.DELETE("/", auth.RequireCsrf, eh.Handle(a.removeApplication))

				application.GET("/activities", eh.Handle(a.getActivities))
				application.PUT("/activities", auth.RequireCsrf, eh.Handle(a.putActivities))

				application.GET("/honors", eh.Handle(a.getHonors))
				application.PUT("/honors", auth.RequireCsrf, eh.Handle(a.putHonors))

				application.GET("/essays", eh.Handle(a.getEssays))
				application.PUT("/essays", auth.RequireCsrf, eh.Handle(a.putEssays))

				application.GET("/supplementalEssays", eh.Handle(a.getSupplementalEssays))
				application.PUT("/supplementalEssays", auth.RequireCsrf, eh.Handle(a.putSupplementalEssays))
			}
		}
	}
}

func (a ApplicationController) getApplications(c *gin.Context) {
	c.JSON(http.StatusOK, a.applicationService.GetApplications(c.Request.Context()))
}

func (a ApplicationController) createApplication(c *gin.Context) {
	c.JSON(http.StatusCreated, a.applicationService.CreateApplication(
		c.Request.Context(),
		httputils.MustBindWith[domain.CreateApplicationRequest](c, binding.JSON).Validated()))
}

func (a ApplicationController) updateApplicationName(c *gin.Context) {
	a.applicationService.UpdateCurrentApplicationName(
		c.Request.Context(),
		httputils.MustBindWith[struct {
			Name string `json:"name" validate:"required,min=1,max=100"`
		}](c, binding.JSON).Validated().Name)
	c.Status(http.StatusOK)
}

func (a ApplicationController) removeApplication(c *gin.Context) {
	a.applicationService.RemoveCurrentApplication(c.Request.Context())
	c.Status(http.StatusNoContent)
}

func (a ApplicationController) getActivities(c *gin.Context) {
	c.JSON(http.StatusOK, a.applicationService.GetActivities(c.Request.Context()))
}

func (a ApplicationController) putActivities(c *gin.Context) {
	a.applicationService.PutActivities(
		c.Request.Context(),
		httputils.MustBindWith[[]domain.UpdateActivityRequest](c, binding.JSON).Validated())
	c.Status(http.StatusOK)
}

func (a ApplicationController) getHonors(c *gin.Context) {
	c.JSON(http.StatusOK, a.applicationService.GetHonors(c.Request.Context()))
}

func (a ApplicationController) putHonors(c *gin.Context) {
	a.applicationService.PutHonors(
		c.Request.Context(),
		httputils.MustBindWith[[]domain.UpdateHonorRequest](c, binding.JSON).Validated())
	c.Status(http.StatusOK)
}

func (a ApplicationController) getEssays(c *gin.Context) {
	c.JSON(http.StatusOK, a.applicationService.GetEssays(c.Request.Context()))
}

func (a ApplicationController) putEssays(c *gin.Context) {
	a.applicationService.PutEssays(c.Request.Context(),
		httputils.MustBindWith[[]domain.UpdateEssayRequest](c, binding.JSON).Validated())
	c.Status(http.StatusOK)
}

func (a ApplicationController) getSupplementalEssays(c *gin.Context) {
	c.JSON(http.StatusOK, a.applicationService.GetSupplementalEssays(c.Request.Context()))
}

func (a ApplicationController) putSupplementalEssays(c *gin.Context) {
	a.applicationService.PutSupplementalEssays(
		c.Request.Context(),
		httputils.MustBindWith[[]domain.UpdateSupplementalEssayRequest](c, binding.JSON).Validated())
	c.Status(http.StatusOK)
}
