package bus

type CommandBus interface {
	Register(interface{}, interface{}) error

	Handlers() Handlers

	Execute(interface{}) error
}
