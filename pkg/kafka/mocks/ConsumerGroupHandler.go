// Code generated by mockery v2.40.3. DO NOT EDIT.

package mocks

import (
	sarama "github.com/IBM/sarama"
	mock "github.com/stretchr/testify/mock"
)

// ConsumerGroupHandler is an autogenerated mock type for the ConsumerGroupHandler type
type ConsumerGroupHandler struct {
	mock.Mock
}

type ConsumerGroupHandler_Expecter struct {
	mock *mock.Mock
}

func (_m *ConsumerGroupHandler) EXPECT() *ConsumerGroupHandler_Expecter {
	return &ConsumerGroupHandler_Expecter{mock: &_m.Mock}
}

// Cleanup provides a mock function with given fields: _a0
func (_m *ConsumerGroupHandler) Cleanup(_a0 sarama.ConsumerGroupSession) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Cleanup")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(sarama.ConsumerGroupSession) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ConsumerGroupHandler_Cleanup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Cleanup'
type ConsumerGroupHandler_Cleanup_Call struct {
	*mock.Call
}

// Cleanup is a helper method to define mock.On call
//   - _a0 sarama.ConsumerGroupSession
func (_e *ConsumerGroupHandler_Expecter) Cleanup(_a0 interface{}) *ConsumerGroupHandler_Cleanup_Call {
	return &ConsumerGroupHandler_Cleanup_Call{Call: _e.mock.On("Cleanup", _a0)}
}

func (_c *ConsumerGroupHandler_Cleanup_Call) Run(run func(_a0 sarama.ConsumerGroupSession)) *ConsumerGroupHandler_Cleanup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(sarama.ConsumerGroupSession))
	})
	return _c
}

func (_c *ConsumerGroupHandler_Cleanup_Call) Return(_a0 error) *ConsumerGroupHandler_Cleanup_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ConsumerGroupHandler_Cleanup_Call) RunAndReturn(run func(sarama.ConsumerGroupSession) error) *ConsumerGroupHandler_Cleanup_Call {
	_c.Call.Return(run)
	return _c
}

// ConsumeClaim provides a mock function with given fields: _a0, _a1
func (_m *ConsumerGroupHandler) ConsumeClaim(_a0 sarama.ConsumerGroupSession, _a1 sarama.ConsumerGroupClaim) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for ConsumeClaim")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ConsumerGroupHandler_ConsumeClaim_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConsumeClaim'
type ConsumerGroupHandler_ConsumeClaim_Call struct {
	*mock.Call
}

// ConsumeClaim is a helper method to define mock.On call
//   - _a0 sarama.ConsumerGroupSession
//   - _a1 sarama.ConsumerGroupClaim
func (_e *ConsumerGroupHandler_Expecter) ConsumeClaim(_a0 interface{}, _a1 interface{}) *ConsumerGroupHandler_ConsumeClaim_Call {
	return &ConsumerGroupHandler_ConsumeClaim_Call{Call: _e.mock.On("ConsumeClaim", _a0, _a1)}
}

func (_c *ConsumerGroupHandler_ConsumeClaim_Call) Run(run func(_a0 sarama.ConsumerGroupSession, _a1 sarama.ConsumerGroupClaim)) *ConsumerGroupHandler_ConsumeClaim_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(sarama.ConsumerGroupSession), args[1].(sarama.ConsumerGroupClaim))
	})
	return _c
}

func (_c *ConsumerGroupHandler_ConsumeClaim_Call) Return(_a0 error) *ConsumerGroupHandler_ConsumeClaim_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ConsumerGroupHandler_ConsumeClaim_Call) RunAndReturn(run func(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error) *ConsumerGroupHandler_ConsumeClaim_Call {
	_c.Call.Return(run)
	return _c
}

// Setup provides a mock function with given fields: _a0
func (_m *ConsumerGroupHandler) Setup(_a0 sarama.ConsumerGroupSession) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Setup")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(sarama.ConsumerGroupSession) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ConsumerGroupHandler_Setup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Setup'
type ConsumerGroupHandler_Setup_Call struct {
	*mock.Call
}

// Setup is a helper method to define mock.On call
//   - _a0 sarama.ConsumerGroupSession
func (_e *ConsumerGroupHandler_Expecter) Setup(_a0 interface{}) *ConsumerGroupHandler_Setup_Call {
	return &ConsumerGroupHandler_Setup_Call{Call: _e.mock.On("Setup", _a0)}
}

func (_c *ConsumerGroupHandler_Setup_Call) Run(run func(_a0 sarama.ConsumerGroupSession)) *ConsumerGroupHandler_Setup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(sarama.ConsumerGroupSession))
	})
	return _c
}

func (_c *ConsumerGroupHandler_Setup_Call) Return(_a0 error) *ConsumerGroupHandler_Setup_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ConsumerGroupHandler_Setup_Call) RunAndReturn(run func(sarama.ConsumerGroupSession) error) *ConsumerGroupHandler_Setup_Call {
	_c.Call.Return(run)
	return _c
}

// NewConsumerGroupHandler creates a new instance of ConsumerGroupHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConsumerGroupHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *ConsumerGroupHandler {
	mock := &ConsumerGroupHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}