package eventstore_test

import (
	"fmt"

	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/eventstore"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/testdata/mock"

	"github.com/screwyprof/roshambo/pkg/domain"
)

func ExampleInMemoryEventStore_LoadEventsFor() {
	ID := mock.StringIdentifier("TestAgg")
	aggregate.NewBase(mock.NewTestAggregate(ID), nil, nil)

	es := eventstore.NewInInMemoryEventStore()
	_ = es.StoreEventsFor(ID, 0, []domain.DomainEvent{mock.SomethingHappened{}})

	events, _ := es.LoadEventsFor(ID)
	fmt.Printf("%#v", events)

	// Output:
	// []domain.DomainEvent{mock.SomethingHappened{}}
}

func ExampleInMemoryEventStoreStoreEventsFor_ConcurrencyError() {
	ID := mock.StringIdentifier("TestAgg")
	aggregate.NewBase(mock.NewTestAggregate(ID), nil, nil)

	es := eventstore.NewInInMemoryEventStore()
	err := es.StoreEventsFor(ID, 1, []domain.DomainEvent{mock.SomethingHappened{}})

	fmt.Printf("%v", err)

	// Output:
	// concurrency error: aggregate versions differ
}
