package internal

import (
	"io"
	"log"
	"net/mail"
	"strings"
)

func ParseMail(s string) (*Message, error) {
	r := strings.NewReader(s)
	m, err := mail.ReadMessage(r)
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := io.ReadAll(m.Body)
	if err != nil {
		return nil, err
	}
	return &Message{
		Date:    m.Header.Get("Date"),
		From:    m.Header.Get("From"),
		To:      strings.Split(m.Header.Get("To"), ","),
		Subject: m.Header.Get("Subject"),
		Body:    string(bytes),
	}, nil
}
