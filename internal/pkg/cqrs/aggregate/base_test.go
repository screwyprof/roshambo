package aggregate_test

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testdata"
	. "github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testfixture"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// ensure that basic aggregate implements domain.AdvancedAggregate interface.
var _ domain.AdvancedAggregate = (*aggregate.Base)(nil)

func TestNewBase(t *testing.T) {
	t.Run("ItPanicsIfThePureAggregateIsNotGiven", func(t *testing.T) {
		factory := func() {
			aggregate.NewBase(nil, nil, nil)
		}
		assert.Panic(t, factory)
	})
}

func TestBaseHandle(t *testing.T) {
	t.Run("ItHandlesTheGivenCommandAndAppliesEventsIfTheHandlerExists", func(t *testing.T) {
		Test(t)(
			Given(createTestAggWithCustomCommandHandlerAndEventApplier()),
			When(testdata.MakeSomethingHappen{}),
			Then(testdata.SomethingHappened{}),
		)
	})

	t.Run("ItReturnsAnErrorIfTheHandlerIsNotFound", func(t *testing.T) {
		Test(t)(
			Given(createTestAggWithDefaultCommandHandlerAndEventApplier()),
			When(testdata.MakeSomethingHappen{}),
			ThenFailWith(testdata.ErrMakeSomethingHandlerNotFound),
		)
	})

	t.Run("ItReturnsAnErrorIfTheEventAppliersNotFound", func(t *testing.T) {
		Test(t)(
			Given(createTestAggWithCustomCommandHandler()),
			When(testdata.MakeSomethingHappen{}),
			ThenFailWith(testdata.ErrOnSomethingHappenedApplierNotFound),
		)
	})
}

func TestBaseApply(t *testing.T) {
	t.Run("ItAppliesEventsAndReturnsSomeBusinessError", func(t *testing.T) {
		Test(t)(
			Given(createTestAggWithCustomCommandHandlerAndEventApplier(), testdata.SomethingHappened{}),
			When(testdata.MakeSomethingHappen{}),
			ThenFailWith(testdata.ErrItCanHappenOnceOnly),
		)
	})

	t.Run("ItReturnsAnErrorIfTheEventAppliersNotFound", func(t *testing.T) {
		Test(t)(
			Given(createTestAggWithCustomCommandHandler(), testdata.SomethingHappened{}),
			When(testdata.MakeSomethingHappen{}),
			ThenFailWith(testdata.ErrOnSomethingHappenedApplierNotFound),
		)
	})
}

func createTestAggWithDefaultCommandHandlerAndEventApplier() *aggregate.Base {
	ID := testdata.StringIdentifier("TestAgg1")
	pureAgg := testdata.NewTestAggregate(ID)

	return aggregate.NewBase(pureAgg, nil, nil)
}

func createTestAggWithCustomCommandHandler() *aggregate.Base {
	ID := testdata.StringIdentifier("TestAgg1")
	pureAgg := testdata.NewTestAggregate(ID)

	return aggregate.NewBase(pureAgg, createCommandHandler(pureAgg), nil)
}

func createTestAggWithCustomCommandHandlerAndEventApplier() *aggregate.Base {
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
