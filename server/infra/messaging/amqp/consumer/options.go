package amqpconsumer

type Option func(*amqpConsumerClient)

func ExchangeName(name string) Option {
	return func(c *amqpConsumerClient) {
		c.exchangeName = name
	}
}

func QueueName(name string) Option {
	return func(c *amqpConsumerClient) {
		c.queueName = name
	}
}

func BindingKey(bindingKey string) Option {
	return func(c *amqpConsumerClient) {
		c.bindingKey = bindingKey
	}
}

func ConsumerTag(tag string) Option {
	return func(c *amqpConsumerClient) {
		c.consumerTag = tag
	}
}

func WorkerPoolSize(size int) Option {
	return func(c *amqpConsumerClient) {
		c.workerPoolSize = size
	}
}
