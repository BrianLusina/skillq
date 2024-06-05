package di

import (
	"github.com/BrianLusina/skillq/server/infra/clients/email"
	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
	amqpconsumer "github.com/BrianLusina/skillq/server/infra/messaging/amqp/consumer"
	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
	"github.com/BrianLusina/skillq/server/infra/storage/minio"
	"github.com/google/wire"
)

// Logger
var LoggerSet = wire.NewSet(logger.New)

// AMQP Client provider and AMQP Event publisher set
var AmqpClientSet = wire.NewSet(amqp.NewAmqpClient)
var AmqpEventPublisherSet = wire.NewSet(amqppublisher.NewPublisher)
var AmqpEventConsumerSet = wire.NewSet(amqpconsumer.NewConsumer)

// Storage clients
var StorageMinioClientSet = wire.NewSet(minio.NewClient)

var EmailClientSet = wire.NewSet(email.New)
