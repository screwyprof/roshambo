package testfixture

import "github.com/screwyprof/roshambo/pkg/domain"

// AggregateFactoryMock mocks aggregate factory.
type AggregateFactoryMock struct {
	Creator   func(aggregateType string, ID domain.Identifier) domain.AdvancedAggregate
	Registrar func(factory domain.FactoryFn)
}

// RegisterAggregate registers an aggregate factory method.
func (m *AggregateFactoryMock) RegisterAggregate(factory domain.FactoryFn) {
	m.Registrar(factory)
}

// CreateAggregate creates an aggregate of a given type.
func (m *AggregateFactoryMock) CreateAggregate(aggregateType string, ID domain.Identifier) domain.AdvancedAggregate {
	return m.Creator(aggregateType, ID)
}

type EventStoreMock struct {
	Loader func(aggregateID domain.Identifier) ([]domain.DomainEvent, error)
	Saver  func(aggregateID domain.Identifier, version int, events []domain.DomainEvent) error
}

func (m *EventStoreMock) LoadEventsFor(aggregateID domain.Identifier) ([]domain.DomainEvent, error) {
	return m.Loader(aggregateID)
}
func (m *EventStoreMock) StoreEventsFor(aggregateID domain.Identifier, version int, events []domain.DomainEvent) error {
	return m.Saver(aggregateID, version, events)
}
