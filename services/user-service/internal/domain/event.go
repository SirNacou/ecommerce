package domain

import "time"

// DomainEvent represents a state change within an aggregate
type DomainEvent interface {
	EventName() string
	OccurredAt() time.Time
}

// AggregateRoot is embedded by domain entities to collect domain events
type AggregateRoot struct {
	events []DomainEvent
}

func (a *AggregateRoot) AddEvent(event DomainEvent) {
	a.events = append(a.events, event)
}

func (a *AggregateRoot) PopEvents() []DomainEvent {
	events := a.events
	a.events = nil // Clear accumulated events
	return events
}
