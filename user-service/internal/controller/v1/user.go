package v1

import (
	"github.com/compendium-tech/compendium/common/pkg/error"
	"net/http"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	"github.com/compendium-tech/compendium/common/pkg/validate"
	"github.com/compendium-tech/compendium/user-service/internal/domain"
	"github.com/compendium-tech/compendium/user-service/internal/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.UserService
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
		authenticated.GET("/account", errorutils.Handle(u.getAccount))

		authenticated.Use(auth.RequireCsrf)
		authenticated.PUT("/account", errorutils.Handle(u.updateAccount))
	}
}

func (u UserController) getAccount(c *gin.Context) error {
	response, err := u.userService.GetAccountAsAuthenticatedUser(c.Request.Context())

	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response)
	return nil
}

func (u UserController) updateAccount(c *gin.Context) error {
	var request domain.UpdateAccount

	if err := c.BindJSON(&request); err != nil {
		return err
	}

	if err := validate.Validate.Struct(request); err != nil {
		return err
	}

	response, err := u.userService.UpdateAccountAsAuthenticatedUser(c.Request.Context(), request)

	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response)
	return nil
}
