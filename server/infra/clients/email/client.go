package email

import (
	"net/smtp"

	"github.com/BrianLusina/skillq/server/infra/logger"
)

type Client struct {
	from    string
	auth    smtp.Auth
	addr    string
	log     logger.Logger
	enabled bool
}

func New(config EmailClientConfig, log logger.Logger) EmailClient {
	auth := smtp.PlainAuth("", config.From, config.Password, config.Host)
	addr := config.Host + ":" + config.Port
	return &Client{
		from:    config.From,
		auth:    auth,
		addr:    addr,
		log:     log,
		enabled: config.Enabled,
	}
}

func (c *Client) Send(to string, body []byte) error {
	c.log.Infof("Sending email to %s", to)
	if c.enabled {
		toList := []string{to}
		err := smtp.SendMail(c.addr, c.auth, c.from, toList, body)
		if err != nil {
			c.log.Errorf("Failed to send email to %s with error: %s", to, err)
			return err
		}

		return nil
	}
	c.log.Infof("Email stubbed, skipping sending email")
	return nil
}
