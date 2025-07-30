package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/compendium-tech/compendium/common/pkg/error"
	"io"
	"time"

	"github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/PaddleHQ/paddle-go-sdk/v4/pkg/paddlenotification"
	"github.com/compendium-tech/compendium/subscription-service/internal/domain"
	"github.com/compendium-tech/compendium/subscription-service/internal/error"
	"github.com/compendium-tech/compendium/subscription-service/internal/service"
	"github.com/compendium-tech/compendium/subscription-service/internal/webhook"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ztrue/tracerr"
)

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
	e.POST("/v1/billingEvents", errorutils.Handle(p.handle))
}

type event struct {
	EventID   string                           `json:"event_id"`
	EventType paddlenotification.EventTypeName `json:"event_type"`
}

func (p *BillingWebhookController) handle(c *gin.Context) error {
	ok, err := p.webhookVerifier.Verify(c.Request)

	if err != nil && (errors.Is(err, paddle.ErrMissingSignature) || errors.Is(err, paddle.ErrInvalidSignatureFormat)) {
		return myerror.NewWithReason(myerror.InvalidWebhookSignatureError, "Failed to verify request signature")
	}

	if err != nil {
		return tracerr.Wrap(err)
	}

	if !ok {
		return myerror.NewWithReason(myerror.InvalidWebhookSignatureError, "Failed to verify request signature")
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return tracerr.Wrap(err)
	}

	var webhook event
	if err = unmarshal(body, &webhook); err != nil {
		return err
	}

	switch webhook.EventType {
	case paddlenotification.EventTypeNameSubscriptionCreated:
		return p.handleSubscriptionCreated(c, body)
	case paddlenotification.EventTypeNameSubscriptionUpdated:
		return p.handleSubscriptionUpdate(c, body)

	default:
		return myerror.NewWithReason(myerror.RequestValidationError, fmt.Sprintf("Unsupported event type %s", webhook.EventType))
	}
}

func (p *BillingWebhookController) handleSubscriptionCreated(c *gin.Context, body []byte) error {
	var event paddlenotification.SubscriptionCreated
	if err := unmarshal(body, &event); err != nil {
		return err
	}

	userIDString, ok := event.Data.CustomData["userId"].(string)
	if !ok {
		return myerror.NewWithReason(myerror.RequestValidationError, "userId wasn't provided")
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return myerror.NewWithReason(myerror.RequestValidationError, fmt.Sprintf("invalid user id %s", userIDString))
	}

	since, err := time.Parse(time.RFC3339, *event.Data.StartedAt)
	if err != nil {
		return myerror.NewWithReason(myerror.RequestValidationError, "Invalid date time at `next_billed_at`")
	}

	till, err := time.Parse(time.RFC3339, *event.Data.NextBilledAt)
	if err != nil {
		return myerror.NewWithReason(myerror.RequestValidationError, "Invalid date time at `next_billed_at`")
	}

	if len(event.Data.Items) == 0 {
		return myerror.NewWithReason(myerror.RequestValidationError, "User didn't subscribe to anything")
	}

	if len(event.Data.Items) > 1 {
		return myerror.NewWithReason(myerror.RequestValidationError, "User shouldn't be able to purchase more than 1 item")
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
		UserID:         userID,
		Items:          items,
		Till:           till,
		Since:          since,
	})
}

func (p *BillingWebhookController) handleSubscriptionUpdate(c *gin.Context, body []byte) error {
	var event paddlenotification.SubscriptionUpdated
	if err := unmarshal(body, &event); err != nil {
		return err
	}

	switch event.Data.Status {
	case paddlenotification.SubscriptionStatusPastDue:
		return p.subscriptionService.RemoveSubscription(c.Request.Context(), event.Data.ID)
	default:
		return nil
	}
}

func unmarshal(body []byte, v any) error {
	if err := json.Unmarshal(body, v); err != nil {
		return myerror.NewWithReason(myerror.RequestValidationError, "invalid request body")
	}

	return nil
}
