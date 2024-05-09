package amqppublisher

// Option allows adding options to the AMQP publisher
type Option func(*amqpPublisherClient)

// ExchangeName adds an exchange name to the publisher
func ExchangeName(exchangeName string) Option {
	return func(p *amqpPublisherClient) {
		p.exchangeName = exchangeName
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
