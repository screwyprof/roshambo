package eventstore_test

import (
	"fmt"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testdata"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/eventstore"

	"github.com/screwyprof/roshambo/pkg/domain"
)

func ExampleInMemoryEventStore_LoadEventsFor() {
	ID := testdata.StringIdentifier("TestAgg")
	aggregate.NewBase(testdata.NewTestAggregate(ID), nil, nil)

	es := eventstore.NewInInMemoryEventStore()
	_ = es.StoreEventsFor(ID, 0, []domain.DomainEvent{testdata.SomethingHappened{}})

	events, _ := es.LoadEventsFor(ID)
	fmt.Printf("%#v", events)

	// Output:
	// []domain.DomainEvent{testdata.SomethingHappened{}}
}

func ExampleInMemoryEventStoreStoreEventsFor_ConcurrencyError() {
	ID := testdata.StringIdentifier("TestAgg")
	aggregate.NewBase(testdata.NewTestAggregate(ID), nil, nil)

	es := eventstore.NewInInMemoryEventStore()
	err := es.StoreEventsFor(ID, 1, []domain.DomainEvent{testdata.SomethingHappened{}})

	fmt.Printf("%v", err)

	// Output:
	// concurrency error: aggregate versions differ
}