// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	queue "github.com/itimky/spindle/pkg/sys/queue"
)

// Mockhandler is an autogenerated mock type for the handler type
type Mockhandler struct {
	mock.Mock
}

type Mockhandler_Expecter struct {
	mock *mock.Mock
}

func (_m *Mockhandler) EXPECT() *Mockhandler_Expecter {
	return &Mockhandler_Expecter{mock: &_m.Mock}
}

// Handle provides a mock function with given fields: ctx, msg
func (_m *Mockhandler) Handle(ctx context.Context, msg queue.Message) error {
	ret := _m.Called(ctx, msg)

	if len(ret) == 0 {
		panic("no return value specified for Handle")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, queue.Message) error); ok {
		r0 = rf(ctx, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Mockhandler_Handle_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Handle'
type Mockhandler_Handle_Call struct {
	*mock.Call
}

// Handle is a helper method to define mock.On call
//   - ctx context.Context
//   - msg queue.Message
func (_e *Mockhandler_Expecter) Handle(ctx interface{}, msg interface{}) *Mockhandler_Handle_Call {
	return &Mockhandler_Handle_Call{Call: _e.mock.On("Handle", ctx, msg)}
}

func (_c *Mockhandler_Handle_Call) Run(run func(ctx context.Context, msg queue.Message)) *Mockhandler_Handle_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(queue.Message))
	})
	return _c
}

func (_c *Mockhandler_Handle_Call) Return(_a0 error) *Mockhandler_Handle_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Mockhandler_Handle_Call) RunAndReturn(run func(context.Context, queue.Message) error) *Mockhandler_Handle_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockhandler creates a new instance of Mockhandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockhandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *Mockhandler {
	mock := &Mockhandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}