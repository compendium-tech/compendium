package email

import (
	"fmt"
	"net/smtp"
)

type smtpEmailSender struct {
	host     string
	port     uint16
	username string
	password string
	from     string
}

func NewSmtpEmailSender(host string, port uint16, username string, password string, from string) Sender {
	return &smtpEmailSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *smtpEmailSender) SendMessage(msg Message) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	emailContent := []byte(
		"To: " + msg.To + "\r\n" +
			"From: " + s.from + "\r\n" +
			"Subject: " + msg.Subject + "\r\n" +
			"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
			msg.Body,
	)

	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.host, s.port),
		auth,
		s.from,
		[]string{msg.To},
		emailContent,
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
