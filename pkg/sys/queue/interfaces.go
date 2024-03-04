package queue

import (
	"context"
)

type HandlerFunc func(ctx context.Context, msg Message) error

type Middleware func(h HandlerFunc) HandlerFunc

type router interface {
	Routes() map[MessageType]HandlerFunc
}
