package aggregate

import (
	"reflect"
	"strings"

	"github.com/screwyprof/roshambo/pkg/domain"
)

type dynamicEventApplier struct {
	*StaticEventApplier
}

func newDynamicEventApplier() *dynamicEventApplier {
	return &dynamicEventApplier{
		StaticEventApplier: NewStaticEventApplier(),
	}
}

// RegisterAppliers registers all the event appliers found in the aggregate.
func (a *dynamicEventApplier) RegisterAppliers(aggregate domain.Aggregate) {
	aggregateType := reflect.TypeOf(aggregate)
	for i := 0; i < aggregateType.NumMethod(); i++ {
		method := aggregateType.Method(i)
		if !strings.HasPrefix(method.Name, "On") {
			continue
		}

		a.StaticEventApplier.RegisterApplier(method.Name, func(e domain.DomainEvent) {
			method.Func.Call([]reflect.Value{reflect.ValueOf(aggregate), reflect.ValueOf(e)})
		})
	}
}
