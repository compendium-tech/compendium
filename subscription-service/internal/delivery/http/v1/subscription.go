package httpv1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/compendium-tech/compendium/common/pkg/auth"
	httputils "github.com/compendium-tech/compendium/common/pkg/http"
	myerror "github.com/compendium-tech/compendium/subscription-service/internal/error"
	"github.com/compendium-tech/compendium/subscription-service/internal/service"
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
	var eh httputils.ErrorHandler

	v1 := e.Group("/v1")
	{
		authenticated := v1.Group("")
		authenticated.Use(auth.RequireAuth)
		authenticated.GET("/subscription", eh.Handle(p.getSubscription))
		authenticated.GET("/subscription/invitationCode", eh.Handle(p.getSubscriptionInvitationCode))

		authenticated.Use(auth.RequireCsrf)
		authenticated.DELETE("/subscription", eh.Handle(p.cancelSubscription))
		authenticated.DELETE("/subscription/members/:id", eh.Handle(p.removeSubscriptionMember))
		authenticated.POST("/subscription/members/me", eh.Handle(p.joinSubscription))
		authenticated.PUT("/subscription/invitationCode", eh.Handle(p.updateSubscriptionInvitationCode))
		authenticated.DELETE("/subscription/invitationCode", eh.Handle(p.removeSubscriptionInvitationCode))
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
		return myerror.NewWithDetails(myerror.RequestValidationError, "member ID is required")
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
		return myerror.NewWithReason(myerror.RequestValidationError, "member ID is required")
	}

	memberID, err := uuid.Parse(memberIDString)
	if err != nil {
		return myerror.NewWithReason(myerror.RequestValidationError, fmt.Sprintf("invalid member ID format: %v", err))
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
