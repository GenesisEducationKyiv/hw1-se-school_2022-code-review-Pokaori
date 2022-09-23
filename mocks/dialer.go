package mocks

import "gopkg.in/gomail.v2"

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
		// If the string  is not in diff bail out early
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

func (dialer *DialerMock) GetEmails() []string {
	return dialer.emails
}
