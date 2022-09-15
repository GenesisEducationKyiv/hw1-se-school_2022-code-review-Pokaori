package mocks

import (
	"log"
)

type EmailNotifierMock struct {
}

func (notifier *EmailNotifierMock) SendEmails(emails []string, rate float64) {
	log.Println("Emails Sent!")
}
