package bus

type Handlers map[string]interface{}

type CommandBus interface {
	Register(interface{}, interface{}) error

	Handlers() Handlers

	Execute(interface{}) error
}
