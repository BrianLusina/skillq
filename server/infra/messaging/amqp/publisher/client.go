package amqppublisher

import (
	"context"
	"fmt"

	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/messaging"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
	"github.com/pkg/errors"
	rabbitmq "github.com/rabbitmq/amqp091-go"
)

// amqpPublisherClient handles defines the methods used to handle publication of messages to a topic on a broker
type amqpPublisherClient struct {
	client             *amqp.AmqpClient
	exchangeName       string
	exchangeKind       string
	exchangeDurable    bool
	exchangeAutoDelete bool
	exchangeInternal   bool
	exchangeNoWait     bool
	exchangeArgs       map[string]any
	bindingKey         string
	messageTypeName    string
	publishMandatory   bool
	publishImmediate   bool
	logger             logger.Logger
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
	body, err := message.PayloadToBytes()
	if err != nil {
		p.logger.Errorf("Failed to parse message: %v", err)
		return errors.Wrapf(err, "failed to parse message event")
	}

	amqpChan, err := p.client.AmqpConn.Channel()
	if err != nil {
		p.logger.Errorf("Failed to open channel: %v", err)
		return fmt.Errorf("failed to open a channel: %w", err)
	}

	defer func() {
		err := amqpChan.Close()
		if err != nil {
			p.logger.Errorf("Failed to close channel with error %v", err)
		}
	}()

	p.logger.Infof("Publishing message %v to exchange: %s, with routingKey: %s", message, p.exchangeName, p.bindingKey)

	err = amqpChan.PublishWithContext(
		ctx,
		p.exchangeName,
		p.bindingKey,
		p.publishMandatory,
		p.publishImmediate,
		rabbitmq.Publishing{
			ContentType:  message.ContentType,
			DeliveryMode: rabbitmq.Persistent,
			MessageId:    message.ID,
			Timestamp:    message.Timestamp,
			Body:         body,
			Type:         message.Topic,
		},
	)
	if err != nil {
		p.logger.Errorf("Failed to publish message to exchange %s with error: %v", p.exchangeName, err)
		return errors.Wrapf(err, "failed to publish message: %v", err)
	}

	p.logger.Infof("Successfully published message %v to exchange: %s, with routingKey: %s", message, p.exchangeName, p.bindingKey)

	return nil
}

// Close closes connection to a broker
func (p *amqpPublisherClient) Close() error {
	if err := p.client.Close(); err != nil {
		p.logger.Errorf("Failed to close client connection with error %v", err)
		return err
	}
	return nil
}

func (p *amqpPublisherClient) Configure(opts ...Option) AmqpEventPublisher {
	for _, opt := range opts {
		opt(p)
	}

	return p
}
