package amqppublisher

import (
	"context"
	"fmt"
	"time"

	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/messaging"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	rabbitmq "github.com/rabbitmq/amqp091-go"
)

// amqpPublisherClient handles defines the methods used to handle publication of messages to a topic on a broker
type amqpPublisherClient struct {
	client           *amqp.AmqpClient
	exchangeName     string
	bindingKey       string
	messageTypeName  string
	publishMandatory bool
	publishImmediate bool
	logger           logger.Logger
}

// NewPublisher creates a new AMQP Publisher
func NewPublisher(client *amqp.AmqpClient, log logger.Logger) (AmqpEventPublisher, error) {
	publisher := &amqpPublisherClient{
		client:           client,
		logger:           log,
		exchangeName:     _exchangeName,
		bindingKey:       _bindingKey,
		messageTypeName:  _messageTypeName,
		publishMandatory: _publishMandatory,
		publishImmediate: _publishImmediate,
	}

	return publisher, nil
}

// Publish publishes a message to a given topic
func (p *amqpPublisherClient) Publish(ctx context.Context, message messaging.Message) error {
	body, err := message.ToBytes()
	if err != nil {
		return errors.Wrapf(err, "failed to parse message event")
	}

	var messageTypeName string
	if message.Topic == "" {
		messageTypeName = p.messageTypeName
	} else {
		messageTypeName = message.Topic
	}

	amqpChan, err := p.client.AmqpConn.Channel()
	if err != nil {
		p.logger.Errorf("Failed to open channel: %v", err)
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer amqpChan.Close()

	p.logger.Infof("Publishing message to exchange: %s, with routingKey: %s", p.exchangeName, p.bindingKey)

	err = amqpChan.PublishWithContext(
		ctx,
		p.exchangeName,
		p.bindingKey,
		p.publishMandatory,
		p.publishImmediate,
		rabbitmq.Publishing{
			ContentType:  message.ContentType,
			DeliveryMode: rabbitmq.Persistent,
			MessageId:    uuid.New().String(),
			Timestamp:    time.Now(),
			Body:         body,
			Type:         messageTypeName,
		},
	)
	if err != nil {
		p.logger.Errorf("Failed to publish message to exchange %s with error: %v", p.exchangeName, err)
		return errors.Wrapf(err, "failed to publish message: %v", err)
	}

	p.logger.Infof("Successfully published message to exchange: %s, with routingKey: %s", p.exchangeName, p.bindingKey)

	return nil
}

// Close closes connection to a broker
func (p *amqpPublisherClient) Close() error {
	return p.client.Close()
}

func (p *amqpPublisherClient) Configure(opts ...Option) AmqpEventPublisher {
	for _, opt := range opts {
		opt(p)
	}

	return p
}
