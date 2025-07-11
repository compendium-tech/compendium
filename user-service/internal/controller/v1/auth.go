package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seacite-tech/compendium/user-service/internal/domain"
	apperr "github.com/seacite-tech/compendium/user-service/internal/error"
	"github.com/seacite-tech/compendium/user-service/internal/service"
	"github.com/seacite-tech/compendium/user-service/internal/validate"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return AuthController{
		authService: authService,
	}
}

func (a AuthController) MakeRoutes(e *gin.Engine) {
	e.POST("/api/v1/users", apperr.HandleAppErr(a.signUp))
	e.POST("/api/v1/sessions", apperr.HandleAppErr(a.createSession))
}

func (a *AuthController) signUp(c *gin.Context) error {
	var request domain.SignUpRequest

	if err := c.BindJSON(&request); err != nil {
		return err
	}

	if err := validate.Validate.Struct(request); err != nil {
		return err
	}

	err := a.authService.SignUp(c.Request.Context(), request)

	if err != nil {
		return err
	}

	c.Status(http.StatusCreated)

	return nil
}

func (a *AuthController) createSession(c *gin.Context) error {
	flow := c.Query("flow")
	if flow != "direct" && flow != "mfa" {
		return apperr.Errorf(apperr.RequestValidationError, "Flow parameter must be equal to `mfa` or `direct`.")
	}

	if flow == "mfa" {
		var request domain.SubmitMfaOtpRequest

		if err := c.BindJSON(&request); err != nil {
			return err
		}

		if err := validate.Validate.Struct(request); err != nil {
			return err
		}

		sessionResponse, err := a.authService.SubmitMfaOtp(c.Request.Context(), request)
		if err != nil {
			return err
		}

		cookieExpiry := 30 * 365 * 24 * 3600
		c.SetCookie("csrfToken", sessionResponse.CsrfToken, cookieExpiry, "/", "", false, false)
		c.SetCookie("accessToken", sessionResponse.AccessToken, cookieExpiry, "/", "", false, true)
		c.SetCookie("refreshToken", sessionResponse.RefreshToken, cookieExpiry, "/", "", false, true)

		c.JSON(http.StatusCreated, sessionResponse.JsonResponse())
	}

	return nil
}
