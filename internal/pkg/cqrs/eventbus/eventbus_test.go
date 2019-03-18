package eventbus_test

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/eventbus"
	"github.com/screwyprof/roshambo/pkg/domain"
)

// ensure that EventBus implements domain.EventPublisher interface.
var _ domain.EventPublisher = (*eventbus.InMemoryEventBus)(nil)

func TestNewInMemoryEventBus(t *testing.T) {
	t.Run("ItCreatesNewInstance", func(t *testing.T) {
		assert.True(t, eventbus.NewInMemoryEventBus() != nil)
	})
}
