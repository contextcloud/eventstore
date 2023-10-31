package events

import (
	"github.com/contextcloud/eventstore/es"
	"github.com/contextcloud/eventstore/es/types"
	"github.com/contextcloud/eventstore/examples/users/models"
)

type UserCreated struct {
	Username string
	Password string
}
type UserDeleted struct {
	Deleted bool
}

type EmailAdded struct {
	Email string
}

type ConnectionAdded struct {
	Connections types.SliceItem[models.Connection]
}

type ConnectionUpdated struct {
	Connections types.SliceItem[models.ConnectionUpdate]
}

type GroupAdded struct {
	es.BaseEvent `es:"publish"`

	Groups types.SliceItem[models.Group]
}
