package es

import (
	"testing"
)

func Test_SagaHandle(t *testing.T) {
	saga := &demoSaga{}
	handles := NewSagaHandles(saga)

	if len(handles) != 1 {
		t.Errorf("expected 1 handle, got %d", len(handles))
	}
}
