package email

import (
	"strings"
	"text/template"

	"github.com/ztrue/tracerr"
)

type MessageBuilder interface {
	SignUpEmail(to string, otp string) Message
	SignInEmail(to string, otp string) Message
	PasswordResetEmail(to string, otp string) Message
}

type emailMessageBuilder struct {
	templates *template.Template
}

func NewMessageBuilder() (MessageBuilder, error) {
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return nil, err
	}

	if templates == nil {
		return nil, tracerr.Errorf("email templates were not initialized correctly. perhaps `templates` folder doesn't exist")
	}

	return &emailMessageBuilder{
		templates: templates,
	}, nil
}

func (b *emailMessageBuilder) executeTemplate(name string, data any) string {
	body := new(strings.Builder)

	if err := b.templates.ExecuteTemplate(body, name, data); err != nil {
		panic(err)
	}

	return body.String()
}

func (b *emailMessageBuilder) SignUpEmail(to, otp string) Message {
	type signUpEmailData struct {
		VerificationCode string
	}

	data := signUpEmailData{VerificationCode: otp}
	body := b.executeTemplate("sign_up.html", data)

	return Message{
		To:      to,
		Subject: "Verification code",
		Body:    body,
	}
}

func (b *emailMessageBuilder) SignInEmail(to, otp string) Message {
	type signInEmailData struct {
		VerificationCode string
	}

	data := signInEmailData{VerificationCode: otp}
	body := b.executeTemplate("sign_in.html", data)

	return Message{
		To:      to,
		Subject: "Verification code",
		Body:    body,
	}
}

func (b *emailMessageBuilder) PasswordResetEmail(to string, otp string) Message {
	type passwordResetEmailData struct {
		VerificationCode string
	}

	data := passwordResetEmailData{VerificationCode: otp}
	body := b.executeTemplate("password_reset.html", data)

	return Message{
		To:      to,
		Subject: "Verification code",
		Body:    body,
	}
}
