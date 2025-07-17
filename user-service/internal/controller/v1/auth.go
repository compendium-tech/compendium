package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seacite-tech/compendium/common/pkg/httphelp"
	"github.com/seacite-tech/compendium/user-service/internal/domain"
	appErr "github.com/seacite-tech/compendium/user-service/internal/error"
	"github.com/seacite-tech/compendium/user-service/internal/service"
	"github.com/seacite-tech/compendium/user-service/internal/validate"
	"github.com/seacite-tech/compendium/user-service/pkg/auth"
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
	e.POST("/api/v1/users", appErr.HandleAppErr(a.signUp))
	e.POST("/api/v1/sessions", appErr.HandleAppErr(a.createSession))
	e.PUT("/api/v1/password", appErr.HandleAppErr(a.resetPassword))
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
	if flow != "password" && flow != "mfa" && flow != "refresh" {
		return appErr.Errorf(appErr.RequestValidationError, "Flow parameter must be equal to `mfa`, `password` or `refresh`.")
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

		setAuthCookies(c, response)

		c.JSON(http.StatusCreated, response.IntoBody())
	} else if flow == "password" {
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
			setAuthCookies(c, response.Session)
		}

		if !response.IsMfaRequired {
			c.JSON(http.StatusCreated, response.IntoBody())
		} else {
			c.JSON(http.StatusAccepted, response.IntoBody())
		}
	} else {
		auth.RequireAuth(c)

		refreshTokenCookie, err := c.Request.Cookie("refreshToken")
		if err != nil {
			return appErr.Errorf(appErr.InvalidSessionError, "Invalid session")
		}

		response, err := a.authService.Refresh(c.Request.Context(), domain.RefreshTokenRequest{
			RefreshToken: refreshTokenCookie.Value,
			IpAddress:    httphelp.GetClientIP(c),
			UserAgent:    httphelp.GetUserAgent(c),
		})
		if err != nil {
			return err
		}

		setAuthCookies(c, response)

		c.JSON(http.StatusCreated, response.IntoBody())
	}

	return nil
}

func (a *AuthController) resetPassword(c *gin.Context) error {
	flow := c.Query("flow")
	if flow != "init" && flow != "finish" {
		return appErr.Errorf(appErr.RequestValidationError, "Flow parameter must be equal to `init` or `finish`.")
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

func setAuthCookies(c *gin.Context, session *domain.SessionResponse) {
	cookieExpiry := 30 * 365 * 24 * 3600

	c.SetCookie("csrfToken", session.CsrfToken, cookieExpiry, "/", "", false, false)
	c.SetCookie("accessToken", session.AccessToken, cookieExpiry, "/", "", false, true)
	c.SetCookie("refreshToken", session.RefreshToken, cookieExpiry, "/", "", false, true)
}
