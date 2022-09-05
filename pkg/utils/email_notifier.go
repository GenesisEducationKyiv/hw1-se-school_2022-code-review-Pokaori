package utils

import (
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

type EmailNotifier interface {
	SendEmails(emails []string)
}
type Dialer interface {
	DialAndSend(m ...*gomail.Message) error
}

type EmailBTCtoUAHNotifier struct {
	Dialer   Dialer
	Host     string
	Port     int
	From     string
	Password string
	Rate     float64
}

func (notifier *EmailBTCtoUAHNotifier) SendEmails(emails []string) {
	subject := "Changes in BTC to UAH currency rate"
	body := "The rate is " + fmt.Sprint(notifier.Rate)

	c := make(chan string)
	for _, email := range emails {
		msg := notifier.formMessage(email, subject, body)
		go notifier.sendEmail(c, msg)
	}
	for i := 0; i < len(emails); i++ {
		log.Println(<-c)
	}
	log.Println("Emails Sent!")
}

func (notifier *EmailBTCtoUAHNotifier) sendEmail(c chan string, msg *gomail.Message) {
	if err := notifier.Dialer.DialAndSend(msg); err != nil {
		c <- "Error for email " + msg.GetHeader("To")[0] + ":\n" + err.Error()
	} else {
		c <- "Successed for email " + msg.GetHeader("To")[0]
	}
}

func (notifier *EmailBTCtoUAHNotifier) formMessage(email, subject, body string) *gomail.Message {
	msg := gomail.NewMessage()
	msg.SetHeader("From", notifier.From)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)
	return msg
}
