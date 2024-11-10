package bus

import (
	"distributed-chat/src/shared/domain/bus"
	"errors"
	"reflect"
)

type memoryCommandBus struct {
	handlers bus.Handlers
}

func NewMemoryCommandBus() bus.CommandBus {
	return &memoryCommandBus{
		handlers: make(bus.Handlers),
	}
}

func (m *memoryCommandBus) Register(command interface{}, handler interface{}) error {
	key := reflect.TypeOf(command).String()
	if _, exists := m.handlers[key]; exists {
		return errors.New("handler already registered for this command")
	}
	m.handlers[key] = handler
	return nil
}

func (m *memoryCommandBus) Handlers() bus.Handlers {
	return m.handlers
}

func (m *memoryCommandBus) Execute(command interface{}) error {
	key := reflect.TypeOf(command).String()
	handler, exists := m.handlers[key]
	if !exists {
		return errors.New("no handler registered for this command")
	}

	if fn, ok := handler.(func(interface{}) error); ok {
		return fn(command)
	}

	return errors.New("handler has an invalid signature")
}
