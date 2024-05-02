package consumer

import (
	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
)

// AmqpConsumer defines a consumer that handles consumption of messages from an AMQP Broker
type AmqpConsumer struct {
	exchangeName   string
	queueName      string
	bindingKey     string
	consumerTag    string
	workerPoolSize int
	client         amqp.AmqpClient
	logger         logger.Logger
}

// NewConsumer creates a new AMQP consumer
func NewConsumer(client amqp.AmqpClient, log logger.Logger, opts ...Option) (*AmqpConsumer, error) {
	sub := &AmqpConsumer{
		client:         client,
		logger:         log,
		exchangeName:   _exchangeName,
		queueName:      _queueName,
		bindingKey:     _bindingKey,
		consumerTag:    _consumerTag,
		workerPoolSize: _workerPoolSize,
	}

	for _, opt := range opts {
		opt(sub)
	}

	return sub, nil
}
