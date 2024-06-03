package email

import "net/smtp"

type Client struct {
	from string
	auth smtp.Auth
	addr string
}

func New(config EmailClientConfig) EmailClient {
	// Load env vars
	auth := smtp.PlainAuth("", config.From, config.Password, config.Host)
	addr := config.Host + ":" + config.Port
	return &Client{from: config.From, auth: auth, addr: addr}
}

func (c *Client) Send(to string, body []byte) error {
	toList := []string{to}
	return smtp.SendMail(c.addr, c.auth, c.from, toList, body)
}
