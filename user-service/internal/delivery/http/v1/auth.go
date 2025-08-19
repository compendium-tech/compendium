package httpv1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	httputils "github.com/compendium-tech/compendium/common/pkg/http"

	"github.com/compendium-tech/compendium/user-service/internal/domain"
	myerror "github.com/compendium-tech/compendium/user-service/internal/error"
	"github.com/compendium-tech/compendium/user-service/internal/service"
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
	var eh httputils.ErrorHandler

	v1 := e.Group("/v1/")
	{
		v1.POST("/users", eh.Handle(a.signUp))
		v1.POST("/sessions", eh.Handle(a.createSession))
		v1.PUT("/password", eh.Handle(a.resetPassword))
		v1.DELETE("/session", eh.Handle(a.logout))

		authenticated := v1.Group("/")
		{
			authenticated.Use(auth.RequireAuth)
			authenticated.GET("/sessions", eh.Handle(a.getSessions))
			authenticated.DELETE("/sessions/:id", eh.Handle(a.removeSession))
		}
	}
}

func (a AuthController) signUp(c *gin.Context) {
	a.authService.SignUp(c.Request.Context(),
		httputils.MustBindWith[domain.SignUpRequest](c, binding.JSON).Validated())
	c.Status(http.StatusCreated)
}

func (a AuthController) createSession(c *gin.Context) {
	switch c.Query("flow") {
	case "mfa":
		body := httputils.MustBindWith[domain.SubmitMfaOtpRequestBody](c, binding.JSON).Validated()

		setSessionCreatedResponse(c,
			a.authService.SubmitMfaOtp(c.Request.Context(),
				domain.SubmitMfaOtpRequest{
					Email:     body.Email,
					Otp:       body.Otp,
					IPAddress: httputils.GetClientIP(c),
					UserAgent: httputils.GetUserAgent(c),
				}))
	case "password":
		body := httputils.MustBindWith[domain.SignInRequestBody](c, binding.JSON).Validated()

		response := a.authService.SignIn(c.Request.Context(), domain.SignInRequest{
			Email:     body.Email,
			Password:  body.Password,
			IPAddress: httputils.GetClientIP(c),
			UserAgent: httputils.GetUserAgent(c),
		})

		if response.Session != nil {
			setAuthCookies(c, *response.Session)
		}

		if !response.IsMfaRequired {
			c.JSON(http.StatusCreated, response.IntoBody())
		} else {
			c.JSON(http.StatusAccepted, response.IntoBody())
		}
	case "refresh":
		refreshTokenCookie, err := c.Request.Cookie(refreshTokenCookieName)
		if err != nil {
			myerror.New(myerror.InvalidSessionError).Throw()
		}

		setSessionCreatedResponse(c,
			a.authService.Refresh(c.Request.Context(), domain.RefreshTokenRequest{
				RefreshToken: refreshTokenCookie.Value,
				IPAddress:    httputils.GetClientIP(c),
				UserAgent:    httputils.GetUserAgent(c),
			}))
	default:
		myerror.NewWithReason(myerror.RequestValidationError, "Flow parameter must be equal to `mfa`, `password` or `refresh`.").Throw()
	}
}

func (a AuthController) resetPassword(c *gin.Context) {
	switch c.Query("flow") {
	case "init":
		a.authService.InitPasswordReset(c.Request.Context(),
			httputils.MustBindWith[domain.InitPasswordResetRequest](c, binding.JSON).Validated())
		c.Status(http.StatusAccepted)
	case "finish":
		a.authService.FinishPasswordReset(c.Request.Context(),
			httputils.MustBindWith[domain.FinishPasswordResetRequest](c, binding.JSON).Validated())
		c.Status(http.StatusOK)
	default:
		myerror.NewWithReason(myerror.RequestValidationError, "Flow parameter must be equal to `init` or `finish`.").Throw()
	}
}

func (a AuthController) logout(c *gin.Context) {
	refreshTokenCookie, err := c.Request.Cookie(refreshTokenCookieName)
	if err != nil {
		myerror.New(myerror.InvalidSessionError).Throw()
	}

	a.authService.Logout(c.Request.Context(), refreshTokenCookie.Value)
	removeAuthCookies(c)

	c.Status(http.StatusOK)
}

func (a AuthController) getSessions(c *gin.Context) {
	refreshTokenCookie, err := c.Request.Cookie(refreshTokenCookieName)
	if err != nil {
		myerror.New(myerror.InvalidSessionError).Throw()
	}

	response := a.authService.GetSessionsForAuthenticatedUser(c.Request.Context(), refreshTokenCookie.Value)
	c.JSON(http.StatusOK, response)
}

func (a AuthController) removeSession(c *gin.Context) {
	sessionIdString := c.Param("id")
	if sessionIdString == "" {
		myerror.New(myerror.RequestValidationError).Throw()
	}

	sessionId, err := uuid.Parse(sessionIdString)
	if err != nil {
		myerror.New(myerror.RequestValidationError).Throw()
	}

	refreshTokenCookie, err := c.Request.Cookie(refreshTokenCookieName)
	if err != nil {
		myerror.New(myerror.InvalidSessionError).Throw()
	}

	a.authService.RemoveSessionByID(c.Request.Context(), sessionId, refreshTokenCookie.Value)
	c.Status(http.StatusNoContent)
}

func setAuthCookies(c *gin.Context, session domain.SessionResponse) {
	cookieExpiresAt := 30 * 365 * 24 * 3600

	c.SetCookie(csrfTokenCookieName, session.CsrfToken, cookieExpiresAt, "/", "", false, false)
	c.SetCookie(accessTokenCookieName, session.AccessToken, cookieExpiresAt, "/", "", false, true)
	c.SetCookie(refreshTokenCookieName, session.RefreshToken, cookieExpiresAt, "/", "", false, true)
}

func setSessionCreatedResponse(c *gin.Context, session domain.SessionResponse) {
	setAuthCookies(c, session)
	c.JSON(http.StatusCreated, session.IntoBody())
}

func removeAuthCookies(c *gin.Context) {
	c.SetCookie(csrfTokenCookieName, "", 0, "/", "", false, false)
	c.SetCookie(accessTokenCookieName, "", 0, "/", "", false, true)
	c.SetCookie(refreshTokenCookieName, "", 0, "/", "", false, true)
}
