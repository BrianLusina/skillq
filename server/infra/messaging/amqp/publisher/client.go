package amqppublisher

import (
	"context"
	"fmt"
	"time"

	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
	"github.com/google/uuid"
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
	p.logger.Infof("Publishing message Exchange: %s, RoutingKey: %s", p.exchangeName, p.bindingKey)

	if err := p.client.AmqpChan.PublishWithContext(
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
			Type:         p.messageTypeName, //"barista.ordered",
		},
	); err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

// CloseChan closes connection to a broker
func (p *amqpPublisherClient) CloseChan() {
	p.client.CloseChan()
}

func (p *amqpPublisherClient) Configure(opts ...Option) AmqpEventPublisher {
	for _, opt := range opts {
		opt(p)
	}

	return p
}
