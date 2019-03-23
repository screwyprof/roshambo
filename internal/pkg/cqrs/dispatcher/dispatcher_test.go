package dispatcher_test

import (
	"testing"

	"github.com/segmentio/ksuid"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/dispatcher"
	. "github.com/screwyprof/roshambo/internal/pkg/cqrs/dispatcher/testdata/fixture"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/testdata/mock"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// ensure that Dispatcher  implements domain.CommandHandler interface.
var _ domain.CommandHandler = (*dispatcher.Dispatcher)(nil)

func TestNewDispatcher(t *testing.T) {
	t.Run("ItPanicsIfAggregateStoreIsNotGiven", func(t *testing.T) {
		factory := func() {
			dispatcher.NewDispatcher(nil, nil)
		}
		assert.Panic(t, factory)
	})

	t.Run("ItPanicsIfEventPublisherIsNotGiven", func(t *testing.T) {
		factory := func() {
			dispatcher.NewDispatcher(
				createAggregateStoreMock(nil, nil, nil),
				nil,
			)
		}
		assert.Panic(t, factory)
	})
}

func TestNewDispatcherHandle(t *testing.T) {
	t.Run("ItFailsIfItCannotLoadAggregate", func(t *testing.T) {
		ID := ksuid.New()
		Test(t)(
			Given(createDispatcher(
				ID,
				withAggregateStoreLoadErr(mock.ErrAggregateStoreCannotLoadAggregate),
			)),
			When(mock.MakeSomethingHappen{AggID: ID}),
			ThenFailWith(mock.ErrAggregateStoreCannotLoadAggregate),
		)
	})

	t.Run("ItFailsIfAggregateCannotHandleTheGivenCommand", func(t *testing.T) {
		ID := ksuid.New()
		Test(t)(
			Given(createDispatcher(
				ID,
				withLoadedEvents([]domain.DomainEvent{mock.SomethingHappened{}}),
			)),
			When(mock.MakeSomethingHappen{AggID: ID}),
			ThenFailWith(mock.ErrItCanHappenOnceOnly),
		)
	})

	t.Run("ItFailsIfItCannotStoreAggregate", func(t *testing.T) {
		ID := ksuid.New()
		Test(t)(
			Given(createDispatcher(
				ID,
				withAggregateStoreSaveErr(mock.ErrAggregateStoreCannotStoreAggregate),
			)),
			When(mock.MakeSomethingHappen{AggID: ID}),
			ThenFailWith(mock.ErrAggregateStoreCannotStoreAggregate),
		)
	})

	t.Run("ItFailsIfItCannotPublishEvents", func(t *testing.T) {
		ID := ksuid.New()
		Test(t)(
			Given(createDispatcher(ID, withPublisherErr(mock.ErrCannotPublishEvents))),
			When(mock.MakeSomethingHappen{AggID: ID}),
			ThenFailWith(mock.ErrCannotPublishEvents),
		)
	})

	t.Run("ItReturnsEvents", func(t *testing.T) {
		ID := ksuid.New()
		Test(t)(
			Given(createDispatcher(ID)),
			When(mock.MakeSomethingHappen{AggID: ID}),
			Then(mock.SomethingHappened{}),
		)
	})
}

type dispatcherOptions struct {
	emptyFactory       bool
	staticEventApplier bool
	loadedEvents       []domain.DomainEvent

	loadErr      error
	storeErr     error
	publisherErr error
}

type option func(*dispatcherOptions)

func withLoadedEvents(loadedEvents []domain.DomainEvent) option {
	return func(o *dispatcherOptions) {
		o.loadedEvents = loadedEvents
	}
}

func withAggregateStoreLoadErr(err error) option {
	return func(o *dispatcherOptions) {
		o.loadErr = err
	}
}

func withAggregateStoreSaveErr(err error) option {
	return func(o *dispatcherOptions) {
		o.storeErr = err
	}
}

func withPublisherErr(err error) option {
	return func(o *dispatcherOptions) {
		o.publisherErr = err
	}
}

func createDispatcher(ID domain.Identifier, opts ...option) *dispatcher.Dispatcher {
	config := &dispatcherOptions{}
	for _, opt := range opts {
		opt(config)
	}

	pureAgg := mock.NewTestAggregate(ID)

	commandHandler := aggregate.NewCommandHandler()
	commandHandler.RegisterHandlers(pureAgg)

	eventApplier := aggregate.NewDynamicEventApplier()
	eventApplier.RegisterAppliers(pureAgg)

	agg := aggregate.NewBase(pureAgg, commandHandler, eventApplier)
	if config.loadedEvents != nil {
		_ = agg.Apply(config.loadedEvents...)
	}

	return dispatcher.NewDispatcher(
		createAggregateStoreMock(agg, config.loadErr, config.storeErr),
		createEventPublisherMock(config.publisherErr),
	)
}

func createAggregateStoreMock(want domain.AdvancedAggregate, loadErr error, storeErr error) *mock.AggregateStoreMock {
	eventStore := &mock.AggregateStoreMock{
		Loader: func(aggregateID domain.Identifier, aggregateType string) (domain.AdvancedAggregate, error) {
			return want, loadErr
		},
		Saver: func(aggregate domain.AdvancedAggregate, events ...domain.DomainEvent) error {
			return storeErr
		},
	}
	return eventStore
}

func createEventPublisherMock(err error) *mock.EventPublisherMock {
	eventPublisher := &mock.EventPublisherMock{
		Publisher: func(e ...domain.DomainEvent) error {
			return err
		},
	}
	return eventPublisher
}
