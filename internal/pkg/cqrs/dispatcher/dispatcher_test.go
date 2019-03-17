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
			dispatcher.NewDispatcher(createEventStoreMock(nil, nil, nil), nil)
		}
		assert.Panic(t, factory)
	})
}

func TestNewDispatcherHandle(t *testing.T) {
	t.Run("ItFailsIfItCannotLoadEventsForAggregate", func(t *testing.T) {
		aggID := testdata.StringIdentifier("TestAgg")
		d := createDispatcher(aggID, nil, esMock.ErrEventStoreCannotLoadEvents)

		Test(t)(
			Given(d),
			When(testdata.MakeSomethingHappen{AggID: aggID}),
			ThenFailWith(esMock.ErrEventStoreCannotLoadEvents),
		)
	})

	t.Run("ItFailsIfItCannotApplyEvents", func(t *testing.T) {
		// arrange
		aggID := testdata.StringIdentifier("TestAgg")
		d := createDispatcherWithStaticEventApplier(aggID, []domain.DomainEvent{testdata.SomethingHappened{}}, nil)

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

	t.Run("ItFailsIfItCannotStoreEvents", func(t *testing.T) {
		aggID := testdata.StringIdentifier("TestAgg")
		d := createDispatcherWithStoreEventStore(aggID, nil, esMock.ErrEventStoreCannotStoreEvents)

		Test(t)(
			Given(d),
			When(testdata.MakeSomethingHappen{AggID: aggID}),
			ThenFailWith(esMock.ErrEventStoreCannotStoreEvents),
		)
	})

	t.Run("ItReturnsEvents", func(t *testing.T) {
		aggID := testdata.StringIdentifier("TestAgg")
		d := createDispatcherWithStoreEventStore(aggID, nil, nil)
		Test(t)(
			Given(d),
			When(testdata.MakeSomethingHappen{AggID: aggID}),
			Then(testdata.SomethingHappened{}),
		)
	})
}

func createDispatcher(
	aggID testdata.StringIdentifier, loadedEvents []domain.DomainEvent, eventStoreErr error) *dispatcher.Dispatcher {
	agg := aggregate.NewBase(testdata.NewTestAggregate(aggID), nil, nil)
	eventStore := createEventStoreMock(loadedEvents, eventStoreErr, nil)
	aggFactory := createAggFactoryMock(agg)
	return dispatcher.NewDispatcher(eventStore, aggFactory)
}

func createDispatcherWithStoreEventStore(
	aggID testdata.StringIdentifier, loadedEvents []domain.DomainEvent, eventStoreErr error) *dispatcher.Dispatcher {
	agg := aggregate.NewBase(testdata.NewTestAggregate(aggID), nil, nil)
	eventStore := createEventStoreMock(loadedEvents, nil, eventStoreErr)
	aggFactory := createAggFactoryMock(agg)
	return dispatcher.NewDispatcher(eventStore, aggFactory)
}

func createDispatcherWithStaticEventApplier(
	aggID testdata.StringIdentifier, loadedEvents []domain.DomainEvent, eventStoreErr error) *dispatcher.Dispatcher {
	agg := aggregate.NewBase(testdata.NewTestAggregate(aggID), nil, aggregate.NewStaticEventApplier())
	eventStore := createEventStoreMock(loadedEvents, eventStoreErr, nil)
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

func createEventStoreMock(want []domain.DomainEvent, loadErr error, storeErr error) *esMock.EventStoreMock {
	eventStore := &esMock.EventStoreMock{
		Loader: func(aggregateID domain.Identifier) ([]domain.DomainEvent, error) {
			return want, loadErr
		},
		Saver: func(aggregateID domain.Identifier, version int, events []domain.DomainEvent) error {
			return storeErr
		},
	}
	return eventStore
}
