package bus

import (
	"distributed-chat/src/shared/domain/bus"
	"errors"
	"reflect"
)

type memoryQueryBus struct {
	handlers bus.Handlers
}

func NewMemoryQueryBus() bus.QueryBus {
	return &memoryQueryBus{
		handlers: make(bus.Handlers),
	}
}

func (m *memoryQueryBus) Register(command interface{}, handler interface{}) error {
	key := reflect.TypeOf(command).String()
	if _, exists := m.handlers[key]; exists {
		return errors.New("handler already registered for this query")
	}
	m.handlers[key] = handler
	return nil
}

func (m *memoryQueryBus) Handlers() bus.Handlers {
	return m.handlers
}

func (m *memoryQueryBus) Execute(query interface{}) (interface{}, error) {
	key := reflect.TypeOf(query).String()
	handler, exists := m.handlers[key]
	if !exists {
		return nil, errors.New("no handler registered for this query")
	}

	if fn, ok := handler.(func(interface{}) (interface{}, error)); ok {
		return fn(query)
	}

	return nil, errors.New("handler has an invalid signature")
}
