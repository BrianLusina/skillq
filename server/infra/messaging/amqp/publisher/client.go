package amqppublisher

import (
	"context"
	"time"

	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	rabbitmq "github.com/rabbitmq/amqp091-go"
)

// amqpPublisherClient handles defines the methods used to handle publication of messages to a topic on a broker
type amqpPublisherClient struct {
	client          *amqp.AmqpClient
	exchangeName    string
	bindingKey      string
	messageTypeName string
	logger          logger.Logger
}

// NewPublisher creates a new AMQP Publisher
func NewPublisher(client *amqp.AmqpClient, log logger.Logger) (AmqpEventPublisher, error) {
	publisher := &amqpPublisherClient{
		client:          client,
		logger:          log,
		exchangeName:    _exchangeName,
		bindingKey:      _bindingKey,
		messageTypeName: _messageTypeName,
	}

	return publisher, nil
}

// Publish publishes a message to a given topic
func (p *amqpPublisherClient) Publish(ctx context.Context, body []byte, contentType string) error {
	amqpChan, err := p.client.AmqpConn.Channel()
	if err != nil {
		p.logger.Errorf("Failed to open channel: %v", err)
		return errors.Wrapf(err, "failed to open a channel: %w", err)
	}
	defer amqpChan.Close()

	p.logger.Infof("Publishing message to exchange: %s, with routingKey: %s", p.exchangeName, p.bindingKey)

	err = amqpChan.PublishWithContext(
		ctx,
		p.exchangeName,
		p.bindingKey,
		_publishMandatory,
		_publishImmediate,
		rabbitmq.Publishing{
			ContentType:  contentType,
			DeliveryMode: rabbitmq.Persistent,
			MessageId:    uuid.New().String(),
			Timestamp:    time.Now(),
			Body:         body,
			Type:         p.messageTypeName,
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
func (p *amqpPublisherClient) Close() {
	p.client.Close()
}

func (p *amqpPublisherClient) Configure(opts ...Option) AmqpEventPublisher {
	for _, opt := range opts {
		opt(p)
	}

	return p
}
