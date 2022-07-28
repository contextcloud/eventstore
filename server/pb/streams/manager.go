package streams

import (
	"context"
	"sync"

	"github.com/contextcloud/eventstore/server/pb/store"
	"gorm.io/gorm"
)

type Manager interface {
	NewStream(stream store.Store_EventStreamServer) Sender
	DeleteSender(sender Sender)
	Register(serviceName string, eventTypes []string, stream store.Store_EventStreamServer) error
	Stop() error
}

type manager struct {
	mux sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc
	gormDb *gorm.DB

	senders map[Sender]bool
}

func (m *manager) loop() {
	for {
		select {
		case <-m.ctx.Done():
			return
		}
	}
}

func (m *manager) NewStream(stream store.Store_EventStreamServer) Sender {
	m.mux.Lock()
	defer m.mux.Unlock()

	sender := NewSender(m, stream)
	m.senders[sender] = true
	return sender
}
func (m *manager) DeleteSender(sender Sender) {
	m.mux.Lock()
	defer m.mux.Unlock()

	delete(m.senders, sender)
}
func (m *manager) Register(serviceName string, eventTypes []string, stream store.Store_EventStreamServer) error {
	return nil
}
func (m *manager) Stop() error {
	m.cancel()
	return nil
}
func (m *manager) Start() error {
	go m.loop()
	return nil
}

func NewManager(gormDb *gorm.DB) (Manager, error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	return &manager{
		ctx:     ctx,
		cancel:  cancel,
		gormDb:  gormDb,
		senders: make(map[Sender]bool),
	}, nil
}