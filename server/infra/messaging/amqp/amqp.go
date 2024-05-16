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
			log.Fatalf("failed to retry: %v", err)

			return nil, ErrCannotConnectRabbitMQ
		}

		log.Info("Backing off for 2 seconds...")
		time.Sleep(_backOffSeconds * time.Second)

		continue
	}

	log.Infof("Connected to RabbitMQ running on host: %s:%s", config.Host, config.Port)

	return &AmqpClient{
		AmqpConn: amqpConn,
		logger:   log,
	}, nil
}

// Close closes connection to a broker
func (p *AmqpClient) Close() {
	if err := p.AmqpConn.Close(); err != nil {
		p.logger.Errorf("Publisher CloseChan: %v", err)
	}
}
