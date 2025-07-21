package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/adslmgrv/compendium/subscription-service/internal/repository"
	"github.com/gin-gonic/gin"
)

type PaddleWebhookController struct {
	repository repository.SubscriptionRepository
	verifier   paddle.WebhookVerifier
}

func NewPaddleWebhookController(
	repository repository.SubscriptionRepository,
	verifier paddle.WebhookVerifier) *PaddleWebhookController {
	return &PaddleWebhookController{
		repository: repository,
		verifier:   verifier,
	}
}

func (p *PaddleWebhookController) MakeRoutes(e *gin.Engine) {
	e.POST("/paddleWebhook", p.Handle)
}

func (p *PaddleWebhookController) Handle(c *gin.Context) {
	ok, err := p.verifier.Verify(c.Request)

	if err != nil && (errors.Is(err, paddle.ErrMissingSignature) || errors.Is(err, paddle.ErrInvalidSignatureFormat)) {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if !ok {
		c.AbortWithError(http.StatusForbidden, fmt.Errorf("signature mismatch"))
		return
	}
}
