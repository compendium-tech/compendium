package email

import (
	"strings"
	"text/template"

	"github.com/ztrue/tracerr"
)

type MessageBuilder interface {
	BuildSignUpMfaEmailMessage(to string, otp string) (Message, error)
	BuildSignInMfaEmailMessage(to string, otp string) (Message, error)
	BuildPasswordResetMfaEmailMessage(to string, otp string) (Message, error)
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

func (b *emailMessageBuilder) executeTemplate(name string, data any) (string, error) {
	body := new(strings.Builder)

	if err := b.templates.ExecuteTemplate(body, name, data); err != nil {
		return "", tracerr.Errorf("failed to execute email template, cause: %s", err.Error())
	}

	return body.String(), nil
}

func (b *emailMessageBuilder) BuildSignUpMfaEmailMessage(to, otp string) (Message, error) {
	type signUpEmailData struct {
		VerificationCode string
	}

	data := signUpEmailData{VerificationCode: otp}
	body, err := b.executeTemplate("sign_up.html", data)

	if err != nil {
		return Message{}, err
	}

	return Message{
		To:      to,
		Subject: "Verification code",
		Body:    body,
	}, nil
}

func (b *emailMessageBuilder) BuildSignInMfaEmailMessage(to, otp string) (Message, error) {
	type signInEmailData struct {
		VerificationCode string
	}

	data := signInEmailData{VerificationCode: otp}
	body, err := b.executeTemplate("sign_in.html", data)

	if err != nil {
		return Message{}, err
	}

	return Message{
		To:      to,
		Subject: "Verification code",
		Body:    body,
	}, err
}

func (b *emailMessageBuilder) BuildPasswordResetMfaEmailMessage(to string, otp string) (Message, error) {
	type passwordResetEmailData struct {
		VerificationCode string
	}

	data := passwordResetEmailData{VerificationCode: otp}
	body, err := b.executeTemplate("password_reset.html", data)

	if err != nil {
		return Message{}, err
	}

	return Message{
		To:      to,
		Subject: "Verification code",
		Body:    body,
	}, err
}
