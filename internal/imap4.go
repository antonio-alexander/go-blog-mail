package internal

import (
	imap "github.com/emersion/go-imap"
	client "github.com/emersion/go-imap/client"
)

func ReceiveEmailImap4(client *client.Client, mailbox string) ([]Message, error) {
	var section imap.BodySectionName
	var messages []Message

	mailboxStatus, err := client.Select("INBOX", false)
	if err != nil {
		return nil, err
	}
	if mailboxStatus.Messages == 0 {
		return []Message{}, nil
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(1, mailboxStatus.Messages)
	chMessage := make(chan *imap.Message, 10)
	chDone := make(chan error, 1)
	go func() {
		chDone <- client.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, chMessage)
	}()
	for msg := range chMessage {
		var bytes []byte

		if _, err = msg.GetBody(&section).Read(bytes); err != nil {
			return nil, err
		}
		message, err := ParseMail(string(bytes))
		if err != nil {
			return nil, err
		}
		messages = append(messages, *message)
	}
	if err := <-chDone; err != nil {
		return nil, err
	}
	return messages, nil
}
