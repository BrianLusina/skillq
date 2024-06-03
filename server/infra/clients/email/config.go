package email

// EmailClientConfig to setup email client for sending email
type EmailClientConfig struct {
	Host     string
	Port     string
	Password string
	From     string
}
