package aggregate

import (
	"errors"
	"sync"

	"github.com/screwyprof/roshambo/internal/pkg/cqrs/identifier"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// Factory handles aggregate creation.
type Factory struct {
	factories   map[string]domain.FactoryFn
	factoriesMu sync.RWMutex
}

func NewFactory() *Factory {
	return &Factory{
		factories: make(map[string]domain.FactoryFn),
	}
}

// RegisterAggregate registers an aggregate factory method.
func (f *Factory) RegisterAggregate(factory domain.FactoryFn) {
	f.factoriesMu.Lock()
	defer f.factoriesMu.Unlock()

	agg := factory(identifier.NewUUID())
	f.factories[agg.AggregateType()] = factory
}

// CreateAggregate creates an aggregate of a given type.
func (f *Factory) CreateAggregate(aggregateType string, ID domain.Identifier) (domain.AdvancedAggregate, error) {
	f.factoriesMu.Lock()
	defer f.factoriesMu.Unlock()

	factory, ok := f.factories[aggregateType]
	if !ok {
		return nil, errors.New(aggregateType + " is not registered")
	}
	return factory(ID), nil
}
