package eventbus

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Bus struct {
	nc *nats.Conn
	js jetstream.JetStream
}

func New(url string) (*Bus, error) {
	nc, err := nats.Connect(url,
		nats.Timeout(5*time.Second),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("nats jetstream: %w", err)
	}

	return &Bus{nc: nc, js: js}, nil
}

func (b *Bus) Close() {
	if b.nc != nil {
		b.nc.Close()
	}
}

func (b *Bus) EnsureStream(ctx context.Context, stream, subjectPrefix string) error {
	_, err := b.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     stream,
		Subjects: []string{subjectPrefix + ".*"},
		Storage:  jetstream.FileStorage,
	})
	return err
}

func (b *Bus) Publish(ctx context.Context, subject string, data []byte) error {
	_, err := b.js.Publish(ctx, subject, data)
	return err
}

type MessageHandler func(ctx context.Context, data []byte) error

// Subscribe creates a durable consumer on the given stream and continuously
// delivers matching messages to the handler. The durable consumer is created
// explicitly and is never deleted on shutdown, so it resumes from its last
// acknowledged position across restarts and messages published while the
// consumer is offline are still delivered.
func (b *Bus) Subscribe(ctx context.Context, stream, subject, durable string, handler MessageHandler) error {
	consumer, err := b.js.CreateOrUpdateConsumer(ctx, stream, jetstream.ConsumerConfig{
		Durable:       durable,
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy,
		FilterSubject: subject,
	})
	if err != nil {
		return fmt.Errorf("nats create consumer %s: %w", durable, err)
	}

	consumeCtx, err := consumer.Consume(func(msg jetstream.Msg) {
		if err := handler(ctx, msg.Data()); err != nil {
			// Negative ack so NATS redelivers and the event is retried.
			_ = msg.Nak()
			return
		}
		_ = msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("nats consume %s: %w", durable, err)
	}

	go func() {
		<-ctx.Done()
		consumeCtx.Drain()
	}()

	return nil
}
