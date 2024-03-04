package logadapter

import (
	"github.com/itimky/spindle/pkg/sys/log"
)

type Nop struct{}

func NewNop() Nop {
	return Nop{}
}

func (l Nop) Info(_ string) {}

func (l Nop) Infof(_ string, _ ...interface{}) {}

func (l Nop) Error(_ string) {}

func (l Nop) Errorf(_ string, _ ...interface{}) {}

func (l Nop) WithField(_ string, _ string) log.Logger { //nolint: ireturn
	return l
}
