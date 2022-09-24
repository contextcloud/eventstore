package groups

import (
	"context"
	"net/http"

	"github.com/contextcloud/eventstore/es"
	"github.com/contextcloud/eventstore/es/providers"
	"github.com/contextcloud/eventstore/examples/groups/aggregates"
	"github.com/contextcloud/eventstore/examples/groups/commands"
	"github.com/contextcloud/eventstore/examples/groups/events"
	"github.com/contextcloud/eventstore/examples/groups/sagas"
	"github.com/contextcloud/eventstore/pkg/db"
	"github.com/contextcloud/eventstore/pkg/pub"
	"github.com/contextcloud/graceful/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/riandyrn/otelchi"
)

func NewHandler(ctx context.Context, cfg *config.Config) (http.Handler, error) {
	pCfg := &providers.Config{}
	if err := cfg.Parse(pCfg); err != nil {
		return nil, err
	}

	gormDb, err := db.Open(&ourCfg.Db)
	if err != nil {
		return nil, err
	}

	conn, err := local.NewConn(gormDb)
	if err != nil {
		return nil, err
	}

	gpub, err := gstream.Open(&ourCfg.Streamer)
	if err != nil {
		return nil, err
	}

	streamer, err := pub.NewStreamer(gpub)
	if err != nil {
		return nil, err
	}

	esCfg, err := es.NewConfig(
		pCfg,
		&aggregates.Group{},
		sagas.NewUserSaga(),
		es.NewAggregateConfig(
			&aggregates.Community{},
			es.EntityDisableProject(),
			es.EntitySnapshotEvery(1),
			es.EntityEventTypes(
				&events.CommunityCreated{},
				&events.CommunityDeleted{},
				&events.CommunityStaffAdded{},
			),
			&commands.CommunityNewCommand{},
			&commands.CommunityDeleteCommand{},
		),
	)
	if err != nil {
		return nil, err
	}

	cli, err := es.NewClient(esCfg, conn, streamer)
	if err != nil {
		return nil, err
	}

	if err := cli.Initialize(ctx); err != nil {
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(otelchi.Middleware("server", otelchi.WithChiRoutes(r)))
	r.Use(es.CreateUnit(cli))
	r.Use(middleware.Logger)
	r.Post("/commands/newcommunity", es.NewCommander[*commands.CommunityNewCommand]())
	r.Post("/commands/creategroup", es.NewCommander[*commands.CreateGroup]())

	return r, nil
}
