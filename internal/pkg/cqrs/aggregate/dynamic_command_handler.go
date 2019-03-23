package aggregate

import (
	"reflect"

	"github.com/screwyprof/roshambo/pkg/domain"
)

type DynamicCommandHandler struct {
	*StaticCommandHandler
}

func NewDynamicCommandHandler() *DynamicCommandHandler {
	return &DynamicCommandHandler{
		StaticCommandHandler: NewStaticCommandHandler(),
	}
}

// RegisterHandlers registers all the command handlers found in the aggregate.
func (h *DynamicCommandHandler) RegisterHandlers(aggregate domain.Aggregate) {
	aggregateType := reflect.TypeOf(aggregate)
	for i := 0; i < aggregateType.NumMethod(); i++ {
		method := aggregateType.Method(i)
		if !h.methodHasValidSignature(method) {
			continue
		}

		h.StaticCommandHandler.RegisterHandler(method.Name, func(c domain.Command) ([]domain.DomainEvent, error) {
			return h.invokeCommandHandler(method, aggregate, c)
		})
	}
}

func (h *DynamicCommandHandler) methodHasValidSignature(method reflect.Method) bool {
	if method.Type.NumIn() != 2 {
		return false
	}

	// ensure that the method has a domain.Command as a parameter.
	cmdIntfType := reflect.TypeOf((*domain.Command)(nil)).Elem()
	cmdType := method.Type.In(1)
	if !cmdType.Implements(cmdIntfType) {
		return false
	}

	return true
}

func (h *DynamicCommandHandler) invokeCommandHandler(
	method reflect.Method, aggregate domain.Aggregate, c domain.Command) ([]domain.DomainEvent, error) {
	result := method.Func.Call([]reflect.Value{reflect.ValueOf(aggregate), reflect.ValueOf(c)})
	resErr := result[1].Interface()
	if resErr != nil {
		return nil, resErr.(error)
	}
	eventsIntf := result[0].Interface()
	events := eventsIntf.([]domain.DomainEvent)
	return events, nil
}
