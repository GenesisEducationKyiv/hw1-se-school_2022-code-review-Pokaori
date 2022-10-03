package interfaces

type EmailNotifier interface {
	SendEmails(emails []string, rate float64)
}
