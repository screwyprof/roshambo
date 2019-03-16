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
	t.Run("ItUsesCustomCommandHandlerAndEventApplierWhenProvided", func(t *testing.T) {
		Test(t)(
			Given(createTestAggregateWithCustomCommandHandlerAndEventApplier()),
			When(testdata.MakeSomethingHappen{}),
			Then(testdata.SomethingHappened{}),
		)
	})

	t.Run("ItReturnsAnErrorIfTheHandlerIsNotFound", func(t *testing.T) {
		Test(t)(
			Given(createTestAggWithEmptyCommandHandler()),
			When(testdata.MakeSomethingHappen{}),
			ThenFailWith(testdata.ErrMakeSomethingHandlerNotFound),
		)
	})

	t.Run("ItReturnsAnErrorIfTheEventAppliersNotFound", func(t *testing.T) {
		Test(t)(
			Given(createTestAggWithEmptyEventApplier()),
			When(testdata.MakeSomethingHappen{}),
			ThenFailWith(testdata.ErrOnSomethingHappenedApplierNotFound),
		)
	})
}

func TestBaseVersion(t *testing.T) {
	t.Run("ItReturnsVersion", func(t *testing.T) {
		agg := createTestAggWithDefaultCommandHandlerAndEventApplier()

		assert.Equals(t, 0, agg.Version())
	})
}

func TestBaseApply(t *testing.T) {
	t.Run("ItAppliesEventsAndReturnsSomeBusinessError", func(t *testing.T) {
		Test(t)(
			Given(createTestAggWithDefaultCommandHandlerAndEventApplier(), testdata.SomethingHappened{}),
			When(testdata.MakeSomethingHappen{}),
			ThenFailWith(testdata.ErrItCanHappenOnceOnly),
		)
	})

	t.Run("ItReturnsAnErrorIfTheEventAppliersNotFound", func(t *testing.T) {
		Test(t)(
			Given(createTestAggWithEmptyEventApplier(), testdata.SomethingHappened{}),
			When(testdata.MakeSomethingHappen{}),
			ThenFailWith(testdata.ErrOnSomethingHappenedApplierNotFound),
		)
	})

	t.Run("ItIncrementsVersion", func(t *testing.T) {
		agg := createTestAggWithEmptyCommandHandler()

		err := agg.Apply(testdata.SomethingHappened{})

		assert.Ok(t, err)
		assert.Equals(t, 1, agg.Version())
	})
}

func createTestAggWithDefaultCommandHandlerAndEventApplier() *aggregate.Base {
	ID := testdata.StringIdentifier("TestAgg1")
	pureAgg := testdata.NewTestAggregate(ID)

	return aggregate.NewBase(pureAgg, nil, nil)
}

func createTestAggregateWithCustomCommandHandlerAndEventApplier() *aggregate.Base {
	ID := testdata.StringIdentifier("TestAgg1")
	a := testdata.NewTestAggregate(ID)

	return aggregate.NewBase(a, createCommandHandler(a), createEventApplier(a))
}

func createTestAggWithEmptyCommandHandler() *aggregate.Base {
	ID := testdata.StringIdentifier("TestAgg1")
	pureAgg := testdata.NewTestAggregate(ID)

	return aggregate.NewBase(pureAgg, aggregate.NewStaticCommandHandler(), nil)
}

func createTestAggWithEmptyEventApplier() *aggregate.Base {
	ID := testdata.StringIdentifier("TestAgg1")
	pureAgg := testdata.NewTestAggregate(ID)

	return aggregate.NewBase(pureAgg, nil, aggregate.NewStaticEventApplier())
}

func createEventApplier(pureAgg *testdata.TestAggregate) *aggregate.StaticEventApplier {
	eventApplier := aggregate.NewStaticEventApplier()
	eventApplier.RegisterApplier("OnSomethingHappened", func(e domain.DomainEvent) {
		pureAgg.OnSomethingHappened(e.(testdata.SomethingHappened))
	},
	)
	return eventApplier
}

func createCommandHandler(pureAgg *testdata.TestAggregate) *aggregate.StaticCommandHandler {
	commandHandler := aggregate.NewStaticCommandHandler()
	commandHandler.RegisterHandler("MakeSomethingHappen", func(c domain.Command) ([]domain.DomainEvent, error) {
		return pureAgg.MakeSomethingHappen(c.(testdata.MakeSomethingHappen))
	},
	)
	return commandHandler
}
