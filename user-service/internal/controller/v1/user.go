package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appErr "github.com/seacite-tech/compendium/user-service/internal/error"
	"github.com/seacite-tech/compendium/user-service/internal/service"
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
	e.GET("/api/v1/account", appErr.HandleAppErr(u.getAccount))
}

func (u *UserController) getAccount(c *gin.Context) error {
	response, err := u.userService.GetAccount(c.Request.Context())

	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, response)
	return nil
}
