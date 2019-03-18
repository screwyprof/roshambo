package eventhandler_test

import (
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/testdata/mock"
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/eventhandler"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// ensure that event handler implements domain.EventHandler interface.
var _ domain.EventHandler = (*eventhandler.Static)(nil)

func TestNewStatic(t *testing.T) {
	t.Run("ItCreatesNewInstance", func(t *testing.T) {
		assert.True(t, eventhandler.NewStatic() != nil)
	})
}

func TestStaticHandle(t *testing.T) {
	t.Run("ItHandlesTheGivenEvent", func(t *testing.T) {
		// arrange
		eh := &mock.TestEventHandler{}

		s := eventhandler.NewStatic()
		s.RegisterHandler("OnSomethingHappened", func(e domain.DomainEvent) error {
			return eh.OnSomethingHappened(e.(mock.SomethingHappened))
		})

		// act
		err := s.Handle(mock.SomethingHappened{})

		// assert
		assert.Ok(t, err)
		assert.Equals(t, "test", eh.SomethingHappened)
	})

	t.Run("ItFailsIfEventHandlerIsNotRegistered", func(t *testing.T) {
		// arrange
		s := eventhandler.NewStatic()

		// act
		err := s.Handle(mock.SomethingElseHappened{})

		// assert
		assert.Equals(t, mock.ErrEventHandlerNotFound, err)
	})

	t.Run("ItFailsIfEventHandlerReturnsAnError", func(t *testing.T) {
		// arrange
		eh := &mock.TestEventHandler{}

		s := eventhandler.NewStatic()
		s.RegisterHandler("OnSomethingElseHappened", func(e domain.DomainEvent) error {
			return eh.OnSomethingElseHappened(e.(mock.SomethingElseHappened))
		})

		// act
		err := s.Handle(mock.SomethingElseHappened{})

		// assert
		assert.Equals(t, mock.ErrCannotHandleEvent, err)
	})
}

func TestStaticSubscribedTo(t *testing.T) {
	t.Run("ItReturnersTheEventsItSubscribedTo", func(t *testing.T) {
		// arrange
		eh := &mock.TestEventHandler{}

		s := eventhandler.NewStatic()
		s.RegisterHandler("OnSomethingHappened", func(e domain.DomainEvent) error {
			return eh.OnSomethingHappened(e.(mock.SomethingHappened))
		})
		s.RegisterHandler("OnSomethingElseHappened", func(e domain.DomainEvent) error {
			return eh.OnSomethingElseHappened(e.(mock.SomethingElseHappened))
		})

		want := domain.MatchAnyEventOf("SomethingHappened", "SomethingElseHappened")

		// act
		got := s.SubscribedTo()

		// assert
		assert.Equals(t, want, got)
	})
}
