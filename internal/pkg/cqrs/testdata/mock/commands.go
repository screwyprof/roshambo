package mock

import "github.com/screwyprof/roshambo/pkg/domain"

type MakeSomethingHappen struct{
	AggID domain.Identifier
}

func (c MakeSomethingHappen) AggregateID() domain.Identifier {
	return c.AggID
}

func (c MakeSomethingHappen) AggregateType() string {
	return "mock.TestAggregate"
}
func (c MakeSomethingHappen) CommandType() string {
	return "MakeSomethingHappen"
}