package eventbus_test

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/eventbus"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/testdata/mock"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// ensure that EventBus implements domain.EventPublisher interface.
var _ domain.EventPublisher = (*eventbus.InMemoryEventBus)(nil)

func TestNewInMemoryEventBus(t *testing.T) {
	t.Run("ItCreatesNewInstance", func(t *testing.T) {
		assert.True(t, eventbus.NewInMemoryEventBus() != nil)
	})
}

func TestInMemoryEventBus_Publish(t *testing.T) {
	t.Run("ItPublishesEvent", func(t *testing.T) {
		// arrange
		eventHandler := &mock.EventHandlerMock{
			Err: mock.ErrCannotHandleEvent,
		}

		b := eventbus.NewInMemoryEventBus()
		b.Register(eventHandler)

		// act
		err := b.Publish(mock.SomethingHappened{}, mock.SomethingElseHappened{})

		// assert
		assert.Equals(t, mock.ErrCannotHandleEvent, err)
	})

	t.Run("ItFailsIfItCannotHandleAnEvent", func(t *testing.T) {
		// arrange
		want := []domain.DomainEvent{mock.SomethingHappened{}, mock.SomethingElseHappened{}}
		eventHandler := &mock.EventHandlerMock{}

		b := eventbus.NewInMemoryEventBus()
		b.Register(eventHandler)

		// act
		err := b.Publish(mock.SomethingHappened{}, mock.SomethingElseHappened{})

		// assert
		assert.Ok(t, err)
		assert.Equals(t, want, eventHandler.Happened)
	})
}
