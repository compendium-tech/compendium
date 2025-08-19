package httpv1

import (
	"net/http"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	httputils "github.com/compendium-tech/compendium/common/pkg/http"
	"github.com/compendium-tech/compendium/user-service/internal/domain"
	"github.com/compendium-tech/compendium/user-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserController struct {
	userService service.UserService
	e           httputils.ErrorHandler
}

func NewUserController(userService service.UserService) UserController {
	return UserController{
		userService: userService,
	}
}

func (u UserController) MakeRoutes(e *gin.Engine) {
	v1 := e.Group("/v1/")
	{
		authenticated := v1.Group("/")
		authenticated.Use(auth.RequireAuth)
		authenticated.GET("/account", u.e.Handle(u.getAccount))

		authenticated.Use(auth.RequireCsrf)
		authenticated.PUT("/account", u.e.Handle(u.updateAccount))
	}
}

func (u UserController) getAccount(c *gin.Context) {
	c.JSON(http.StatusOK, u.userService.GetAccountAsAuthenticatedUser(c.Request.Context()))
}

func (u UserController) updateAccount(c *gin.Context) {
	c.JSON(http.StatusOK,
		u.userService.UpdateAccountAsAuthenticatedUser(c.Request.Context(),
			httputils.MustBindWith[domain.UpdateAccount](c, binding.JSON).Validated()))
}
