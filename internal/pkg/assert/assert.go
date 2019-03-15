package assert

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert/deep"
)

// True fails the test if the condition is false.
func True(tb testing.TB, condition bool) {
	tb.Helper()
	if !condition {
		tb.Errorf("\033[31mcondition is false\033[39m\n\n")
	}
}

// Ok fails the test if an err is not nil.
func Ok(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Errorf("\033[31munexpected error: %v\033[39m\n\n", err)
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	tb.Helper()
	if diff := deep.Equal(exp, act); diff != nil {
		tb.Error("\033[31m", diff, "\033[39m")
	}
}

// Panic fails the test if it didn't panic.
func Panic(tb testing.TB, f func()) {
	tb.Helper()
	defer func() {
		tb.Helper()
		if r := recover(); r == nil {
			tb.Errorf("\033[31mpanic is expected\033[39m")
		}
	}()
	f()
}
