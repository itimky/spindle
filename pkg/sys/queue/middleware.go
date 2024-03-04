package queue

import (
	"context"
	"fmt"
)

func MiddlewareRecover(
	handlerFunc HandlerFunc,
) HandlerFunc {
	return func(ctx context.Context, msg Message) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%w: %v", ErrPanicRecovered, r)
			}
		}()

		err = handlerFunc(ctx, msg)
		if err != nil {
			return fmt.Errorf("recovery middlware: %w", err)
		}

		return nil
	}
}
