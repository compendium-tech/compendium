package httpv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	httputils "github.com/compendium-tech/compendium/common/pkg/http"

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

func (cc CollegeController) searchColleges(c *gin.Context) {
	c.JSON(http.StatusCreated,
		cc.collegeService.SearchColleges(c.Request.Context(),
			httputils.MustBindWith[domain.SearchCollegesRequest](c, binding.JSON, true)))
}
