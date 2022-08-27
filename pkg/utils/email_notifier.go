package utils

import (
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

type EmailNotifier interface {
	SendEmails(emails []string)
}

type EmailBTCtoUAHNotifier struct {
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
		go notifier.sendEmail(c, msg, email)
	}
	for i := 0; i < len(emails); i++ {
		log.Println(<-c)
	}
	log.Println("Emails Sent!")
}

func (notifier *EmailBTCtoUAHNotifier) sendEmail(c chan string, msg *gomail.Message, email string) {
	dialer := gomail.NewDialer(notifier.Host, notifier.Port, notifier.From, notifier.Password)
	if err := dialer.DialAndSend(msg); err != nil {
		c <- "Error for email " + email + ":\n" + err.Error()
	} else {
		c <- "Successed for email " + email
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
