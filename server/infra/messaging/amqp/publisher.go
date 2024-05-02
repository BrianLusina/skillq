package amqp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/messaging"
	"github.com/google/uuid"
	rabbitmq "github.com/rabbitmq/amqp091-go"
)

const (
	_retryTimes     = 5
	_backOffSeconds = 2

	_publishMandatory = false
	_publishImmediate = false

	_exchangeName    = "skillq-exchange"
	_bindingKey      = "skillq-routing-key"
	_messageTypeName = "skillq"
)

// AmqpPublisher handles defines the methods used to handle publication of messages to a topic on a broker
type AmqpPublisher struct {
	amqpConn        *rabbitmq.Connection
	amqpChan        *rabbitmq.Channel
	exchangeName    string
	bindingKey      string
	messageTypeName string
	logger          logger.Logger
}

var ErrCannotConnectRabbitMQ = errors.New("cannot connect to rabbit")

// NewPublisher creates a new AMQP Publisher
func NewPublisher(config Config, log logger.Logger, opts ...AmqpPublisherOption) (messaging.Publisher, error) {
	connString := fmt.Sprintf("amqp://%v:%v@%v:%v/", config.Username, config.Password, config.Host, config.Port)
	var (
		amqpConn *rabbitmq.Connection
		counts   int64
	)

	for {
		conn, err := rabbitmq.Dial(connString)
		if err != nil {
			log.Errorf("RabbitMq at %s:%s not ready...\n", config.Host, config.Port)
			counts++
		} else {
			amqpConn = conn
			break
		}

		if counts > _retryTimes {
			log.Fatalf(err)

			return nil, ErrCannotConnectRabbitMQ
		}

		log.Info("Backing off for 2 seconds...")
		time.Sleep(_backOffSeconds * time.Second)

		continue
	}

	log.Info("Connected to RabbitMQ!")

	amqpChan, err := amqpConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	publisher := &AmqpPublisher{
		amqpConn:        amqpConn,
		amqpChan:        amqpChan,
		logger:          log,
		exchangeName:    _exchangeName,
		bindingKey:      _bindingKey,
		messageTypeName: _messageTypeName,
	}

	for _, opt := range opts {
		opt(publisher)
	}

	return publisher, nil
}

// Publish publishes a message to a given topic
func (p *AmqpPublisher) Publish(ctx context.Context, body []byte, contentType string) error {
	p.logger.Infof("Publishing message Exchange: %s, RoutingKey: %s", p.exchangeName, p.bindingKey)

	if err := p.amqpChan.PublishWithContext(
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
func (p *AmqpPublisher) CloseChan() {
	if err := p.amqpChan.Close(); err != nil {
		p.logger.Errorf("Publisher CloseChan: %v", err)
	}
}
