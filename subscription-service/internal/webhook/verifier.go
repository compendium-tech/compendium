package webhook

import "net/http"

type WebhookVerifier interface {
	Verify(req *http.Request) (bool, error)
}
