package internal

import (
	"net/smtp"

	gomail "github.com/go-mail/mail"
)

func SendEmailGomail(d *gomail.Dialer, message *Message) error {
	return d.DialAndSend(message.ToGomail())
}

func SendEmailSmtp(client *smtp.Client, message *Message) error {
	if err := client.Mail(message.From); err != nil {
		return err
	}
	for _, addr := range message.To {
		if err := client.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(message.Subject + message.Body))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return client.Quit()
}
