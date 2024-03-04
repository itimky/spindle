package queue

import (
	"context"
	"fmt"
)

type Handler struct {
	router router

	middlewares []Middleware
}

func NewHandler(
	router router,
	middlewares ...Middleware,
) *Handler {
	return &Handler{
		router:      router,
		middlewares: middlewares,
	}
}

func (r *Handler) Handle(ctx context.Context, msg Message) error {
	handlerFn, ok := r.router.Routes()[msg.Type]
	if !ok {
		return fmt.Errorf("%w: %s", ErrNoRoute, msg.Type)
	}

	// decorate middlewares in reverse order
	// to execute in declaration order
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		middlewareFn := r.middlewares[i]

		handlerFn = middlewareFn(handlerFn)
	}

	return handlerFn(ctx, msg)
}
