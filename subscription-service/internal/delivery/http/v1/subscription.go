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

func (p *SubscriptionController) getSubscription(c *gin.Context) {
	c.JSON(http.StatusOK, p.subscriptionService.GetSubscription(c.Request.Context()))
}

func (p *SubscriptionController) getSubscriptionInvitationCode(c *gin.Context) {
	c.JSON(http.StatusOK, p.subscriptionService.GetSubscriptionInvitationCode(c.Request.Context()))
}

func (p *SubscriptionController) cancelSubscription(c *gin.Context) {
	p.subscriptionService.CancelSubscription(c.Request.Context())
	c.Status(http.StatusNoContent)
}

func (p *SubscriptionController) joinSubscription(c *gin.Context) {
	invitationCode := c.Query("invitationCode")
	if invitationCode == "" {
		myerror.NewWithDetails(myerror.RequestValidationError, "member ID is required").Throw()
	}

	c.JSON(http.StatusOK, p.subscriptionService.JoinCollectiveSubscription(c.Request.Context(), invitationCode))
}

func (p *SubscriptionController) removeSubscriptionMember(c *gin.Context) {
	memberIDString := c.Param("id")
	if memberIDString == "" {
		myerror.NewWithReason(myerror.RequestValidationError, "member ID is required").Throw()
	}

	memberID, err := uuid.Parse(memberIDString)
	if err != nil {
		myerror.NewWithReason(myerror.RequestValidationError, fmt.Sprintf("invalid member ID format: %v", err)).Throw()
	}

	p.subscriptionService.RemoveSubscriptionMember(c.Request.Context(), memberID)
	c.Status(http.StatusNoContent)
}

func (p *SubscriptionController) updateSubscriptionInvitationCode(c *gin.Context) {
	c.JSON(http.StatusOK, p.subscriptionService.UpdateSubscriptionInvitationCode(c.Request.Context()))
}

func (p *SubscriptionController) removeSubscriptionInvitationCode(c *gin.Context) {
	c.JSON(http.StatusOK, p.subscriptionService.RemoveSubscriptionInvitationCode(c.Request.Context()))
}
