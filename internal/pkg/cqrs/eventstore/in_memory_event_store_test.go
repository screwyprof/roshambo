package eventstore_test

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/eventstore"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// ensure that event store implements domain.EventStore interface.
var _ domain.EventStore = (*eventstore.InMemoryEventStore)(nil)

func TestNewInInMemoryEventStore(t *testing.T) {
	t.Run("ItCreatesEventStore", func(t *testing.T) {
		es := eventstore.NewInInMemoryEventStore()
		assert.True(t, es != nil)
	})
}
