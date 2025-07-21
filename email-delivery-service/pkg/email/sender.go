package email

type EmailMessage struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type EmailSender interface {
	SendMessage(msg EmailMessage) error
}
