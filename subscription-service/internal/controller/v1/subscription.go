package v1

import (
	"net/http"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	appErr "github.com/compendium-tech/compendium/subscription-service/internal/error"
	"github.com/compendium-tech/compendium/subscription-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		authenticated.DELETE("/subscription/members/:id", appErr.HandleAppErr(p.removeSubscriptionMember))
	}
}

func (p *SubscriptionController) getSubscription(c *gin.Context) error {
	subscription, err := p.subscriptionService.GetSubscription(c.Request.Context())
	if err != nil {
		return err
	}

	c.JSON(http.StatusCreated, subscription)
	return nil
}

func (p *SubscriptionController) cancelSubscription(c *gin.Context) error {
	err := p.subscriptionService.CancelSubscription(c.Request.Context())
	if err != nil {
		return err
	}

	c.Status(http.StatusNoContent)
	return nil
}

func (p *SubscriptionController) removeSubscriptionMember(c *gin.Context) error {
	memberIDString := c.Param("id")
	if memberIDString == "" {
		return appErr.Errorf(appErr.RequestValidationError, "member ID is required")
	}

	memberID, err := uuid.Parse(memberIDString)
	if err != nil {
		return appErr.Errorf(appErr.RequestValidationError, "invalid member ID format: %v", err)
	}

	err = p.subscriptionService.RemoveSubscriptionMember(c.Request.Context(), memberID)
	if err != nil {
		return err
	}

	c.Status(http.StatusNoContent)
	return nil
}
