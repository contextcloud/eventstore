package gpub

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/contextcloud/eventstore/es"
	"github.com/contextcloud/eventstore/pkg/gcppubsub"
	"go.opentelemetry.io/otel"
)

func wrapped(callback func(context.Context, []byte) error) func(ctx context.Context, msg *pubsub.Message) {
	return func(ctx context.Context, msg *pubsub.Message) {
		if err := callback(ctx, msg.Data); err != nil {
			msg.Nack()
			return
		}
		msg.Ack()
	}
}

type streamer struct {
	service string
	p       *gcppubsub.Pub

	started bool
}

func (s *streamer) Start(ctx context.Context, callback es.EventCallback) error {
	pctx, span := otel.Tracer("gpub").Start(ctx, "Start")
	defer span.End()

	if callback == nil {
		return fmt.Errorf("callback is required")
	}

	sub, err := s.p.Subscription(pctx, s.service)
	if err != nil {
		return err
	}

	sub.ReceiveSettings.MaxOutstandingMessages = 100

	handle := func(ctx context.Context, data []byte) error {
		pctx, span := otel.Tracer("gpub").Start(ctx, "Handle")
		defer span.End()

		evt, err := es.UnmarshalEvent(pctx, data)
		if errors.Is(err, es.ErrNotFound) {
			return nil
		}
		if err != nil {
			return err
		}

		return callback(pctx, evt)
	}

	go func() {
		ctx := context.Background()
		pctx, span := otel.Tracer("gpub").Start(ctx, "Receive")
		defer span.End()

		if err := sub.Receive(pctx, wrapped(handle)); err != nil {
			// no sure what to do here yet
			panic(err)
		}
	}()

	s.started = true
	return nil
}

func (s *streamer) Publish(ctx context.Context, evts ...*es.Event) error {
	pctx, span := otel.Tracer("gpub").Start(ctx, "Publish")
	defer span.End()

	if !s.started {
		return fmt.Errorf("streamer is not started")
	}

	messages := make([]*pubsub.Message, len(evts))
	for i, evt := range evts {
		orderingKey := fmt.Sprintf("%s:%s:%s:%d", evt.Namespace, evt.AggregateId.String(), evt.AggregateType, evt.Version)
		data, err := es.MarshalEvent(ctx, evt)
		if err != nil {
			return err
		}

		msg := &pubsub.Message{
			Data:        data,
			OrderingKey: orderingKey,
		}
		messages[i] = msg
	}

	_, err := s.p.Publish(pctx, messages...)
	if err != nil {
		// todo add some logging
		return err
	}

	return nil
}

func (s *streamer) Close(ctx context.Context) error {
	_, span := otel.Tracer("gpub").Start(ctx, "Close")
	defer span.End()

	return s.p.Close()
}

func NewStreamer(service string, p *gcppubsub.Pub) (es.Streamer, error) {
	return &streamer{
		service: service,
		p:       p,
	}, nil
}
