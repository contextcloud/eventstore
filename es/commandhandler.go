package es

import (
	"context"
	"fmt"
)

type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) error
}

// CommandHandlerFunc is a function that can be used as a command handler.
type CommandHandlerFunc func(context.Context, Command) error

// Handle is a method of the CommandHandler.
func (h CommandHandlerFunc) Handle(ctx context.Context, cmd Command) error {
	return h(ctx, cmd)
}

type IsCommandHandler interface {
	IsCommandHandler()
}

type BaseCommmandHandler struct {
}

func (BaseCommmandHandler) IsCommandHandler() {}

type commandHandler struct {
	h       IsCommandHandler
	handles CommandHandles
}

func (h *commandHandler) Handle(ctx context.Context, cmd Command) error {
	if h.handles == nil {
		return fmt.Errorf("no handler for command: %T", cmd)
	}

	return h.handles.Handle(h.h, ctx, cmd)
}

func NewCommandHandler(h IsCommandHandler, handles CommandHandles) CommandHandler {
	return &commandHandler{
		h:       h,
		handles: handles,
	}
}
