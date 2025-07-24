package v1

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/PaddleHQ/paddle-go-sdk/v4/pkg/paddlenotification"
	"github.com/compendium-tech/compendium/subscription-service/internal/domain"
	appErr "github.com/compendium-tech/compendium/subscription-service/internal/error"
	"github.com/compendium-tech/compendium/subscription-service/internal/service"
	"github.com/compendium-tech/compendium/subscription-service/internal/webhook"
	"github.com/gin-gonic/gin"
	"github.com/ztrue/tracerr"
)

const dateTimeLayout = time.RFC3339

type BillingWebhookController struct {
	subscriptionService service.SubscriptionService
	webhookVerifier     webhook.WebhookVerifier
}

func NewBillingWebhookController(
	subscriptionService service.SubscriptionService,
	webhookVerifier webhook.WebhookVerifier) *BillingWebhookController {
	return &BillingWebhookController{
		subscriptionService: subscriptionService,
		webhookVerifier:     webhookVerifier,
	}
}

func (p *BillingWebhookController) MakeRoutes(e *gin.Engine) {
	e.POST("/v1/billingEvents", appErr.Handle(p.handle))
}

type event struct {
	EventID   string                           `json:"event_id"`
	EventType paddlenotification.EventTypeName `json:"event_type"`
}

func (p *BillingWebhookController) handle(c *gin.Context) error {
	ok, err := p.webhookVerifier.Verify(c.Request)

	if err != nil && (errors.Is(err, paddle.ErrMissingSignature) || errors.Is(err, paddle.ErrInvalidSignatureFormat)) {
		return appErr.Errorf(appErr.InvalidWebhookSignatureError, "Failed to verify request signature")
	}

	if err != nil {
		return tracerr.Wrap(err)
	}

	if !ok {
		return appErr.Errorf(appErr.InvalidWebhookSignatureError, "Failed to verify request signature")
	}

	var webhook event
	if err = unmarshal(c.Request, &webhook); err != nil {
		return err
	}

	switch webhook.EventType {
	case paddlenotification.EventTypeNameSubscriptionCreated:
		p.handleSubscriptionCreated(c)
	case paddlenotification.EventTypeNameSubscriptionUpdated:
		p.handleSubscriptionUpdate(c)
	}

	return nil
}

func (p *BillingWebhookController) handleSubscriptionCreated(c *gin.Context) error {
	var event paddlenotification.SubscriptionCreated
	if err := unmarshal(c.Request, &event); err != nil {
		return err
	}

	since, err := time.Parse(dateTimeLayout, *event.Data.StartedAt)
	if err != nil {
		return appErr.Errorf(appErr.RequestValidationError, "Invalid date time at `next_billed_at`")
	}

	till, err := time.Parse(dateTimeLayout, *event.Data.NextBilledAt)
	if err != nil {
		return appErr.Errorf(appErr.RequestValidationError, "Invalid date time at `next_billed_at`")
	}

	if len(event.Data.Items) == 0 {
		return appErr.Errorf(appErr.RequestValidationError, "User didn't subscribe to anything")
	}

	if len(event.Data.Items) > 1 {
		return appErr.Errorf(appErr.RequestValidationError, "User shouldn't be able to purchase more than 1 item")
	}

	items := make([]domain.SubscriptionItem, len(event.Data.Items))
	for i, item := range event.Data.Items {
		items[i] = domain.SubscriptionItem{
			Quantity:  item.Quantity,
			PriceID:   item.Price.ID,
			ProductID: item.Product.ID,
		}
	}

	return p.subscriptionService.HandleUpdatedSubscription(c.Request.Context(), domain.HandleUpdatedSubscriptionRequest{
		SubscriptionID: event.Data.ID,
		CustomerID:     event.Data.CustomerID,
		Items:          items,
		Till:           till,
		Since:          since,
	})
}

func (p *BillingWebhookController) handleSubscriptionUpdate(c *gin.Context) error {
	var event paddlenotification.SubscriptionUpdated
	if err := unmarshal(c.Request, &event); err != nil {
		return err
	}

	switch event.Data.Status {
	case paddlenotification.SubscriptionStatusPastDue:
		return p.subscriptionService.RemoveSubscription(c.Request.Context(), event.Data.ID)
	}

	since, err := time.Parse(dateTimeLayout, *event.Data.StartedAt)
	if err != nil {
		return appErr.Errorf(appErr.RequestValidationError, "Invalid date time at `next_billed_at`")
	}

	till, err := time.Parse(dateTimeLayout, *event.Data.NextBilledAt)
	if err != nil {
		return appErr.Errorf(appErr.RequestValidationError, "Invalid date time at `next_billed_at`")
	}

	if len(event.Data.Items) == 0 {
		return appErr.Errorf(appErr.RequestValidationError, "User didn't subscribe to anything")
	}

	if len(event.Data.Items) > 1 {
		return appErr.Errorf(appErr.RequestValidationError, "User shouldn't be able to purchase more than 1 item")
	}

	items := make([]domain.SubscriptionItem, len(event.Data.Items))
	for i, item := range event.Data.Items {
		items[i] = domain.SubscriptionItem{
			Quantity:  item.Quantity,
			PriceID:   item.Price.ID,
			ProductID: item.Product.ID,
		}
	}

	return p.subscriptionService.HandleUpdatedSubscription(c.Request.Context(), domain.HandleUpdatedSubscriptionRequest{
		CustomerID:     event.Data.CustomerID,
		SubscriptionID: event.Data.ID,
		Till:           till,
		Since:          since,
		Items:          items,
	})
}

func unmarshal(r *http.Request, v any) error {
	rawBody, err := io.ReadAll(r.Body)

	if err != nil {
		return tracerr.Wrap(err)
	}

	if err := json.Unmarshal(rawBody, v); err != nil {
		return appErr.Errorf(appErr.RequestValidationError, "invalid request body")
	}

	return nil
}
