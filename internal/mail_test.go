package internal_test

import (
	"fmt"
	"net/smtp"
	"testing"
	"time"

	mail "github.com/antonio-alexander/go-blog-mail/internal"

	imapclient "github.com/emersion/go-imap/client"
	gomail "github.com/go-mail/mail"
	assert "github.com/stretchr/testify/assert"
	pop3 "github.com/taknb2nch/go-pop3"
)

const (
	host       = "localhost"
	smtp_port  = 587
	pop3_port  = 110
	imap4_port = 143
	user       = "user@example.com"
	pass       = "password"
)

var (
	message = &mail.Message{
		Date:    fmt.Sprint(time.Now()),
		From:    "user@example.com",
		To:      []string{"user@example.com"},
		Subject: "Hello!",
		Body:    "Hello, World!",
	}
)

func TestSendEmailSmtp(t *testing.T) {
	//REVIEW: need to figure out why this doesn't work (headers are shit)
	//
	address := fmt.Sprintf("%s:%d", host, smtp_port)
	client, err := smtp.Dial(address)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	defer client.Close()

	//
	err = client.Auth(smtp.PlainAuth("", message.From, pass, host))
	assert.Nil(t, err)

	//
	err = mail.SendEmailSmtp(client, message)
	assert.Nil(t, err)
}

func TestSendEmailGomail(t *testing.T) {
	//
	d := gomail.NewDialer(host, smtp_port, user, pass)
	assert.NotNil(t, d)
	err := mail.SendEmailGomail(d, message)
	assert.Nil(t, err)
}

func TestReceiveEmailPop3(t *testing.T) {
	//
	address := fmt.Sprintf("%s:%d", host, pop3_port)
	client, err := pop3.Dial(address)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	defer func() {
		_ = client.Quit()
		_ = client.Close()
	}()

	//
	err = client.User(user)
	assert.Nil(t, err)
	err = client.Pass(pass)
	assert.Nil(t, err)

	//
	messages, err := mail.ReceiveEmailPop3(client)
	assert.Nil(t, err)
	for _, message := range messages {
		fmt.Println(message)
	}
}

func TestReceiveEmailImap(t *testing.T) {
	const mailbox string = "INBOX"

	// Connect to server
	address := fmt.Sprintf("%s:%d", host, imap4_port)
	client, err := imapclient.Dial(address)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	defer client.Logout()

	//login
	err = client.Login(user, pass)
	assert.Nil(t, err)

	//get email
	messages, err := mail.ReceiveEmailImap4(client, mailbox)
	assert.Nil(t, err)
	for _, message := range messages {
		fmt.Println(message)
	}
}
