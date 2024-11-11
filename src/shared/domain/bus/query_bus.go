package bus

type QueryBus interface {
	Register(interface{}, interface{}) error

	Handlers() Handlers

	Execute(interface{}) (interface{}, error)
}
