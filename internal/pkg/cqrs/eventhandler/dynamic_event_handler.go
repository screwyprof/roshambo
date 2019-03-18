package eventhandler

import (
	"reflect"
	"strings"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// Dynamic handles events.
//
// It automatically registers all event handlers for the given entity.
type Dynamic struct {
	*Static
}

// NewDynamic creates a new instance of Dynamic
func NewDynamic() *Dynamic {
	return &Dynamic{
		Static: NewStatic(),
	}
}

// RegisterHandlers registers all the event handlers found in .
func (h *Dynamic) RegisterHandlers(entity interface{}) {
	entityType := reflect.TypeOf(entity)
	for i := 0; i < entityType.NumMethod(); i++ {
		method := entityType.Method(i)
		if !strings.HasPrefix(method.Name, "On") {
			continue
		}

		h.Static.RegisterHandler(method.Name, func(e domain.DomainEvent) error {
			result := method.Func.Call([]reflect.Value{reflect.ValueOf(entity), reflect.ValueOf(e)})
			resErr := result[0].Interface()
			if resErr != nil {
				return resErr.(error)
			}
			return nil
		})
	}
}
