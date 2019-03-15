package testfixture

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// GivenFn is a test init function.
type GivenFn func() (domain.AdvancedAggregate, []domain.DomainEvent)

// WhenFn is a command handler function.
type WhenFn func(agg domain.AdvancedAggregate, err error) ([]domain.DomainEvent, error)

// ThenFn prepares the Checker.
type ThenFn func(t *testing.T) Checker

// Checker asserts the given results.
type Checker func(got []domain.DomainEvent, err error)

// AggregateTester defines an aggregate tester.
type AggregateTester func(given GivenFn, when WhenFn, then ThenFn)

// Test runs the test.
//
// Example:
//  Test(t)(
//	  Given(agg),
//	  When(testdata.TestCommand{Param: "param"}),
//	  Then(testdata.TestEvent{Data: "param"}),
//  )
func Test(t *testing.T) AggregateTester {
	return func(given GivenFn, when WhenFn, then ThenFn) {
		then(t)(when(applyEvents(given)))
	}
}

// Given prepares the given aggregate for testing.
func Given(agg domain.AdvancedAggregate, events ...domain.DomainEvent) GivenFn {
	return func() (domain.AdvancedAggregate, []domain.DomainEvent) {
		return agg, events
	}
}

// When prepares the command handler for the given command.
func When(c domain.Command) WhenFn {
	return func(agg domain.AdvancedAggregate, err error) ([]domain.DomainEvent, error) {
		if err != nil {
			return nil, err
		}
		return agg.Handle(c)
	}
}

// Then asserts that the expected events are applied.
func Then(want ...domain.DomainEvent) ThenFn {
	return func(t *testing.T) Checker {
		return func(got []domain.DomainEvent, err error) {
			assert.Ok(t, err)
			assert.Equals(t, want, got)
		}
	}
}

// ThenFailWith asserts that the expected error occurred.
func ThenFailWith(want error) ThenFn {
	return func(t *testing.T) Checker {
		return func(got []domain.DomainEvent, err error) {
			assert.Equals(t, want, err)
		}
	}
}

func applyEvents(given GivenFn) (domain.AdvancedAggregate, error) {
	agg, events := given()
	err := agg.Apply(events...)
	if err != nil {
		return nil, err
	}

	return agg, nil
}
