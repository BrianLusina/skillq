package email

// EmailClient is an interface for sending out email
type EmailClient interface {
	// Send sends out body to an email to the provided email address
	Send(to string, body []byte) error
}
