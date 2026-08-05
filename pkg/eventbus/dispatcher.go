package eventbus

import (
	"context"
	"log"
	"time"
)

type OutboxEvent struct {
	ID        string
	EventType string
	Payload   []byte
}

// OutboxReader abstracts the producer's outbox table access.
type OutboxReader interface {
	GetPendingOutboxEvents(ctx context.Context, limit int32) ([]OutboxEvent, error)
	MarkOutboxEventProcessed(ctx context.Context, id string) error
}

type Dispatcher struct {
	bus        *Bus
	reader     OutboxReader
	stream     string
	subject    string
	pollEvery  time.Duration
	batchSize  int32
}

func NewDispatcher(bus *Bus, reader OutboxReader, stream, subject string) *Dispatcher {
	return &Dispatcher{
		bus:       bus,
		reader:    reader,
		stream:    stream,
		subject:   subject,
		pollEvery: 2 * time.Second,
		batchSize: 50,
	}
}

func (d *Dispatcher) WithPollInterval(interval time.Duration) *Dispatcher {
	d.pollEvery = interval
	return d
}

func (d *Dispatcher) WithBatchSize(size int32) *Dispatcher {
	d.batchSize = size
	return d
}

func (d *Dispatcher) Start(ctx context.Context) error {
	if err := d.bus.EnsureStream(ctx, d.stream, d.subject); err != nil {
		return err
	}

	go d.run(ctx)
	return nil
}

func (d *Dispatcher) run(ctx context.Context) {
	ticker := time.NewTicker(d.pollEvery)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := d.dispatch(ctx); err != nil {
				log.Printf("eventbus: outbox dispatch error: %v", err)
			}
		}
	}
}

func (d *Dispatcher) dispatch(ctx context.Context) error {
	events, err := d.reader.GetPendingOutboxEvents(ctx, d.batchSize)
	if err != nil {
		return err
	}

	for _, event := range events {
		subject := d.subject + "." + event.EventType
		if err := d.bus.Publish(ctx, subject, event.Payload); err != nil {
			return err
		}

		if err := d.reader.MarkOutboxEventProcessed(ctx, event.ID); err != nil {
			return err
		}
	}

	return nil
}
