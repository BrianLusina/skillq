package amqpconsumer

import "github.com/BrianLusina/skillq/server/infra/messaging/amqp"

type Option func(*amqpConsumerClient)

// Exchange configures an exchange for the consumer to construct
func Exchange(params amqp.ExchangeOptionParams) Option {
	return func(c *amqpConsumerClient) {
		c.exchangeName = params.Name
		c.exchangeKind = params.Kind
		c.exchangeDurable = params.Durable
		c.exchangeAutoDelete = params.AutoDelete
		c.exchangeInternal = params.Internal
		c.exchangeNoWait = params.NoWait
	}
}

func Queue(params amqp.QueueOptionParams) Option {
	return func(c *amqpConsumerClient) {
		c.queueName = params.Name
		c.queueDurable = params.Durable
		c.queueAutoDelete = params.AutoDelete
		c.queueExclusive = params.Exclusive
		c.queueNoWait = params.NoWait
		c.queueArgs = params.Args
	}
}

func BindingKey(bindingKey string) Option {
	return func(c *amqpConsumerClient) {
		c.bindingKey = bindingKey
	}
}

// Consumer sets up the consumer options
func Consumer(params amqp.ConsumerOptionParams) Option {
	return func(c *amqpConsumerClient) {
		c.consumerTag = params.Tag
		c.consumeAutoAck = params.AutoAck
		c.consumeExclusive = params.Exclusive
		c.consumeNoLocal = params.NoLocal
		c.consumeNoWait = params.NoWait
		c.consumerArgs = params.Args
	}
}

// Qos sets up the consumer QOS for message consumption
func Qos(params amqp.QosOptionParams) Option {
	return func(c *amqpConsumerClient) {
		c.qosPrefetchCount = params.PrefetchCount
		c.qosPrefetchSize = params.PrefetchSize
		c.qosPrefetchGlobal = params.PrefetchGlobal
	}
}

func WorkerPoolSize(size int) Option {
	return func(c *amqpConsumerClient) {
		c.workerPoolSize = size
	}
}
