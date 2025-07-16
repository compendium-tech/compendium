package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seacite-tech/compendium/common/pkg/httphelp"
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
	e.PUT("/api/v1/password", apperr.HandleAppErr(a.resetPassword))
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
	if flow != "password" && flow != "mfa" {
		return apperr.Errorf(apperr.RequestValidationError, "Flow parameter must be equal to `mfa` or `password`.")
	}

	if flow == "mfa" {
		var body domain.SubmitMfaOtpRequestBody

		if err := c.BindJSON(&body); err != nil {
			return err
		}

		if err := validate.Validate.Struct(body); err != nil {
			return err
		}

		request := domain.SubmitMfaOtpRequest{
			Email:     body.Email,
			Otp:       body.Otp,
			IpAddress: httphelp.GetClientIP(c),
			UserAgent: httphelp.GetUserAgent(c),
		}

		response, err := a.authService.SubmitMfaOtp(c.Request.Context(), request)
		if err != nil {
			return err
		}

		setAuthCookies(c, response.CsrfToken, response.AccessToken, response.RefreshToken)

		c.JSON(http.StatusCreated, response.IntoBody())
	} else {
		var body domain.SignInRequestBody

		if err := c.BindJSON(&body); err != nil {
			return err
		}

		if err := validate.Validate.Struct(body); err != nil {
			return err
		}

		request := domain.SignInRequest{
			Email:     body.Email,
			Password:  body.Password,
			IpAddress: httphelp.GetClientIP(c),
			UserAgent: httphelp.GetUserAgent(c),
		}

		response, err := a.authService.SignIn(c.Request.Context(), request)
		if err != nil {
			return err
		}

		if response.Session != nil {
			setAuthCookies(c, response.Session.CsrfToken, response.Session.AccessToken, response.Session.RefreshToken)
		}

		if !response.IsMfaRequired {
			c.JSON(http.StatusCreated, response.IntoBody())
		} else {
			c.JSON(http.StatusAccepted, response.IntoBody())
		}
	}

	return nil
}

func (a *AuthController) resetPassword(c *gin.Context) error {
	flow := c.Query("flow")
	if flow != "init" && flow != "finish" {
		return apperr.Errorf(apperr.RequestValidationError, "Flow parameter must be equal to `init` or `finish`.")
	}

	if flow == "init" {
		var request domain.InitPasswordResetRequest

		if err := c.BindJSON(&request); err != nil {
			return err
		}

		if err := validate.Validate.Struct(request); err != nil {
			return err
		}

		err := a.authService.InitPasswordReset(c.Request.Context(), request)
		if err != nil {
			return err
		}

		c.Status(http.StatusAccepted)
	} else {
		var request domain.FinishPasswordResetRequest

		if err := c.BindJSON(&request); err != nil {
			return err
		}

		if err := validate.Validate.Struct(request); err != nil {
			return err
		}

		err := a.authService.FinishPasswordReset(c.Request.Context(), request)
		if err != nil {
			return err
		}

		c.Status(http.StatusOK)
	}

	return nil
}

func setAuthCookies(c *gin.Context, csrfToken, accessToken, refreshToken string) {
	cookieExpiry := 30 * 365 * 24 * 3600

	c.SetCookie("csrfToken", csrfToken, cookieExpiry, "/", "", false, false)
	c.SetCookie("accessToken", accessToken, cookieExpiry, "/", "", false, true)
	c.SetCookie("refreshToken", refreshToken, cookieExpiry, "/", "", false, true)
}
