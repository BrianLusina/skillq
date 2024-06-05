package di

import (
	"github.com/BrianLusina/skillq/server/app/internal/publishers"
	"github.com/google/wire"
)

var SendEmailEventPublisherSet = wire.NewSet(publishers.NewSendEmailEventPublisher)

var StoreImageEventPublisherSet = wire.NewSet(publishers.NewStoreImageEventPublisher)
