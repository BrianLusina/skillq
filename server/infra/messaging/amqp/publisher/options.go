package amqppublisher

import "github.com/BrianLusina/skillq/server/infra/messaging/amqp"

// Option allows adding options to the AMQP publisher
type Option func(*amqpPublisherClient)

// Exchange adds an exchange name to the publisher to be used when publishing messages
func Exchange(params amqp.ExchangeOptionParams) Option {
	return func(p *amqpPublisherClient) {
		p.exchangeName = params.Name
		p.exchangeKind = params.Kind
		p.exchangeDurable = params.Durable
		p.exchangeAutoDelete = params.AutoDelete
		p.exchangeInternal = params.Internal
		p.exchangeNoWait = params.NoWait
		p.exchangeArgs = params.Args
	}
}

func PublishConfig(params amqp.PublishOptionsParams) Option {
	return func(apc *amqpPublisherClient) {
		apc.publishMandatory = params.PublishMandatory
		apc.publishImmediate = params.PublishImmediate
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
