package sender

import (
	"context"

	"github.com/compendium-tech/compendium/email-delivery-service/pkg/domain"
)

type EmailSender interface {
	SendMessage(ctx context.Context, msg domain.EmailMessage) error
}
