package email

import (
	"fmt"
	"net/smtp"
	"strings"
	"text/template"

	"github.com/ztrue/tracerr"
)

type smtpConnectionSettings struct {
	host     string
	port     uint16
	username string
	password string
	from     string
}

type SmtpEmailSender struct {
	smtpConnectionSettings

	emailTemplates *template.Template
}

func NewSmtpEmailSender(host string, port uint16, username, password, from string) (*SmtpEmailSender, error) {
	emailTemplates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return nil, err
	}

	if emailTemplates == nil {
		return nil, fmt.Errorf("email templates were not initialized correctly. perhaps `templates` folder doesn't exist?")
	}

	return &SmtpEmailSender{
		smtpConnectionSettings: smtpConnectionSettings{
			host:     host,
			port:     port,
			username: username,
			password: password,
			from:     from,
		},
		emailTemplates: emailTemplates,
	}, nil
}

func (s *SmtpEmailSender) sendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	msg := []byte("To: " + to + "\r\n" +
		"From: " + s.from + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(fmt.Sprintf("%s:%d", s.host, s.port), auth, s.from, []string{to}, msg)
	if err != nil {
		return tracerr.Errorf("failed to send email: %v", err)
	}

	return nil
}

func (s *SmtpEmailSender) executeTemplate(name string, data any) (string, error) {
	body := new(strings.Builder)

	if err := s.emailTemplates.ExecuteTemplate(body, name, data); err != nil {
		return "", tracerr.Errorf("failed to execute email template, cause: %s", err.Error())
	}

	return body.String(), nil
}

type signUpEmailData struct {
	VerificationCode string
}

func (s *SmtpEmailSender) SendSignUpMfaEmail(to, otp string) error {
	data := signUpEmailData{VerificationCode: otp}
	body, err := s.executeTemplate("sign_up.html", data)

	if err != nil {
		return err
	}

	return s.sendEmail(to, "Verification code", body)
}

type signInEmailData struct {
	VerificationCode string
}

func (s *SmtpEmailSender) SendSignInMfaEmail(to, otp string) error {
	data := signInEmailData{VerificationCode: otp}
	body, err := s.executeTemplate("sign_in.html", data)

	if err != nil {
		return err
	}

	return s.sendEmail(to, "Verification code", body)
}
