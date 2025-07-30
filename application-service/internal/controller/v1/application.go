package v1

import (
	"net/http"

	"github.com/compendium-tech/compendium/application-service/internal/domain"
	"github.com/compendium-tech/compendium/application-service/internal/middleware"
	"github.com/compendium-tech/compendium/application-service/internal/service"
	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/http"
	"github.com/compendium-tech/compendium/common/pkg/validate"
	"github.com/gin-gonic/gin"
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

				application.GET("/supplemental-essays", eh.Handle(a.getSupplementalEssays))
				application.PUT("/supplemental-essays", auth.RequireCsrf, eh.Handle(a.putSupplementalEssays))
			}
		}
	}
}

func (a ApplicationController) getApplications(c *gin.Context) error {
	response, err := a.applicationService.GetApplications(c.Request.Context())
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response)
	return nil
}

func (a ApplicationController) createApplication(c *gin.Context) error {
	var request domain.CreateApplicationRequest

	if err := c.BindJSON(&request); err != nil {
		return err
	}

	if err := validate.Validate.Struct(request); err != nil {
		return err
	}

	response, err := a.applicationService.CreateApplication(c.Request.Context(), request)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, response)
	return nil
}

func (a ApplicationController) updateApplicationName(c *gin.Context) error {
	var request struct {
		Name string `json:"name"`
	}

	if err := c.BindJSON(&request); err != nil {
		return err
	}

	err := a.applicationService.UpdateCurrentApplicationName(c.Request.Context(), request.Name)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}

func (a ApplicationController) removeApplication(c *gin.Context) error {
	err := a.applicationService.RemoveCurrentApplication(c.Request.Context())
	if err != nil {
		return err
	}

	c.Status(http.StatusNoContent)
	return nil
}

func (a ApplicationController) getActivities(c *gin.Context) error {
	response, err := a.applicationService.GetActivities(c.Request.Context())
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response)
	return nil
}

func (a ApplicationController) putActivities(c *gin.Context) error {
	var request []domain.UpdateActivityRequest

	if err := c.BindJSON(&request); err != nil {
		return err
	}

	if err := validate.Validate.Struct(request); err != nil {
		return err
	}

	err := a.applicationService.PutActivities(c.Request.Context(), request)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}

func (a ApplicationController) getHonors(c *gin.Context) error {
	response, err := a.applicationService.GetHonors(c.Request.Context())
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response)
	return nil
}

func (a ApplicationController) putHonors(c *gin.Context) error {
	var request []domain.UpdateHonorRequest

	if err := c.BindJSON(&request); err != nil {
		return err
	}

	if err := validate.Validate.Struct(request); err != nil {
		return err
	}

	err := a.applicationService.PutHonors(c.Request.Context(), request)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}

func (a ApplicationController) getEssays(c *gin.Context) error {
	response, err := a.applicationService.GetEssays(c.Request.Context())
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response)
	return nil
}

func (a ApplicationController) putEssays(c *gin.Context) error {
	var request []domain.UpdateEssayRequest

	if err := c.BindJSON(&request); err != nil {
		return err
	}

	if err := validate.Validate.Struct(request); err != nil {
		return err
	}

	err := a.applicationService.PutEssays(c.Request.Context(), request)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}

func (a ApplicationController) getSupplementalEssays(c *gin.Context) error {
	response, err := a.applicationService.GetSupplementalEssays(c.Request.Context())
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response)
	return nil
}

func (a ApplicationController) putSupplementalEssays(c *gin.Context) error {
	var request []domain.UpdateSupplementalEssayRequest

	if err := c.BindJSON(&request); err != nil {
		return err
	}

	if err := validate.Validate.Struct(request); err != nil {
		return err
	}

	err := a.applicationService.PutSupplementalEssays(c.Request.Context(), request)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}
