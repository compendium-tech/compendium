package httpv1

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/PaddleHQ/paddle-go-sdk/v4/pkg/paddlenotification"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	httputils "github.com/compendium-tech/compendium/common/pkg/http"

	"github.com/compendium-tech/compendium/subscription-service/internal/domain"
	myerror "github.com/compendium-tech/compendium/subscription-service/internal/error"
	"github.com/compendium-tech/compendium/subscription-service/internal/service"
	"github.com/compendium-tech/compendium/subscription-service/internal/webhook"
)

type BillingWebhookController struct {
	billingEventHandlerService service.BillingEventHandlerService
	webhookVerifier            webhook.WebhookVerifier
}

func NewBillingWebhookController(
	billingEventHandlerService service.BillingEventHandlerService,
	webhookVerifier webhook.WebhookVerifier) *BillingWebhookController {
	return &BillingWebhookController{
		billingEventHandlerService: billingEventHandlerService,
		webhookVerifier:            webhookVerifier,
	}
}

func (p *BillingWebhookController) MakeRoutes(e *gin.Engine) {
	var eh httputils.ErrorHandler

	e.POST("/v1/billingEvents", eh.Handle(p.handle))
}

type event struct {
	EventID   string                           `json:"event_id"`
	EventType paddlenotification.EventTypeName `json:"event_type"`
}

func (p *BillingWebhookController) handle(c *gin.Context) {
	ok, err := p.webhookVerifier.Verify(c.Request)

	if err != nil && (errors.Is(err, paddle.ErrMissingSignature) || errors.Is(err, paddle.ErrInvalidSignatureFormat)) {
		myerror.NewWithReason(myerror.InvalidWebhookSignatureError, "Failed to verify request signature").Throw()
	}

	if err != nil {
		panic(err)
	}

	if !ok {
		myerror.NewWithReason(myerror.InvalidWebhookSignatureError, "Failed to verify request signature").Throw()
	}

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		panic(err)
	}

	var webhookEvent event
	unmarshal(body, &webhookEvent)

	switch webhookEvent.EventType {
	case paddlenotification.EventTypeNameSubscriptionCreated:
		p.handleSubscriptionCreated(c, body)
	case paddlenotification.EventTypeNameSubscriptionUpdated:
		p.handleSubscriptionUpdate(c, body)

	default:
		myerror.NewWithReason(
			myerror.RequestValidationError,
			fmt.Sprintf("Unsupported event type %s", webhookEvent.EventType)).Throw()
	}
}

func (p *BillingWebhookController) handleSubscriptionCreated(c *gin.Context, body []byte) {
	var event paddlenotification.SubscriptionCreated
	unmarshal(body, &event)

	userIDString, ok := event.Data.CustomData["userId"].(string)
	if !ok {
		myerror.NewWithReason(myerror.RequestValidationError, "userId wasn't provided").Throw()
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		myerror.NewWithReason(myerror.RequestValidationError, fmt.Sprintf("invalid user id %s", userIDString)).Throw()
	}

	since, err := time.Parse(time.RFC3339, *event.Data.StartedAt)
	if err != nil {
		myerror.NewWithReason(myerror.RequestValidationError, "Invalid date time at `next_billed_at`").Throw()
	}

	till, err := time.Parse(time.RFC3339, *event.Data.NextBilledAt)
	if err != nil {
		myerror.NewWithReason(myerror.RequestValidationError, "Invalid date time at `next_billed_at`").Throw()
	}

	if len(event.Data.Items) == 0 {
		myerror.NewWithReason(myerror.RequestValidationError, "User didn't subscribe to anything").Throw()
	}

	if len(event.Data.Items) > 1 {
		myerror.NewWithReason(myerror.RequestValidationError, "User shouldn't be able to purchase more than 1 item").Throw()
	}

	items := make([]domain.SubscriptionItem, len(event.Data.Items))
	for i, item := range event.Data.Items {
		items[i] = domain.SubscriptionItem{
			Quantity:  item.Quantity,
			PriceID:   item.Price.ID,
			ProductID: item.Product.ID,
		}
	}

	p.billingEventHandlerService.HandleUpdatedSubscription(c.Request.Context(), domain.HandleUpdatedSubscriptionRequest{
		SubscriptionID: event.Data.ID,
		UserID:         userID,
		Items:          items,
		Till:           till,
		Since:          since,
	})
}

func (p *BillingWebhookController) handleSubscriptionUpdate(c *gin.Context, body []byte) {
	var event paddlenotification.SubscriptionUpdated
	unmarshal(body, &event)

	switch event.Data.Status {
	case paddlenotification.SubscriptionStatusPastDue:
		p.billingEventHandlerService.CancelSubscription(c.Request.Context(), event.Data.ID)
	}
}

func unmarshal(body []byte, v any) {
	if err := json.Unmarshal(body, v); err != nil {
		myerror.NewWithReason(myerror.RequestValidationError, "invalid request body").Throw()
	}
}
