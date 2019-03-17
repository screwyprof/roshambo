package fixture

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// GivenFn is a test init function.
type GivenFn func() domain.CommandHandler

// WhenFn is a command handler function.
type WhenFn func(dispatcher domain.CommandHandler) ([]domain.DomainEvent, error)

// ThenFn prepares the Checker.
type ThenFn func(t *testing.T) Checker

// Checker asserts the given results.
type Checker func(got []domain.DomainEvent, err error)

// DispatcherTester defines a dispatcher tester.
type DispatcherTester func(given GivenFn, when WhenFn, then ThenFn)

// Test runs the test.
//
// Example:
//  Test(t)(
//	  Given(dispatcher),
//	  When(testdata.TestCommand{Param: "param"}),
//	  Then(testdata.TestEvent{Data: "param"}),
//  )
func Test(t *testing.T) DispatcherTester {
	return func(given GivenFn, when WhenFn, then ThenFn) {
		t.Helper()
		then(t)(when(given()))
	}
}

// Given prepares the given aggregate for testing.
func Given(dispatcher domain.CommandHandler) GivenFn {
	return func() domain.CommandHandler {
		return dispatcher
	}
}

// When prepares the command handler for the given command.
func When(cmd ...domain.Command) WhenFn {
	return func(dispatcher domain.CommandHandler) ([]domain.DomainEvent, error) {
		var events []domain.DomainEvent
		for _, c := range cmd {
			e, err := dispatcher.Handle(c)
			if err != nil {
				return nil, err
			}
			events = append(events, e...)
		}
		return events, nil
	}
}

// Then asserts that the expected events are applied.
func Then(want ...domain.DomainEvent) ThenFn {
	return func(t *testing.T) Checker {
		return func(got []domain.DomainEvent, err error) {
			t.Helper()
			assert.Ok(t, err)
			assert.Equals(t, want, got)
		}
	}
}

// ThenFailWith asserts that the expected error occurred.
func ThenFailWith(want error) ThenFn {
	return func(t *testing.T) Checker {
		return func(got []domain.DomainEvent, err error) {
			t.Helper()
			assert.Equals(t, want, err)
		}
	}
}
