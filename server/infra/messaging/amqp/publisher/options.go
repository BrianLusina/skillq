package publisher

// AmqpPublisherOption allows adding options to the AMQP publisher
type AmqpPublisherOption func(*AmqpPublisher)

// ExchangeName adds an exchange name to the publisher
func ExchangeName(exchangeName string) AmqpPublisherOption {
	return func(p *AmqpPublisher) {
		p.exchangeName = exchangeName
	}
}

// BindingKey allows adding a binding key to the publisher
func BindingKey(bindingKey string) AmqpPublisherOption {
	return func(p *AmqpPublisher) {
		p.bindingKey = bindingKey
	}
}

// MessageTypeName adds the name of the type of the message
func MessageTypeName(messageTypeName string) AmqpPublisherOption {
	return func(p *AmqpPublisher) {
		p.messageTypeName = messageTypeName
	}
}
