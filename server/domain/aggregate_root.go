package sharedkernel

import "github.com/BrianLusina/skillq/server/domain/entity"

// AggregateRoot represents an Aggregate Entity in the system
type AggregateRoot struct {
	entity.Entity
	domainEvents []DomainEvent
}

// NewAggregateRoot creates a new Aggregate root with a given entity and the events
func NewAggregateRoot(e entity.Entity, events []DomainEvent) *AggregateRoot {
	return &AggregateRoot{
		Entity:       e,
		domainEvents: events,
	}
}

// ApplyDomain applies a domain event to an aggregate root
func (ar *AggregateRoot) ApplyDomain(e DomainEvent) {
	ar.domainEvents = append(ar.domainEvents, e)
}

// DomainEvents retrieves the domain events
func (ar *AggregateRoot) DomainEvents() []DomainEvent {
	return ar.domainEvents
}
