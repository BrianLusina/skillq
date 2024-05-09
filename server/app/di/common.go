package di

import (
	"github.com/BrianLusina/skillq/server/infra/logger"
	"github.com/BrianLusina/skillq/server/infra/messaging/amqp"
	amqppublisher "github.com/BrianLusina/skillq/server/infra/messaging/amqp/publisher"
	"github.com/BrianLusina/skillq/server/infra/storage/minio"
	"github.com/google/wire"
)

// Logger
var LoggerSet = wire.NewSet(logger.New)

// AMQP Client provider and AMQP Event publisher set
var AmqpClientSet = wire.NewSet(amqp.NewAmqpClient)
var AmqpEventPublisherSet = wire.NewSet(amqppublisher.NewPublisher)

// Storage clients
var StorageMinioClientSet = wire.NewSet(minio.NewClient)
