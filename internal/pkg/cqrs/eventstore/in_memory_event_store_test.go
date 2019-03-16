package eventstore_test

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testdata"
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

func TestInMemoryEventStoreLoadEventsFor(t *testing.T) {
	t.Run("ItLoadsEventsForTheGivenAggregate", func(t *testing.T) {
		// assert
		ID := testdata.StringIdentifier("TestAgg")
		es := eventstore.NewInInMemoryEventStore()

		want := []domain.DomainEvent{testdata.SomethingHappened{}}

		err := es.StoreEventsFor(ID, 1, want)
		assert.Ok(t, err)

		// act
		got, err := es.LoadEventsFor(ID)

		// arrange
		assert.Ok(t, err)
		assert.Equals(t, want, got)
	})
}
