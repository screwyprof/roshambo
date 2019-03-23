package store_test

import (
	"testing"

	"github.com/segmentio/ksuid"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/testdata/mock"

	"github.com/screwyprof/roshambo/internal/pkg/cqrs/store"
	"github.com/screwyprof/roshambo/pkg/domain"
)

// ensure that AggregateStore implements domain.AggregateStore interface.
var _ domain.AggregateStore = (*store.AggregateStore)(nil)

func TestNewStore(t *testing.T) {
	t.Run("ItPanicsIfEventStoreIsNotGiven", func(t *testing.T) {
		factory := func() {
			store.NewStore(nil, nil)
		}
		assert.Panic(t, factory)
	})

	t.Run("ItPanicsIfAggregateFactoryIsNotGiven", func(t *testing.T) {
		factory := func() {
			store.NewStore(
				createEventStoreMock(nil, nil, nil),
				nil,
			)
		}
		assert.Panic(t, factory)
	})
}

func TestAggregateStoreLoad(t *testing.T) {
	t.Run("ItFailsIfItCannotLoadEventsForAggregate", func(t *testing.T) {
		// arrange
		ID := ksuid.New()
		s := createAggregateStore(ID, withEventStoreLoadErr(mock.ErrEventStoreCannotLoadEvents))

		// act
		_, err := s.Load(ID, mock.TestAggregateType)

		// assert
		assert.Equals(t, mock.ErrEventStoreCannotLoadEvents, err)
	})

	t.Run("ItCannotCreateAggregate", func(t *testing.T) {
		// arrange
		ID := ksuid.New()
		s := createAggregateStore(ID, withEmptyFactory())

		// act
		_, err := s.Load(ID, mock.TestAggregateType)

		// assert
		assert.Equals(t, mock.ErrAggIsNotRegistered, err)
	})

	t.Run("ItFailsIfItCannotApplyEvents", func(t *testing.T) {
		// arrange
		ID := ksuid.New()
		s := createAggregateStore(
			ID,
			withLoadedEvents([]domain.DomainEvent{mock.SomethingHappened{}}),
			withStaticEventApplier(),
		)

		// act
		_, err := s.Load(ID, mock.TestAggregateType)

		// assert
		assert.Equals(t, mock.ErrOnSomethingHappenedApplierNotFound, err)
	})

	t.Run("ItReturnsAggregate", func(t *testing.T) {
		// arrange
		ID := ksuid.New()
		s := createAggregateStore(ID)

		// act
		got, err := s.Load(ID, mock.TestAggregateType)

		// assert
		assert.Ok(t, err)
		assert.True(t, nil != got)
	})
}

func TestAggregateStoreStore(t *testing.T) {
	t.Run("ItFailsIfItCannotSafeEventsForAggregate", func(t *testing.T) {
		// arrange
		ID := ksuid.New()
		s := createAggregateStore(ID, withEventStoreSaveErr(mock.ErrEventStoreCannotStoreEvents))
		agg := createAgg(ID)

		// act
		err := s.Store(agg, nil)

		// assert
		assert.Equals(t, mock.ErrEventStoreCannotStoreEvents, err)
	})
}

func createAgg(ID ksuid.KSUID) *aggregate.Base {
	pureAgg := mock.NewTestAggregate(ID)

	commandHandler := aggregate.NewCommandHandler()
	commandHandler.RegisterHandlers(pureAgg)

	eventApplier := aggregate.NewEventApplier()
	eventApplier.RegisterAppliers(pureAgg)

	return aggregate.NewBase(pureAgg, commandHandler, eventApplier)
}

type aggregateStoreOptions struct {
	emptyFactory       bool
	staticEventApplier bool
	loadedEvents       []domain.DomainEvent

	loadErr      error
	storeErr     error
	publisherErr error
}

type option func(*aggregateStoreOptions)

func withStaticEventApplier() option {
	return func(o *aggregateStoreOptions) {
		o.staticEventApplier = true
	}
}

func withEmptyFactory() option {
	return func(o *aggregateStoreOptions) {
		o.emptyFactory = true
	}
}

func withLoadedEvents(loadedEvents []domain.DomainEvent) option {
	return func(o *aggregateStoreOptions) {
		o.loadedEvents = loadedEvents
	}
}

func withEventStoreLoadErr(err error) option {
	return func(o *aggregateStoreOptions) {
		o.loadErr = err
	}
}

func withEventStoreSaveErr(err error) option {
	return func(o *aggregateStoreOptions) {
		o.storeErr = err
	}
}

func createAggregateStore(ID domain.Identifier, opts ...option) *store.AggregateStore {
	config := &aggregateStoreOptions{}
	for _, opt := range opts {
		opt(config)
	}

	pureAgg := mock.NewTestAggregate(ID)

	applier := aggregate.NewEventApplier()
	if !config.staticEventApplier {
		applier.RegisterAppliers(pureAgg)
	}

	commandHandler := aggregate.NewCommandHandler()
	commandHandler.RegisterHandlers(pureAgg)

	agg := aggregate.NewBase(pureAgg, commandHandler, applier)
	if config.loadedEvents != nil {
		_ = agg.Apply(config.loadedEvents...)
	}
	aggFactory := createAggFactory(agg, config.emptyFactory)
	eventStore := createEventStoreMock(config.loadedEvents, config.loadErr, config.storeErr)

	return store.NewStore(eventStore, aggFactory)
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

func createEventStoreMock(want []domain.DomainEvent, loadErr error, storeErr error) *mock.EventStoreMock {
	eventStore := &mock.EventStoreMock{
		Loader: func(aggregateID domain.Identifier) ([]domain.DomainEvent, error) {
			return want, loadErr
		},
		Saver: func(aggregateID domain.Identifier, version int, events []domain.DomainEvent) error {
			return storeErr
		},
	}
	return eventStore
}
