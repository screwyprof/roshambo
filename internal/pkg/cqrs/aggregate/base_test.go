package aggregate_test

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testdata"

	"github.com/screwyprof/roshambo/pkg/domain"
)

func TestNewBase(t *testing.T) {
	t.Run("ItPanicsIfThePureAggregateIsNotGiven", func(t *testing.T) {
		factory := func() {
			aggregate.NewBase(nil, nil, nil)
		}
		assert.Panic(t, factory)
	})

	t.Run("ItReturnsAggregateID", func(t *testing.T) {
		ID := testdata.StringIdentifier("TestAgg1")
		pureAgg := testdata.NewTestAggregate(ID)
		agg := aggregate.NewBase(pureAgg, nil, nil)

		assert.Equals(t, ID, agg.AggregateID())
	})
}

func TestBaseHandle(t *testing.T) {
	t.Run("ItHandlesTheGivenCommandAndAppliesEventsIfTheHandlerExists", func(t *testing.T) {
		// arrange
		agg := createTestAgg()
		expectedEvents := []domain.DomainEvent{testdata.SomethingHappened{}}

		// act
		events, err := agg.Handle(testdata.MakeSomethingHappen{})

		// assert
		assert.Ok(t, err)
		assert.Equals(t, expectedEvents, events)
	})

	t.Run("ItReturnsAnErrorIfTheHandlerIsNotFound", func(t *testing.T) {
		// arrange
		agg := createTestAggWithDefaultCommandHandlerAndEventApplier()

		// act
		_, err := agg.Handle(testdata.MakeSomethingHappen{})

		// assert
		assert.Equals(t, testdata.ErrMakeSomethingHandlerNotFound, err)
	})

	t.Run("ItReturnsAnErrorIfTheEventAppliersNotFound", func(t *testing.T) {
		// arrange
		agg := createTestAggWithCustomCommandHandler()

		// act
		_, err := agg.Handle(testdata.MakeSomethingHappen{})

		// assert
		assert.Equals(t, testdata.ErrOnSomethingHappenedApplierNotFound, err)
	})
}

func createTestAggWithDefaultCommandHandlerAndEventApplier() *aggregate.Base {
	ID := testdata.StringIdentifier("TestAgg1")
	pureAgg := testdata.NewTestAggregate(ID)
	agg := aggregate.NewBase(pureAgg, nil, nil)
	return agg
}

func createTestAggWithCustomCommandHandler() *aggregate.Base {
	ID := testdata.StringIdentifier("TestAgg1")
	pureAgg := testdata.NewTestAggregate(ID)
	commandHandler := createCommandHandler(pureAgg)
	agg := aggregate.NewBase(pureAgg, commandHandler, nil)
	return agg
}

func createTestAgg() *aggregate.Base {
	ID := testdata.StringIdentifier("TestAgg1")
	pureAgg := testdata.NewTestAggregate(ID)

	commandHandler := createCommandHandler(pureAgg)
	eventApplier := createEventApplier(pureAgg)

	return aggregate.NewBase(pureAgg, commandHandler, eventApplier)
}

func createEventApplier(pureAgg *testdata.TestAggregate) *aggregate.StaticEventApplier {
	eventApplier := aggregate.NewStaticEventApplier()
	eventApplier.RegisterApplier(
		"OnSomethingHappened",
		func(e domain.DomainEvent) {
			pureAgg.OnSomethingHappened(e.(testdata.SomethingHappened))
		},
	)
	return eventApplier
}

func createCommandHandler(pureAgg *testdata.TestAggregate) *aggregate.StaticCommandHandler {
	commandHandler := aggregate.NewStaticCommandHandler()
	commandHandler.RegisterHandler(
		"MakeSomethingHappen",
		func(c domain.Command) ([]domain.DomainEvent, error) {
			return pureAgg.MakeSomethingHappen(c.(testdata.MakeSomethingHappen))
		},
	)
	return commandHandler
}
