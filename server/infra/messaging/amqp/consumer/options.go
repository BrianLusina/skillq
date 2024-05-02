package consumer

type Option func(*AmqpConsumer)

func ExchangeName(exchangeName string) Option {
	return func(p *AmqpConsumer) {
		p.exchangeName = exchangeName
	}
}

func QueueName(queueName string) Option {
	return func(p *AmqpConsumer) {
		p.queueName = queueName
	}
}

func BindingKey(bindingKey string) Option {
	return func(p *AmqpConsumer) {
		p.bindingKey = bindingKey
	}
}

func ConsumerTag(consumerTag string) Option {
	return func(p *AmqpConsumer) {
		p.consumerTag = consumerTag
	}
}

func WorkerPoolSize(workerPoolSize int) Option {
	return func(p *AmqpConsumer) {
		p.workerPoolSize = workerPoolSize
	}
}
