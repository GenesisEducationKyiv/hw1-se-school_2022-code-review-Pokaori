package utils

import (
	"bitcoin-service/pkg/config"
	"testing"

	"gopkg.in/gomail.v2"
)

type DialerMock struct {
	emails []string
}

func (dialer *DialerMock) DialAndSend(m ...*gomail.Message) error {
	dialer.emails = append(dialer.emails, m[0].GetHeader("To")[0])
	return nil
}

func (dialer *DialerMock) AssertCalledEmails(emails []string) bool {
	if len(emails) != len(dialer.emails) {
		return false
	}

	diff := make(map[string]int, len(emails))
	for _, email := range emails {
		diff[email]++
	}
	for _, email := range dialer.emails {
		// If the string _y is not in diff bail out early
		if _, ok := diff[email]; !ok {
			return false
		}
		diff[email] -= 1
		if diff[email] == 0 {
			delete(diff, email)
		}
	}
	return len(diff) == 0

}

func TestNotifierSendToCorrectEmails(t *testing.T) {

	var emails = []string{"ex1@example.com", "ex2@example.com", "ex3@example.com"}
	dialer := &DialerMock{}
	var notifier EmailNotifier = &EmailBTCtoUAHNotifier{
		Dialer:   dialer,
		Host:     config.Settings.EmailHost,
		Port:     config.Settings.EmailPort,
		From:     config.Settings.EmailName,
		Password: config.Settings.EmailPass,
		Rate:     0.1,
	}

	notifier.SendEmails(emails)

	if !dialer.AssertCalledEmails(emails) {
		t.Fatalf("Messages are sent to wrong emails. Should be %v, instead got %v", emails, dialer.emails)
	}
}
