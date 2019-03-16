package dispatcher_test

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testdata"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testfixture"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/dispatcher"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// ensure that Dispatcher  implements domain.CommandHandler interface.
var _ domain.CommandHandler = (*dispatcher.Dispatcher)(nil)

func TestNewDispatcher(t *testing.T) {
	t.Run("ItPanicsIfEventStoreIsNotGiven", func(t *testing.T) {
		factory := func() {
			dispatcher.NewDispatcher(nil, nil)
		}
		assert.Panic(t, factory)
	})

	t.Run("ItPanicsIfAggregateFactoryIsNotGiven", func(t *testing.T) {
		factory := func() {
			dispatcher.NewDispatcher(nil, nil)
		}
		assert.Panic(t, factory)
	})
}

func TestNewDispatcherHandle(t *testing.T) {
	t.Run("ItCreatesAggregate", func(t *testing.T) {
		// arrange
		aggID := testdata.StringIdentifier("TestAgg")
		agg := aggregate.NewBase(testdata.NewTestAggregate(aggID), nil, nil)

		eventStore := createEventStoreMock(t)
		aggFactory := createAggFactoryMock(t, agg)

		d := dispatcher.NewDispatcher(eventStore, aggFactory)

		// act
		_, err := d.Handle(testdata.MakeSomethingHappen{AggID: aggID})

		// assert
		assert.Ok(t, err)
	})

	t.Run("ItLoadsEventsForAggregate", func(t *testing.T) {
		// arrange
		aggID := testdata.StringIdentifier("TestAgg")
		agg := aggregate.NewBase(testdata.NewTestAggregate(aggID), nil, nil)

		eventStore := createEventStoreMock(t)
		aggFactory := createAggFactoryMock(t, agg)

		d := dispatcher.NewDispatcher(eventStore, aggFactory)

		// act
		_, err := d.Handle(testdata.MakeSomethingHappen{AggID: aggID})

		// assert
		assert.Ok(t, err)
	})

	t.Run("ItFailsIfItCannotLoadEventsForAggregate", func(t *testing.T) {
		// arrange
		aggID := testdata.StringIdentifier("TestAgg")
		agg := aggregate.NewBase(testdata.NewTestAggregate(aggID), nil, nil)

		eventStore := createEventStoreMockWithError(t, testfixture.ErrEventStoreCannotLoadEvents)
		aggFactory := createAggFactoryMock(t, agg)

		d := dispatcher.NewDispatcher(eventStore, aggFactory)

		// act
		_, err := d.Handle(testdata.MakeSomethingHappen{AggID: aggID})

		// assert
		assert.Equals(t, testfixture.ErrEventStoreCannotLoadEvents, err)
	})

	t.Run("ItAppliesTheLoadedEventsToTheCreatedAggregate", func(t *testing.T) {
		// arrange
		aggID := testdata.StringIdentifier("TestAgg")
		agg := aggregate.NewBase(testdata.NewTestAggregate(aggID), nil, nil)

		eventStore := createEventStoreMock(t)
		aggFactory := createAggFactoryMock(t, agg)

		d := dispatcher.NewDispatcher(eventStore, aggFactory)

		// act
		_, err := d.Handle(testdata.MakeSomethingHappen{AggID: aggID})

		// assert
		assert.Ok(t, err)
		assert.Equals(t, agg.Version(), 1)
	})

	t.Run("ItFailsIfItCannotApplyEvents", func(t *testing.T) {
		// arrange
		aggID := testdata.StringIdentifier("TestAgg")
		agg := aggregate.NewBase(
			testdata.NewTestAggregate(aggID),
			nil, aggregate.NewStaticEventApplier())

		eventStore := createEventStoreMock(t)
		aggFactory := createAggFactoryMock(t, agg)

		d := dispatcher.NewDispatcher(eventStore, aggFactory)

		// act
		_, err := d.Handle(testdata.MakeSomethingHappen{AggID: aggID})

		// assert
		assert.Equals(t, testdata.ErrOnSomethingHappenedApplierNotFound, err)
	})
}

func createAggFactoryMock(t *testing.T, agg *aggregate.Base) *testfixture.AggregateFactoryMock {
	aggFactory := &testfixture.AggregateFactoryMock{
		Creator: func(aggregateType string, ID domain.Identifier) domain.AdvancedAggregate {
			assert.Equals(t, agg.AggregateID(), ID)
			assert.Equals(t, agg.AggregateType(), aggregateType)
			return agg
		},
	}
	return aggFactory
}

func createEventStoreMock(t *testing.T) *testfixture.EventStoreMock {
	eventStore := &testfixture.EventStoreMock{
		Loader: func(aggregateID domain.Identifier) ([]domain.DomainEvent, error) {
			assert.True(t, true)
			return []domain.DomainEvent{testdata.SomethingHappened{}}, nil
		},
	}
	return eventStore
}

func createEventStoreMockWithError(t *testing.T, err error) *testfixture.EventStoreMock {
	eventStore := &testfixture.EventStoreMock{
		Loader: func(aggregateID domain.Identifier) ([]domain.DomainEvent, error) {
			assert.True(t, true)
			return nil, err
		},
	}
	return eventStore
}
