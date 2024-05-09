package amqpconsumer

type Option func(*amqpConsumerClient)

func ExchangeName(exchangeName string) Option {
	return func(p *amqpConsumerClient) {
		p.exchangeName = exchangeName
	}
}

func QueueName(queueName string) Option {
	return func(p *amqpConsumerClient) {
		p.queueName = queueName
	}
}

func BindingKey(bindingKey string) Option {
	return func(p *amqpConsumerClient) {
		p.bindingKey = bindingKey
	}
}

func ConsumerTag(consumerTag string) Option {
	return func(p *amqpConsumerClient) {
		p.consumerTag = consumerTag
	}
}

func WorkerPoolSize(workerPoolSize int) Option {
	return func(p *amqpConsumerClient) {
		p.workerPoolSize = workerPoolSize
	}
}
