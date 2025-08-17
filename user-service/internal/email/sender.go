package email

type Message struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type Sender interface {
	SendMessage(msg Message)
}
