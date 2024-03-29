package sagas

import (
	"context"
	"fmt"

	"github.com/contextcloud/eventstore/examples/users/data/aggregates"
	"github.com/contextcloud/eventstore/examples/users/data/commands"
	"github.com/contextcloud/eventstore/examples/users/data/events"
	"github.com/contextcloud/eventstore/examples/users/helpers"
	"github.com/google/uuid"

	"github.com/contextcloud/eventstore/es"
)

type ConnectionSaga struct {
	es.BaseSaga
}

func (s *ConnectionSaga) HandleConnectionAdded(ctx context.Context, evt *es.Event, data *events.ConnectionAdded) ([]es.Command, error) {
	skip := helpers.GetSkipSaga(ctx)
	if skip {
		return nil, nil
	}

	item := data.Connections.Value

	id := uuid.NewSHA1(evt.AggregateId, []byte(item.UserId))

	q := es.NewQuery[*aggregates.User]()
	all, err := q.Find(ctx, es.Filter{})
	if err != nil {
		return nil, err
	}
	fmt.Printf("all: %+v", all)

	return es.Commands(&commands.CreateExternalUser{
		BaseCommand: es.BaseCommand{
			AggregateId: id,
		},

		Name:     item.Name,
		UserId:   item.UserId,
		Username: item.Username,
	}), nil
}

func NewConnectionSaga() *ConnectionSaga {
	return &ConnectionSaga{
		BaseSaga: es.BaseSaga{},
	}
}
