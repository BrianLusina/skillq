package amqppublisher

// Option allows adding options to the AMQP publisher
type Option func(*amqpPublisherClient)

type ExchangeOptionParams struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       map[string]any
}

// Exchange adds an exchange name to the publisher to be used when publishing messages
func Exchange(params ExchangeOptionParams) Option {
	return func(p *amqpPublisherClient) {
		amqpChan, err := p.client.AmqpConn.Channel()
		if err != nil {
			p.logger.Errorf("Failed to open channel: %v", err)
		}
		defer amqpChan.Close()

		err = amqpChan.ExchangeDeclare(
			params.Name,
			params.Kind,
			params.Durable,
			params.AutoDelete,
			params.Internal,
			params.NoWait,
			params.Args,
		)
		if err != nil {
			p.logger.Infof("Failed to declare exchange: %s, with kind: %s", params.Name, params.Kind)
		}
		p.exchangeName = params.Name
	}
}

// BindingKey allows adding a binding key to the publisher
func BindingKey(bindingKey string) Option {
	return func(p *amqpPublisherClient) {
		p.bindingKey = bindingKey
	}
}

// MessageTypeName adds the name of the type of the message
func MessageTypeName(messageTypeName string) Option {
	return func(p *amqpPublisherClient) {
		p.messageTypeName = messageTypeName
	}
}
