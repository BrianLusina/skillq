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
	exchangeName       string
	exchangeKind       string
	exchangeDurable    bool
	exchangeAutoDelete bool
	exchangeInternal   bool
	exchangeNoWait     bool
	exchangeArgs       map[string]any
	queueName          string
	queueDurable       bool
	queueAutoDelete    bool
	queueExclusive     bool
	queueNoWait        bool
	queueArgs          map[string]any
	bindingKey         string
	consumerTag        string
	consumeAutoAck     bool
	consumeExclusive   bool
	consumeNoLocal     bool
	consumeNoWait      bool
	consumerArgs       map[string]any
	qosPrefetchCount   int
	qosPrefetchSize    int
	qosPrefetchGlobal  bool
	workerPoolSize     int
	client             *amqp.AmqpClient
	logger             logger.Logger
	handlers           map[string]func(payload []byte) error
}

// NewConsumer creates a new AMQP consumer
func NewConsumer(client *amqp.AmqpClient, log logger.Logger) (AmqpEventConsumer, error) {
	handlers := make(map[string]func(payload []byte) error)
	sub := &amqpConsumerClient{
		client:             client,
		logger:             log,
		exchangeName:       _exchangeName,
		exchangeKind:       _exchangeKind,
		exchangeDurable:    _exchangeDurable,
		exchangeAutoDelete: _exchangeAutoDelete,
		exchangeInternal:   _exchangeInternal,
		exchangeNoWait:     _exchangeNoWait,
		exchangeArgs:       map[string]any{},
		queueName:          _queueName,
		queueDurable:       _queueDurable,
		queueAutoDelete:    _queueAutoDelete,
		queueExclusive:     _queueExclusive,
		queueNoWait:        _queueNoWait,
		queueArgs:          map[string]any{},
		bindingKey:         _bindingKey,
		consumerTag:        _consumerTag,
		consumeAutoAck:     _consumeAutoAck,
		consumeExclusive:   _consumeExclusive,
		consumeNoLocal:     _consumeNoLocal,
		consumeNoWait:      _consumeNoWait,
		consumerArgs:       map[string]any{},
		qosPrefetchCount:   _prefetchCount,
		qosPrefetchSize:    _prefetchSize,
		qosPrefetchGlobal:  _prefetchGlobal,
		workerPoolSize:     _workerPoolSize,
		handlers:           handlers,
	}

	return sub, nil
}

// Consumes a message from a given queue. This is mostly a blocking operation
func (c *amqpConsumerClient) Consume(ctx context.Context, queue string) error {
	ch, err := c.createChannel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	msgs, err := ch.Consume(
		queue,              // queue
		c.consumerTag,      // consumer
		c.consumeAutoAck,   // auto-ack
		c.consumeExclusive, // exclusive
		c.consumeNoLocal,   // no-local
		c.consumeNoWait,    // no-wait
		c.consumerArgs,     // args
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
			if handler, ok := c.handlers[msg.Topic]; ok {
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
				c.logger.Warn("task does not exist: %v", msg.Topic)
				delivery.Nack(false, false)
			}
		}
	}()

	<-forever

	return nil
}

// StartConsumer starts a new consumer worker. Used for async workflows
func (c *amqpConsumerClient) StartConsumer(fn func(ctx context.Context, message <-chan rabbitmq.Delivery)) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := c.createChannel()
	if err != nil {
		c.logger.Errorf("Failed to create channel with error: %s", err.Error())
		return errors.Wrapf(err, "failed to create channel")
	}
	c.logger.Info("Successfully created channel")

	defer func() {
		c.logger.Info("Closing channel connection")
		err := ch.Close()
		if err != nil {
			c.logger.Errorf("Failed to close channel connection %s", err.Error())
		}
	}()

	deliveries, err := ch.Consume(
		c.queueName,
		c.consumerTag,
		c.consumeAutoAck,
		c.consumeExclusive,
		c.consumeNoLocal,
		c.consumeNoWait,
		c.consumerArgs,
	)
	if err != nil {
		c.logger.Errorf("Failed to consume messages with error: %s", err.Error())
		return errors.Wrapf(err, "failed to consume messages")
	}

	c.logger.Infof("Retrieved deliveries of count %d from queue %s", len(deliveries), c.queueName)

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

// createChannel creates a rabbit MQ channel
func (c *amqpConsumerClient) createChannel() (*rabbitmq.Channel, error) {
	ch, err := c.client.AmqpConn.Channel()
	if err != nil {
		c.logger.Errorf("Failed to create AMQP channel with error: %s", err.Error())
		return nil, errors.Wrapf(err, "failed to create AMQP channel")
	}

	c.logger.Infof("Declaring exchange: %s", c.exchangeName)
	err = ch.ExchangeDeclare(
		c.exchangeName,
		c.exchangeKind,
		c.exchangeDurable,
		c.exchangeAutoDelete,
		c.exchangeInternal,
		c.exchangeNoWait,
		c.exchangeArgs,
	)

	if err != nil {
		c.logger.Errorf("Failed to declare exchange: %s with error: %s", c.exchangeName, err.Error())
		return nil, errors.Wrapf(err, "failed to declare exchange: %s", c.exchangeName)
	}

	c.logger.Infof("Declaring queue: %s", c.queueName)
	queue, err := ch.QueueDeclare(c.queueName, c.queueDurable, c.queueAutoDelete, c.queueExclusive, c.queueNoWait, c.queueArgs)
	if err != nil {
		c.logger.Errorf("Failed to declare queue: %s with error %s", c.queueName, err.Error())
		return nil, errors.Wrapf(err, "failed to declare queue: %s", c.queueName)
	}

	c.logger.Infof("Declaring queue %s, binding it to exchange %s, messagesCount: %v, Consumer Count: %v, bindingKey: %v",
		queue.Name, c.exchangeName, queue.Messages, queue.Consumers, c.bindingKey,
	)

	err = ch.QueueBind(
		queue.Name,
		c.bindingKey,
		c.exchangeName,
		c.queueNoWait,
		nil,
	)
	if err != nil {
		c.logger.Errorf("Failed to bind exchange %s to queue %s with error: %s", c.exchangeName, queue.Name, err.Error())
		return nil, errors.Wrapf(err, "failed to bind queue: %s", queue.Name)
	}

	c.logger.Infof("Queue bound to exchange %s, starting to consume from queue %s, consumerTag: %s", c.exchangeName, queue.Name, c.consumerTag)

	err = ch.Qos(
		c.qosPrefetchCount,
		c.qosPrefetchSize,
		c.qosPrefetchGlobal,
	)
	if err != nil {
		c.logger.Errorf("Failed to QOS channel with error: %s", err.Error())
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
