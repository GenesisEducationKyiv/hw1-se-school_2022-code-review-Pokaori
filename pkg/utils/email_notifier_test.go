package utils

import (
	"bitcoin-service/mocks"
	"bitcoin-service/pkg/config"
	"testing"
)

func TestNotifierSendToCorrectEmails(t *testing.T) {

	var emails = []string{"ex1@example.com", "ex2@example.com", "ex3@example.com"}
	dialer := &mocks.DialerMock{}
	notifier := &EmailBTCtoUAHNotifier{
		Dialer:   dialer,
		Host:     config.Settings.EmailHost,
		Port:     config.Settings.EmailPort,
		From:     config.Settings.EmailName,
		Password: config.Settings.EmailPass,
		Rate:     0.1,
	}

	notifier.SendEmails(emails)

	if !dialer.AssertCalledEmails(emails) {
		t.Fatalf("Messages are sent to wrong emails. Should be %v, instead got %v", emails, dialer.GetEmails())
	}
}
