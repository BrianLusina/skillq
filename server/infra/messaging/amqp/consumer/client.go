package amqpconsumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/messaging"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
	"github.com/pkg/errors"
	rabbitmq "github.com/rabbitmq/amqp091-go"
)

// amqpConsumerClient defines a consumer that handles consumption of messages from an AMQP Broker
type amqpConsumerClient struct {
	exchangeName   string
	queueName      string
	bindingKey     string
	consumerTag    string
	workerPoolSize int
	client         amqp.AmqpClient
	logger         logger.Logger
	handlers       map[string]func(payload []byte) error
}

// NewConsumer creates a new AMQP consumer
func NewConsumer[T rabbitmq.Delivery](client amqp.AmqpClient, log logger.Logger) (AmqpEventConsumer, error) {
	handlers := make(map[string]func(payload []byte) error)
	sub := &amqpConsumerClient{
		client:         client,
		logger:         log,
		exchangeName:   _exchangeName,
		queueName:      _queueName,
		bindingKey:     _bindingKey,
		consumerTag:    _consumerTag,
		workerPoolSize: _workerPoolSize,
		handlers:       handlers,
	}

	return sub, nil
}

// Consumes a message from a given queue. This is mostly a blocking operation
func (c *amqpConsumerClient) Consume(ctx context.Context, queue string) error {
	ch, err := c.CreateChannel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	msgs, err := ch.Consume(
		queue,             // queue
		c.consumerTag,     // consumer
		_consumeAutoAck,   // auto-ack
		_consumeExclusive, // exclusive
		_consumeNoLocal,   // no-local
		_consumeNoWait,    // no-wait
		nil,               // args
	)
	if err != nil {
		return fmt.Errorf("failed to register as consumer: %w", err)
	}

	forever := make(chan struct{})

	go func() {
		for delivery := range msgs {
			var msg messaging.ConsumeMessage
			err := json.Unmarshal(delivery.Body, &msg)
			if err != nil {
				c.logger.Errorf("failed to unmarshal message: %v", err)
				continue
			}
			if handler, ok := c.handlers[msg.Task]; ok {
				err = handler(msg.Payload)
				if err != nil {
					c.logger.Errorf("failed to handle message: %v", err)
					delivery.Nack(false, false)
					continue
				}
				err = delivery.Ack(false)
				if err != nil {
					c.logger.Errorf("failed to acknowledge deliver: %v", err)
				}
			} else {
				c.logger.Warn("task does not exist: %v", msg.Task)
				delivery.Nack(false, false)
			}
		}
	}()

	<-forever

	return nil
}

// StartConsumer starts a new consumer worker. Used for async workflows
func (c *amqpConsumerClient) StartConsumer(fn messaging.Worker[any]) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := c.CreateChannel()
	if err != nil {
		return errors.Wrapf(err, "failed to create channel")
	}

	defer ch.Close()

	deliveries, err := ch.Consume(
		c.queueName,
		c.consumerTag,
		_consumeAutoAck,
		_consumeExclusive,
		_consumeNoLocal,
		_consumeNoWait,
		nil,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to consume messages")
	}

	forever := make(chan bool)

	for i := 0; i < c.workerPoolSize; i++ {
		go fn(ctx, deliveries)
	}

	chanErr := <-ch.NotifyClose(make(chan *rabbitmq.Error))
	c.logger.Errorf("notify close: %v", chanErr)
	<-forever
	return chanErr
}

// AddHandler adds a handler that will handle consumption of messages from a queue
func (c *amqpConsumerClient) AddHandler(ctx context.Context, task string, handler func(payload []byte) error) {

}

// CreateChannel creates a rabbit MQ channel
func (c *amqpConsumerClient) CreateChannel() (*rabbitmq.Channel, error) {
	ch, err := c.client.AmqpConn.Channel()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create AMQP channel")
	}

	c.logger.Infof("Declaring exchange: %s", c.exchangeName)
	err = ch.ExchangeDeclare(
		c.exchangeName,
		_exchangeKind,
		_exchangeDurable,
		_exchangeAutoDelete,
		_exchangeInternal,
		_exchangeNoWait,
		nil,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to declare exchange: %s", c.exchangeName)
	}

	queue, err := ch.QueueDeclare(c.queueName, _queueDurable, _queueAutoDelete, _queueExclusive, _queueNoWait, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to declare queue: %s", c.queueName)
	}

	c.logger.Infof("Declaring queue, binding it to exchange: Queue: %v, messagesCount: %v, Consumer Count: %v, exchange: %v, bindingKey: %v",
		queue.Name, queue.Messages, queue.Consumers, c.exchangeName, c.bindingKey,
	)

	err = ch.QueueBind(
		queue.Name,
		c.bindingKey,
		c.exchangeName,
		_queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to bind queue: %s", queue.Name)
	}

	c.logger.Infof("Queue bound to exchange, starting to consumer from queue, consumerTag: %s", c.consumerTag)

	err = ch.Qos(
		_prefetchCount,  // prefetch count
		_prefetchSize,   // prefetch size
		_prefetchGlobal, // global
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failure to qos channel")
	}

	return ch, nil
}

// Configures an AMQP Event Consumer
func (c *amqpConsumerClient) Configure(opts ...Option) AmqpEventConsumer {
	for _, opt := range opts {
		opt(c)
	}
	return c
}
