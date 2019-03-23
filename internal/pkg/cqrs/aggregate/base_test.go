package aggregate_test

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"

	. "github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testdata/fixture"
	. "github.com/screwyprof/roshambo/internal/pkg/cqrs/testdata/mock"

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

	t.Run("ItPanicsIfCommandHandlerIsNotGiven", func(t *testing.T) {
		factory := func() {
			aggregate.NewBase(NewTestAggregate(StringIdentifier("Test")), nil, nil)
		}
		assert.Panic(t, factory)
	})

	t.Run("ItPanicsIfEventApplierIsNotGiven", func(t *testing.T) {
		factory := func() {
			aggregate.NewBase(
				NewTestAggregate(StringIdentifier("Test")),
				aggregate.NewDynamicCommandHandler(),
				nil,
			)
		}
		assert.Panic(t, factory)
	})
}

func TestBaseHandle(t *testing.T) {
	t.Run("ItUsesCustomCommandHandlerAndEventApplierWhenProvided", func(t *testing.T) {
		Test(t)(
			Given(createTestAggregateWithCustomCommandHandlerAndEventApplier()),
			When(MakeSomethingHappen{}),
			Then(SomethingHappened{}),
		)
	})

	t.Run("ItReturnsAnErrorIfTheHandlerIsNotFound", func(t *testing.T) {
		Test(t)(
			Given(createTestAggWithEmptyCommandHandler()),
			When(MakeSomethingHappen{}),
			ThenFailWith(ErrMakeSomethingHandlerNotFound),
		)
	})

	t.Run("ItReturnsAnErrorIfTheEventAppliersNotFound", func(t *testing.T) {
		Test(t)(
			Given(createTestAggWithEmptyEventApplier()),
			When(MakeSomethingHappen{}),
			ThenFailWith(ErrOnSomethingHappenedApplierNotFound),
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
			Given(createTestAggWithDefaultCommandHandlerAndEventApplier(), SomethingHappened{}),
			When(MakeSomethingHappen{}),
			ThenFailWith(ErrItCanHappenOnceOnly),
		)
	})

	t.Run("ItReturnsAnErrorIfTheEventAppliersNotFound", func(t *testing.T) {
		Test(t)(
			Given(createTestAggWithEmptyEventApplier(), SomethingHappened{}),
			When(MakeSomethingHappen{}),
			ThenFailWith(ErrOnSomethingHappenedApplierNotFound),
		)
	})

	t.Run("ItIncrementsVersion", func(t *testing.T) {
		agg := createTestAggWithEmptyCommandHandler()

		err := agg.Apply(SomethingHappened{})

		assert.Ok(t, err)
		assert.Equals(t, 1, agg.Version())
	})
}

func createTestAggWithDefaultCommandHandlerAndEventApplier() *aggregate.Base {
	ID := StringIdentifier("TestAgg1")
	pureAgg := NewTestAggregate(ID)

	handler := aggregate.NewDynamicCommandHandler()
	handler.RegisterHandlers(pureAgg)

	applier := aggregate.NewDynamicEventApplier()
	applier.RegisterAppliers(pureAgg)

	return aggregate.NewBase(pureAgg, handler, applier)
}

func createTestAggregateWithCustomCommandHandlerAndEventApplier() *aggregate.Base {
	ID := StringIdentifier("TestAgg1")
	a := NewTestAggregate(ID)

	return aggregate.NewBase(a, createCommandHandler(a), createEventApplier(a))
}

func createTestAggWithEmptyCommandHandler() *aggregate.Base {
	ID := StringIdentifier("TestAgg1")
	pureAgg := NewTestAggregate(ID)

	applier := aggregate.NewDynamicEventApplier()
	applier.RegisterAppliers(pureAgg)

	return aggregate.NewBase(pureAgg, aggregate.NewStaticCommandHandler(), applier)
}

func createTestAggWithEmptyEventApplier() *aggregate.Base {
	ID := StringIdentifier("TestAgg1")
	pureAgg := NewTestAggregate(ID)

	handler := aggregate.NewDynamicCommandHandler()
	handler.RegisterHandlers(pureAgg)

	return aggregate.NewBase(pureAgg, handler, aggregate.NewStaticEventApplier())
}

func createEventApplier(pureAgg *TestAggregate) *aggregate.StaticEventApplier {
	eventApplier := aggregate.NewStaticEventApplier()
	eventApplier.RegisterApplier("OnSomethingHappened", func(e domain.DomainEvent) {
		pureAgg.OnSomethingHappened(e.(SomethingHappened))
	})
	return eventApplier
}

func createCommandHandler(pureAgg *TestAggregate) *aggregate.StaticCommandHandler {
	commandHandler := aggregate.NewStaticCommandHandler()
	commandHandler.RegisterHandler("MakeSomethingHappen", func(c domain.Command) ([]domain.DomainEvent, error) {
		return pureAgg.MakeSomethingHappen(c.(MakeSomethingHappen))
	})
	return commandHandler
}
