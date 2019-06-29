package handler

import "pas/src/events"

type unsupportedCommand struct {
}

func (c *unsupportedCommand) GetRequestId() string {
	return "command.unsupported_command"
}

type testEventHandler struct {
	dynamicHandle func(event events.Event)
}

func (h *testEventHandler) Handle(event events.Event) {
	h.dynamicHandle(event)
}
