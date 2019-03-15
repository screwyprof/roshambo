package domain

import "fmt"

// Identifier an object identifier.
type Identifier interface {
	fmt.Stringer
}

// Command is an object that is sent to the domain to change state.
//
// People request changes to the domain by sending commands.
// Command are named with a verb in the imperative mood, for example ConfirmOrder.
type Command interface {
	CommandType() string
}

// CommandHandler executes commands.
//
// A command handler receives a command and brokers a result from the appropriate aggregate.
// "A result" is either a successful application of the command, or an error.
//
// Could be implemented like a method on the aggregate.
type CommandHandler interface {
	Handle(c Command) ([]DomainEvent, error)
}

// CommandHandlerFunc is a function that can be used as a command handler.
type CommandHandlerFunc func(Command) ([]DomainEvent, error)

// DomainEvent represents something that took place in the domain.
//
// Events are always named with a past-participle verb, such as OrderConfirmed.
type DomainEvent interface {
	EventType() string
}

// EventApplier applies the given events to an aggregate.
type EventApplier interface {
	Apply(e ...DomainEvent) error
}

// EventApplierFunc is a function that can be used as an event applier.
type EventApplierFunc func(DomainEvent)

// Aggregate is a cluster of domain objects that can be treated as a single unit.
//
// Every transaction is scoped to a single aggregate.
// The lifetimes of the components of an aggregate are bounded by
// the lifetime of the entire aggregate.
//
// Concretely, an aggregate will handle commands, apply events,
// and have a state model encapsulated within it that allows it to implement the required command validation,
// thus upholding the invariants (business rules) of the aggregate.
type Aggregate interface {
	AggregateID() Identifier
}

// AdvancedAggregate is an aggregate which handles commands
// and applies events after it automatically
type AdvancedAggregate interface {
	Aggregate
	CommandHandler
	EventApplier
}

// EventStore stores and loads events.
type EventStore interface {
	LoadEventsFor(aggregateID Identifier) ([]DomainEvent, error)
	StoreEventsFor(aggregateID Identifier, version int, events []DomainEvent) error
}
