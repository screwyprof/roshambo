package dispatcher_test

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testdata"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/dispatcher"
	. "github.com/screwyprof/roshambo/internal/pkg/cqrs/dispatcher/testdata/fixture"
	esMock "github.com/screwyprof/roshambo/internal/pkg/cqrs/eventstore/mock"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/identifier"

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
		ID := identifier.NewUUID()
		Test(t)(
			Given(createDispatcher(
				ID,
				withEventStoreLoadErr(esMock.ErrEventStoreCannotLoadEvents),
			)),
			When(testdata.MakeSomethingHappen{AggID: ID}),
			ThenFailWith(esMock.ErrEventStoreCannotLoadEvents),
		)
	})

	t.Run("ItCannotCreateAggregate", func(t *testing.T) {
		ID := identifier.NewUUID()
		Test(t)(
			Given(createDispatcher(
				ID,
				withEmptyFactory(),
			)),
			When(testdata.MakeSomethingHappen{AggID: ID}),
			ThenFailWith(testdata.ErrAggIsNotRegistered),
		)
	})

	t.Run("ItFailsIfItCannotApplyEvents", func(t *testing.T) {
		ID := identifier.NewUUID()
		Test(t)(
			Given(createDispatcher(
				ID,
				withLoadedEvents([]domain.DomainEvent{testdata.SomethingHappened{}}),
				withStaticEventApplier(),
			)),
			When(testdata.MakeSomethingHappen{AggID: ID}),
			ThenFailWith(testdata.ErrOnSomethingHappenedApplierNotFound),
		)
	})

	t.Run("ItFailsIfAggregateCannotHandleTheGivenCommand", func(t *testing.T) {
		ID := identifier.NewUUID()
		Test(t)(
			Given(createDispatcher(
				ID,
				withLoadedEvents([]domain.DomainEvent{testdata.SomethingHappened{}}),
			)),
			When(testdata.MakeSomethingHappen{AggID: ID}),
			ThenFailWith(testdata.ErrItCanHappenOnceOnly),
		)
	})

	t.Run("ItFailsIfItCannotStoreEvents", func(t *testing.T) {
		ID := identifier.NewUUID()
		Test(t)(
			Given(createDispatcher(
				ID,
				withEventStoreSaveErr(esMock.ErrEventStoreCannotStoreEvents),
			)),
			When(testdata.MakeSomethingHappen{AggID: ID}),
			ThenFailWith(esMock.ErrEventStoreCannotStoreEvents),
		)
	})

	t.Run("ItReturnsEvents", func(t *testing.T) {
		ID := identifier.NewUUID()
		Test(t)(
			Given(createDispatcher(ID)),
			When(testdata.MakeSomethingHappen{AggID: ID}),
			Then(testdata.SomethingHappened{}),
		)
	})
}

type dispatcherOptions struct {
	emptyFactory       bool
	staticEventApplier bool
	loadedEvents       []domain.DomainEvent

	loadErr  error
	storeErr error
}

type option func(*dispatcherOptions)

func withStaticEventApplier() option {
	return func(o *dispatcherOptions) {
		o.staticEventApplier = true
	}
}

func withEmptyFactory() option {
	return func(o *dispatcherOptions) {
		o.emptyFactory = true
	}
}

func withLoadedEvents(loadedEvents []domain.DomainEvent) option {
	return func(o *dispatcherOptions) {
		o.loadedEvents = loadedEvents
	}
}

func withEventStoreLoadErr(err error) option {
	return func(o *dispatcherOptions) {
		o.loadErr = err
	}
}

func withEventStoreSaveErr(err error) option {
	return func(o *dispatcherOptions) {
		o.storeErr = err
	}
}

type eventApplier interface {
	domain.EventApplier
	RegisterAppliers(aggregate domain.Aggregate)
	RegisterApplier(method string, applierFunc domain.EventApplierFunc)
}

func createDispatcher(ID domain.Identifier, opts ...option) *dispatcher.Dispatcher {
	config := &dispatcherOptions{}
	for _, opt := range opts {
		opt(config)
	}

	var applier eventApplier
	if config.staticEventApplier {
		applier = aggregate.NewStaticEventApplier()
	}

	agg := aggregate.NewBase(testdata.NewTestAggregate(ID), nil, applier)
	aggFactory := createAggFactory(agg, config.emptyFactory)
	eventStore := createEventStoreMock(config.loadedEvents, config.loadErr, config.storeErr)

	return dispatcher.NewDispatcher(eventStore, aggFactory)
}

func createAggFactory(agg *aggregate.Base, empty bool) *aggregate.Factory {
	f := aggregate.NewFactory()
	if empty {
		return f
	}
	f.RegisterAggregate(func(ID domain.Identifier) domain.AdvancedAggregate {
		return agg
	})

	return f
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
