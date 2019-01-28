package cq

// CommandBus command bus
//
//
type CommandBus struct {
	genericRequestBus
}

func (CommandBus) New() RequestBus {
	bus := &CommandBus{}
	bus.genericRequestBus.handlerPrefix = "command."

	return bus
}
