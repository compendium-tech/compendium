package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seacite-tech/compendium/common/pkg/httphelp"
	"github.com/seacite-tech/compendium/common/pkg/validate"
	"github.com/seacite-tech/compendium/user-service/internal/domain"
	appErr "github.com/seacite-tech/compendium/user-service/internal/error"
	"github.com/seacite-tech/compendium/user-service/internal/service"
)

const (
	csrfTokenCookieName    = "csrfToken"
	accessTokenCookieName  = "accessToken"
	refreshTokenCookieName = "refreshToken"
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
	v1 := e.Group("/api/v1/")
	{
		v1.POST("/users", appErr.HandleAppErr(a.signUp))
		v1.POST("/sessions", appErr.HandleAppErr(a.createSession))
		v1.PUT("/password", appErr.HandleAppErr(a.resetPassword))
		v1.DELETE("/session", appErr.HandleAppErr(a.logout))
	}
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
	switch c.Query("flow") {
	case "mfa":
		return a.submitMfaOtp(c)
	case "password":
		return a.signIn(c)
	case "refresh":
		return a.refresh(c)
	default:
		return appErr.Errorf(appErr.RequestValidationError, "Flow parameter must be equal to `mfa`, `password` or `refresh`.")
	}
}

func (a *AuthController) submitMfaOtp(c *gin.Context) error {
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
	return nil
}

func (a *AuthController) signIn(c *gin.Context) error {
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

	return nil
}

func (a *AuthController) refresh(c *gin.Context) error {
	refreshTokenCookie, err := c.Request.Cookie(refreshTokenCookieName)
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
	return nil
}

func (a *AuthController) resetPassword(c *gin.Context) error {
	switch c.Query("flow") {
	case "init":
		return a.initPasswordReset(c)
	case "finish":
		return a.finishPasswordReset(c)
	default:
		return appErr.Errorf(appErr.RequestValidationError, "Flow parameter must be equal to `init` or `finish`.")
	}
}

func (a *AuthController) initPasswordReset(c *gin.Context) error {
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
	return nil
}

func (a *AuthController) finishPasswordReset(c *gin.Context) error {
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
	return nil
}

func (a *AuthController) logout(c *gin.Context) error {
	refreshTokenCookie, err := c.Request.Cookie(refreshTokenCookieName)
	if err != nil {
		return appErr.Errorf(appErr.InvalidSessionError, "Invalid session")
	}

	err = a.authService.Logout(c.Request.Context(), refreshTokenCookie.Value)
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	return nil
}

func setAuthCookies(c *gin.Context, session *domain.SessionResponse) {
	cookieExpiry := 30 * 365 * 24 * 3600

	c.SetCookie(csrfTokenCookieName, session.CsrfToken, cookieExpiry, "/", "", false, false)
	c.SetCookie(accessTokenCookieName, session.AccessToken, cookieExpiry, "/", "", false, true)
	c.SetCookie(refreshTokenCookieName, session.RefreshToken, cookieExpiry, "/", "", false, true)
}
