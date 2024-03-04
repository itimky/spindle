package queue

import "errors"

var ErrNoRoute = errors.New("no route")
var ErrPanicRecovered = errors.New("panic recovered")
