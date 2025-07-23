package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/PaddleHQ/paddle-go-sdk/v4"
	"github.com/PaddleHQ/paddle-go-sdk/v4/pkg/paddlenotification"
	"github.com/compendium-tech/compendium/subscription-service/internal/config"
	"github.com/compendium-tech/compendium/subscription-service/internal/domain"
	appErr "github.com/compendium-tech/compendium/subscription-service/internal/error"
	"github.com/compendium-tech/compendium/subscription-service/internal/interop"
	"github.com/compendium-tech/compendium/subscription-service/internal/model"
	"github.com/compendium-tech/compendium/subscription-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/ztrue/tracerr"
)

const dateTimeLayout = time.RFC3339Nano

type PaddleWebhookController struct {
	subscriptionService   service.SubscriptionService
	paddleClient          paddle.SDK
	paddleProductIDs      config.PaddleProductIDs
	paddleWebhookVerifier paddle.WebhookVerifier
	userService           interop.UserService
}

func NewPaddleWebhookController(
	subscriptionService service.SubscriptionService,
	paddleProductIDs config.PaddleProductIDs,
	paddleClient paddle.SDK,
	paddleWebhookVerifier paddle.WebhookVerifier,
	userService interop.UserService) *PaddleWebhookController {
	return &PaddleWebhookController{
		subscriptionService:   subscriptionService,
		paddleClient:          paddleClient,
		paddleProductIDs:      paddleProductIDs,
		paddleWebhookVerifier: paddleWebhookVerifier,
		userService:           userService,
	}
}

func (p *PaddleWebhookController) MakeRoutes(e *gin.Engine) {
	e.POST("/paddleWebhook", appErr.HandleAppErr(p.Handle))
}

type webhook struct {
	EventID   string                           `json:"event_id"`
	EventType paddlenotification.EventTypeName `json:"event_type"`
}

func (p *PaddleWebhookController) Handle(c *gin.Context) error {
	ok, err := p.paddleWebhookVerifier.Verify(c.Request)

	if err != nil && (errors.Is(err, paddle.ErrMissingSignature) || errors.Is(err, paddle.ErrInvalidSignatureFormat)) {
		return appErr.Errorf(appErr.InvalidWebhookSignature, "Failed to verify request signature")
	}

	if err != nil {
		return tracerr.Wrap(err)
	}

	if !ok {
		return appErr.Errorf(appErr.InvalidWebhookSignature, "Failed to verify request signature")
	}

	var webhook webhook
	if err = unmarshal(c.Request, &webhook); err != nil {
		return err
	}

	switch webhook.EventType {
	case paddlenotification.EventTypeNameSubscriptionCreated:
		p.handleSubscriptionCreated(c)
	}

	return nil
}

func (p *PaddleWebhookController) handleSubscriptionCreated(c *gin.Context) error {
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

	// TODO: Move this to a business logic layer
	customer, err := p.paddleClient.GetCustomer(c.Request.Context(), &paddle.GetCustomerRequest{
		CustomerID: event.Data.CustomerID,
	})
	if err != nil {
		return err
	}

	account, err := p.userService.FindAccountByEmail(c.Request.Context(), customer.Email)
	if err != nil {
		return err
	}

	if len(event.Data.Items) == 0 {
		return appErr.Errorf(appErr.RequestValidationError, "User didn't subscribe to anything")
	}

	if len(event.Data.Items) > 1 {
		return appErr.Errorf(appErr.RequestValidationError, "User shouldn't be able to purchase more than 1 item")
	}

	item := event.Data.Items[0]

	var subscriptionLevel model.SubscriptionLevel

	switch item.Product.ID {
	case p.paddleProductIDs.StudentSubscriptionProductID:
		subscriptionLevel = model.SubscriptionLevelStudent
	case p.paddleProductIDs.TeamSubscriptionProductID:
		subscriptionLevel = model.SubscriptionLevelTeam
	case p.paddleProductIDs.CommunitySubscriptionProductID:
		subscriptionLevel = model.SubscriptionLevelCommunity
	default:
		return appErr.Errorf(appErr.RequestValidationError, "Unknown price ID")
	}

	return p.subscriptionService.PutSubscription(domain.PutSubscriptionRequest{
		UserID:            account.ID,
		SubscriptionLevel: subscriptionLevel,
		Till:              till,
		Since:             since,
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
