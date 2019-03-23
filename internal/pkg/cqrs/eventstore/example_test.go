package eventstore_test

import (
	"fmt"

	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/eventstore"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/testdata/mock"

	"github.com/screwyprof/roshambo/pkg/domain"
)

func ExampleInMemoryEventStoreLoadEventsFor() {
	ID := mock.StringIdentifier("TestAgg")

	es := eventstore.NewInInMemoryEventStore()
	_ = es.StoreEventsFor(ID, 0, []domain.DomainEvent{mock.SomethingHappened{}})

	events, _ := es.LoadEventsFor(ID)
	fmt.Printf("%#v", events)

	// Output:
	// []domain.DomainEvent{mock.SomethingHappened{}}
}

func ExampleInMemoryEventStoreStoreEventsForConcurrencyError() {
	ID := mock.StringIdentifier("TestAgg")

	pureAgg := mock.NewTestAggregate(ID)

	commandHandler := aggregate.NewCommandHandler()
	commandHandler.RegisterHandlers(pureAgg)

	eventApplier := aggregate.NewEventApplier()
	eventApplier.RegisterAppliers(pureAgg)

	aggregate.NewAdvanced(pureAgg, commandHandler, eventApplier)

	es := eventstore.NewInInMemoryEventStore()
	err := es.StoreEventsFor(ID, 1, []domain.DomainEvent{mock.SomethingHappened{}})

	fmt.Printf("%v", err)

	// Output:
	// concurrency error: aggregate versions differ
}
