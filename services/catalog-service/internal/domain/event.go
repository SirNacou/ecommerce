package domain

import "time"

type DomainEvent interface {
	EventType() string
	OccurredAt() time.Time
}

type AggregateRoot struct {
	events []DomainEvent
}

func (a *AggregateRoot) RecordEvent(event DomainEvent) {
	a.events = append(a.events, event)
}

func (a *AggregateRoot) PopEvents() []DomainEvent {
	events := a.events
	a.events = nil
	return events
}
