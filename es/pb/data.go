package pb

import (
	"context"
	"encoding/json"

	"github.com/contextcloud/eventstore/es"
	"github.com/contextcloud/eventstore/es/filters"
	"github.com/contextcloud/eventstore/es/pb/store"
)

type data struct {
	storeClient store.StoreClient

	transactionId *string
}

func (d *data) Begin(ctx context.Context) (es.Tx, error) {
	if d.transactionId == nil {
		req := &store.NewTxRequest{}
		resp, err := d.storeClient.NewTx(ctx, req)
		if err != nil {
			return nil, err
		}
		d.transactionId = &resp.TransactionID
	}

	return newTransaction(d.storeClient, *d.transactionId), nil
}
func (d *data) LoadSnapshot(ctx context.Context, serviceName string, aggregateName string, namespace string, id string, out es.SourcedAggregate) error {
	return nil
}
func (d *data) GetEventDatas(ctx context.Context, serviceName string, aggregateName string, namespace string, id string, fromVersion int) ([]json.RawMessage, error) {
	f := int64(fromVersion)
	req := &store.EventsRequest{
		TransactionID: d.transactionId,
		ServiceName:   &serviceName,
		AggregateType: &aggregateName,
		AggregateId:   &id,
		Namespace:     &namespace,
		FromVersion:   &f,
	}
	resp, err := d.storeClient.Events(ctx, req)
	if err != nil {
		return nil, err
	}

	var datas []json.RawMessage
	for _, event := range resp.Events {
		datas = append(datas, event.Data)
	}
	return datas, nil
}
func (d *data) SaveEvents(ctx context.Context, events []es.Event) error {

	return nil
}
func (d *data) SaveEntity(ctx context.Context, entity es.Entity) error {
	return nil
}

func (d *data) Load(ctx context.Context, serviceName string, aggregateName string, namespace string, id string, out interface{}) error {
	return nil
}
func (d *data) Find(ctx context.Context, serviceName string, aggregateName string, namespace string, filter filters.Filter, out interface{}) error {
	return nil
}
func (d *data) Count(ctx context.Context, serviceName string, aggregateName string, namespace string, filter filters.Filter) (int, error) {
	return 0, nil
}

func newData(storeClient store.StoreClient) (es.Data, error) {
	return &data{
		storeClient: storeClient,
	}, nil
}
