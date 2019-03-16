package dispatcher_test

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	fMock "github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/mock"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testdata"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/dispatcher"
	. "github.com/screwyprof/roshambo/internal/pkg/cqrs/dispatcher/testdata/fixture"
	esMock "github.com/screwyprof/roshambo/internal/pkg/cqrs/eventstore/mock"

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

		eventStore := createEventStoreMock(nil, nil)
		aggFactory := createAggFactoryMock(agg)

		d := dispatcher.NewDispatcher(eventStore, aggFactory)

		// act
		_, err := d.Handle(testdata.MakeSomethingHappen{AggID: aggID})

		// assert
		assert.Ok(t, err)
	})

	t.Run("ItLoadsEventsForAggregate", func(t *testing.T) {
		aggID := testdata.StringIdentifier("TestAgg")
		d := createDispatcher(aggID, []domain.DomainEvent{testdata.SomethingElseHappened{}}, nil)

		Test(t)(
			Given(d),
			When(testdata.MakeSomethingHappen{AggID: aggID}),
			Then(),
		)
	})

	t.Run("ItFailsIfItCannotLoadEventsForAggregate", func(t *testing.T) {
		aggID := testdata.StringIdentifier("TestAgg")
		d := createDispatcher(aggID, nil, esMock.ErrEventStoreCannotLoadEvents)

		Test(t)(
			Given(d),
			When(testdata.MakeSomethingHappen{AggID: aggID}),
			ThenFailWith(esMock.ErrEventStoreCannotLoadEvents),
		)
	})

	t.Run("ItAppliesTheLoadedEventsToTheCreatedAggregate", func(t *testing.T) {
		// arrange
		aggID := testdata.StringIdentifier("TestAgg")
		agg := aggregate.NewBase(testdata.NewTestAggregate(aggID), nil, nil)

		eventStore := createEventStoreMock([]domain.DomainEvent{testdata.SomethingElseHappened{}}, nil)
		aggFactory := createAggFactoryMock(agg)

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

		eventStore := createEventStoreMock(nil, nil)
		aggFactory := createAggFactoryMock(agg)

		d := dispatcher.NewDispatcher(eventStore, aggFactory)

		Test(t)(
			Given(d),
			When(testdata.MakeSomethingHappen{AggID: aggID}),
			ThenFailWith(testdata.ErrOnSomethingHappenedApplierNotFound),
		)
	})

	t.Run("ItFailsIfAggregateCannotHandleTheGivenCommand", func(t *testing.T) {
		// arrange
		aggID := testdata.StringIdentifier("TestAgg")
		d := createDispatcher(aggID, []domain.DomainEvent{testdata.SomethingHappened{}}, nil)

		Test(t)(
			Given(d),
			When(testdata.MakeSomethingHappen{AggID: aggID}),
			ThenFailWith(testdata.ErrItCanHappenOnceOnly),
		)
	})
}

func createDispatcher(aggID testdata.StringIdentifier, loadedEvents []domain.DomainEvent, eventStoreErr error) *dispatcher.Dispatcher {
	agg := aggregate.NewBase(testdata.NewTestAggregate(aggID), nil, nil)
	eventStore := createEventStoreMock(loadedEvents, eventStoreErr)
	aggFactory := createAggFactoryMock(agg)
	return dispatcher.NewDispatcher(eventStore, aggFactory)
}

func createAggFactoryMock(agg *aggregate.Base) *fMock.AggregateFactoryMock {
	aggFactory := &fMock.AggregateFactoryMock{
		Creator: func(aggregateType string, ID domain.Identifier) domain.AdvancedAggregate {
			return agg
		},
	}
	return aggFactory
}

func createEventStoreMock(want []domain.DomainEvent, err error) *esMock.EventStoreMock {
	eventStore := &esMock.EventStoreMock{
		Loader: func(aggregateID domain.Identifier) ([]domain.DomainEvent, error) {
			return want, err
		},
	}
	return eventStore
}
