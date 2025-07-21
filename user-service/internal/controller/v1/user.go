package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seacite-tech/compendium/common/pkg/validate"
	"github.com/seacite-tech/compendium/user-service/internal/domain"
	appErr "github.com/seacite-tech/compendium/user-service/internal/error"
	"github.com/seacite-tech/compendium/user-service/internal/service"
	"github.com/seacite-tech/compendium/user-service/pkg/auth"
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
	v1 := e.Group("/api/v1/")
	{
		authenticated := v1.Group("/")
		authenticated.Use(auth.RequireAuth)
		authenticated.GET("/account", appErr.HandleAppErr(u.getAccount))

		authenticated.Use(auth.RequireCsrf)
		authenticated.PUT("/account", appErr.HandleAppErr(u.updateAccount))
	}
}

func (u *UserController) getAccount(c *gin.Context) error {
	response, err := u.userService.GetAccount(c.Request.Context())

	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response)
	return nil
}

func (u *UserController) updateAccount(c *gin.Context) error {
	var request domain.UpdateAccount

	if err := c.BindJSON(&request); err != nil {
		return err
	}

	if err := validate.Validate.Struct(request); err != nil {
		return err
	}

	response, err := u.userService.UpdateAccount(c.Request.Context(), request)

	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response)
	return nil
}
