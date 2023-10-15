package internal

import (
	pop3 "github.com/taknb2nch/go-pop3"
)

func ReceiveEmailPop3(client *pop3.Client) ([]Message, error) {
	var messagesRetrieved []int
	var messages []Message

	messageInfos, err := client.ListAll()
	if err != nil {
		return nil, err
	}
	if len(messageInfos) == 0 {
		return []Message{}, nil
	}
	for _, messageInfo := range messageInfos {

		s, err := client.Retr(messageInfo.Number)
		if err != nil {
			return nil, err
		}
		message, err := ParseMail(s)
		if err != nil {
			return nil, err
		}
		messages = append(messages, *message)
		messagesRetrieved = append(messagesRetrieved, messageInfo.Number)
	}
	for _, messageRetreived := range messagesRetrieved {
		if err = client.Dele(messageRetreived); err != nil {
			return nil, err
		}
	}
	return messages, nil
}
