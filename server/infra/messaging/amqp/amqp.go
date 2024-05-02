package amqp

import (
	"fmt"
	"time"

	"github.com/BrianLusina/skillq/server/infra/logger"
	rabbitmq "github.com/rabbitmq/amqp091-go"
)

// AmqpClient defines an interface that Advanced Message Queuing Protocol clients implement
type AmqpClient struct {
	AmqpConn *rabbitmq.Connection
	AmqpChan *rabbitmq.Channel
	logger   logger.Logger
}

// NewAmqpClient creates an AMQP Client that contains a connection to an AMQP broker
func NewAmqpClient(config Config, log logger.Logger) (*AmqpClient, error) {
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

	return &AmqpClient{
		AmqpConn: amqpConn,
		AmqpChan: amqpChan,
		logger:   log,
	}, nil
}

// CloseChan closes connection to a broker
func (p *AmqpClient) CloseChan() {
	if err := p.AmqpChan.Close(); err != nil {
		p.logger.Errorf("Publisher CloseChan: %v", err)
	}
}
