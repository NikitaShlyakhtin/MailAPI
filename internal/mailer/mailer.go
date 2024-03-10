package mailer

import (
	"time"

	"github.com/go-mail/mail/v2"
)

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

func (m Mailer) Send(recipient, subject, textBody, htmlBody string) error {
	message := mail.NewMessage()
	message.SetHeader("From", m.sender)
	message.SetHeader("To", recipient)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", textBody)
	message.AddAlternative("text/html", htmlBody)

	err := m.dialer.DialAndSend(message)
	if err != nil {
		return err
	}

	return nil
}
