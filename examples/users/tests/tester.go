package tests

import (
	"context"

	"github.com/contextcloud/eventstore/es"
	"github.com/contextcloud/eventstore/examples/users/data"
	"github.com/contextcloud/goutils/xgorm"

	_ "github.com/contextcloud/eventstore/es/providers/data/pg"
	_ "github.com/contextcloud/eventstore/es/providers/stream/apub"
	_ "github.com/contextcloud/eventstore/es/providers/stream/gpub"
	_ "github.com/contextcloud/eventstore/es/providers/stream/noop"
	_ "github.com/contextcloud/eventstore/es/providers/stream/npub"
)

type Tester interface {
	Client() es.Client
}

type tester struct {
	client es.Client
}

func (h *tester) Client() es.Client {
	return h.client
}

func NewTester() (Tester, error) {
	pcfg := &es.ProviderConfig{
		Service: "users",
		Version: "v1",
		Data: es.DataConfig{
			Type: "pg",
			Pg: &xgorm.DbConfig{
				Host:     "localhost",
				Port:     5432,
				Username: "es",
				Password: "es",
				Database: "eventstore",

				MaxIdleConns: 10,
				MaxOpenConns: 10,
			},
			Reset: true,
		},
		Stream: es.StreamConfig{
			Type: "noop",
			PubSub: &es.GcpPubSubConfig{
				ProjectId: "nordic-gaming",
				TopicId:   "test_topic",
			},
			Nats: &es.NatsConfig{
				Url:     "nats://localhost:4222",
				Subject: "test",
			},
			AWS: &es.AwsSnsConfig{
				Region:   "us-east-1",
				TopicArn: "arn:aws:sns:us-east-1:888821167166:deployment.fifo",
				Profile:  "Development",
			},
		},
	}

	ctx := context.Background()
	cli, err := data.NewClient(ctx, pcfg)
	if err != nil {
		return nil, err
	}

	return &tester{
		client: cli,
	}, nil
}
