package httpv1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	httputils "github.com/compendium-tech/compendium/common/pkg/http"
	"github.com/compendium-tech/compendium/common/pkg/validate"

	"github.com/compendium-tech/compendium/college-service/internal/domain"
	"github.com/compendium-tech/compendium/college-service/internal/service"
)

type CollegeController struct {
	collegeService service.CollegeService
}

func NewCollegeController(collegeService service.CollegeService) CollegeController {
	return CollegeController{
		collegeService: collegeService,
	}
}

func (c CollegeController) MakeRoutes(e *gin.Engine) {
	var eh httputils.ErrorHandler

	v1 := e.Group("/v1")
	{
		authenticated := v1.Group("/")
		authenticated.Use(auth.RequireAuth)
		authenticated.GET("/colleges", eh.Handle(c.searchColleges))
	}
}

func (cc CollegeController) searchColleges(c *gin.Context) error {
	var request domain.SearchCollegesRequest

	if err := c.BindJSON(&request); err != nil {
		return err
	}

	if err := validate.Validate.Struct(request); err != nil {
		return err
	}

	response, err := cc.collegeService.SearchColleges(c.Request.Context(), request)
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, response)
	return nil
}
