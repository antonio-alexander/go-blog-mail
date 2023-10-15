package internal

import (
	"fmt"
	"io"
	"net/mail"
	"strings"

	gomail "github.com/go-mail/mail"
)

type Message struct {
	Date    string   `json:"date"`
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

func (m *Message) Parse(s string) error {
	r := strings.NewReader(s)
	message, err := mail.ReadMessage(r)
	if err != nil {
		return err
	}
	header := message.Header
	m.Date = header.Get("Date")
	m.From = header.Get("From")
	m.To = strings.Split(header.Get("To"), ",")
	m.Subject = header.Get("Subject")
	body, err := io.ReadAll(message.Body)
	if err != nil {
		return err
	}
	m.Body = string(body)
	return nil
}

func (m *Message) ToGomail() *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader("From", m.From)
	message.SetHeader("To", strings.Join(m.To, ","))
	message.SetHeader("Subject", m.Subject)
	message.SetBody("text/html", m.Body)
	return message
}

func (m *Message) String() string {
	return fmt.Sprintf("Date: %s\nFrom: %s\nTo: %s\nSubject: %s\nBody: %s\n",
		m.Date, m.From, m.To, m.Subject, m.Body)
}
