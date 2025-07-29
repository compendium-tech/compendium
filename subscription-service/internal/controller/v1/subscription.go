package v1

import (
	"fmt"
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
		authenticated.GET("/subscription", appErr.Handle(p.getSubscription))
		authenticated.GET("/subscription/invitationCode", appErr.Handle(p.getSubscriptionInvitationCode))

		authenticated.Use(auth.RequireCsrf)
		authenticated.DELETE("/subscription", appErr.Handle(p.cancelSubscription))
		authenticated.DELETE("/subscription/members/:id", appErr.Handle(p.removeSubscriptionMember))
		authenticated.POST("/subscription/members/me", appErr.Handle(p.joinSubscription))
		authenticated.PUT("/subscription/invitationCode", appErr.Handle(p.updateSubscriptionInvitationCode))
		authenticated.DELETE("/subscription/invitationCode", appErr.Handle(p.removeSubscriptionInvitationCode))
	}
}

func (p *SubscriptionController) getSubscription(c *gin.Context) error {
	subscription, err := p.subscriptionService.GetSubscription(c.Request.Context())
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, subscription)
	return nil
}

func (p *SubscriptionController) getSubscriptionInvitationCode(c *gin.Context) error {
	code, err := p.subscriptionService.GetSubscriptionInvitationCode(c.Request.Context())
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, code)
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

func (p *SubscriptionController) joinSubscription(c *gin.Context) error {
	invitationCode := c.Query("invitationCode")
	if invitationCode == "" {
		return appErr.NewWithDetails(appErr.RequestValidationError, "member ID is required")
	}

	subscription, err := p.subscriptionService.JoinCollectiveSubscription(c.Request.Context(), invitationCode)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, subscription)
	return nil
}

func (p *SubscriptionController) removeSubscriptionMember(c *gin.Context) error {
	memberIDString := c.Param("id")
	if memberIDString == "" {
		return appErr.NewWithReason(appErr.RequestValidationError, "member ID is required")
	}

	memberID, err := uuid.Parse(memberIDString)
	if err != nil {
		return appErr.NewWithReason(appErr.RequestValidationError, fmt.Sprintf("invalid member ID format: %w", err))
	}

	err = p.subscriptionService.RemoveSubscriptionMember(c.Request.Context(), memberID)
	if err != nil {
		return err
	}

	c.Status(http.StatusNoContent)
	return nil
}

func (p *SubscriptionController) updateSubscriptionInvitationCode(c *gin.Context) error {
	code, err := p.subscriptionService.UpdateSubscriptionInvitationCode(c.Request.Context())
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, code)
	return nil
}

func (p *SubscriptionController) removeSubscriptionInvitationCode(c *gin.Context) error {
	code, err := p.subscriptionService.RemoveSubscriptionInvitationCode(c.Request.Context())
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, code)
	return nil
}
