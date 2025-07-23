package v1

import (
	"net/http"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	appErr "github.com/compendium-tech/compendium/subscription-service/internal/error"
	"github.com/compendium-tech/compendium/subscription-service/internal/service"
	"github.com/gin-gonic/gin"
)

type SubscriptionController struct {
	subscriptionService service.SubscriptionService
}

func NewSubscriptionController(subscriptionService service.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{
		subscriptionService: subscriptionService,
	}
}

func (p *SubscriptionController) MakeRoutes(e *gin.Engine) {
	v1 := e.Group("/v1")
	{
		authenticated := v1.Group("")
		authenticated.Use(auth.RequireAuth)
		authenticated.GET("/subscription", appErr.HandleAppErr(p.getSubscription))

		authenticated.Use(auth.RequireCsrf)
		authenticated.DELETE("/subscription", appErr.HandleAppErr(p.cancelSubscription))
	}
}

func (p *SubscriptionController) getSubscription(c *gin.Context) error {
	subscription, err := p.subscriptionService.GetSubscriptionForAuthenticatedUser(c.Request.Context())
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, subscription)
	return nil
}

func (p *SubscriptionController) cancelSubscription(c *gin.Context) error {
	err := p.subscriptionService.CancelSubscriptionForAuthenticatedUser(c.Request.Context())
	if err != nil {
		return err
	}

	c.Status(http.StatusNoContent)
	return nil
}
