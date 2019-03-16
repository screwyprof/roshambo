package aggregate_test

import (
	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testdata"
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// ensure that factory implements domain.AggregateFactory interface.
var _ domain.AggregateFactory = (*aggregate.Factory)(nil)

func TestNewFactory(t *testing.T) {
	t.Run("ItReturnsNewFactoryInstance", func(t *testing.T) {
		f := aggregate.NewFactory()
		assert.True(t, f != nil)
	})
}

func TestFactoryCreateAggregate(t *testing.T) {
	t.Run("ItPanicsIfTheAggregateIsNotRegistered", func(t *testing.T) {
		f := aggregate.NewFactory()

		factory := func() {
			f.CreateAggregate("testdata.TestAggregate", testdata.StringIdentifier("TestAgg"))
		}

		assert.Panic(t, factory)
	})
}

func TestFactoryRegisterAggregate(t *testing.T) {
	t.Run("ItPanicsIfTheAggregateIsNotRegistered", func(t *testing.T) {
		ID := testdata.StringIdentifier("TestAgg")
		expected := aggregate.NewBase(testdata.NewTestAggregate(ID), nil, nil)

		f := aggregate.NewFactory()
		f.RegisterAggregate(func(ID domain.Identifier) domain.AdvancedAggregate {
			return expected
		})

		newAgg := f.CreateAggregate("testdata.TestAggregate", ID)

		assert.Equals(t, expected, newAgg)
	})
}
