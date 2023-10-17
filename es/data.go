package es

import (
	"context"

	"github.com/contextcloud/eventstore/es/filters"
	"github.com/google/uuid"
)

type EventSearch struct {
	Namespace     string
	AggregateType string
	AggregateId   uuid.UUID
	FromVersion   int
}

type SnapshotSearch struct {
	Namespace     string
	AggregateType string
	AggregateId   uuid.UUID
	Revision      string
}

type ConnFactory func(ctx context.Context, cfg *ProviderConfig) (Conn, error)

type Conn interface {
	NewData(ctx context.Context) (Data, error)
	Close(ctx context.Context) error
}

type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Data interface {
	Begin(ctx context.Context) (Tx, error)

	LoadSnapshot(ctx context.Context, search SnapshotSearch, out AggregateSourced) error
	SaveSnapshot(ctx context.Context, snapshot *Snapshot) error

	GetEvents(ctx context.Context, mapper EventDataMapper, search EventSearch) ([]*Event, error)
	SaveEvents(ctx context.Context, events []*Event) error
	SaveEntity(ctx context.Context, aggregateName string, entity Entity) error
	DeleteEntity(ctx context.Context, aggregateName string, entity Entity) error
	Truncate(ctx context.Context, aggregateName string) error

	Get(ctx context.Context, aggregateName string, namespace string, id uuid.UUID, out interface{}) error
	Find(ctx context.Context, aggregateName string, namespace string, filter filters.Filter, out interface{}) error
	Count(ctx context.Context, aggregateName string, namespace string, filter filters.Filter) (int, error)
}
