// Code generated by mockery v2.33.3. DO NOT EDIT.

package mock_gorm_provider

import (
	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"
)

// QueryOption is an autogenerated mock type for the QueryOption type
type QueryOption struct {
	mock.Mock
}

type QueryOption_Expecter struct {
	mock *mock.Mock
}

func (_m *QueryOption) EXPECT() *QueryOption_Expecter {
	return &QueryOption_Expecter{mock: &_m.Mock}
}

// Apply provides a mock function with given fields: _a0
func (_m *QueryOption) Apply(_a0 *gorm.DB) *gorm.DB {
	ret := _m.Called(_a0)

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func(*gorm.DB) *gorm.DB); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// QueryOption_Apply_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Apply'
type QueryOption_Apply_Call struct {
	*mock.Call
}

// Apply is a helper method to define mock.On call
//   - _a0 *gorm.DB
func (_e *QueryOption_Expecter) Apply(_a0 interface{}) *QueryOption_Apply_Call {
	return &QueryOption_Apply_Call{Call: _e.mock.On("Apply", _a0)}
}

func (_c *QueryOption_Apply_Call) Run(run func(_a0 *gorm.DB)) *QueryOption_Apply_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*gorm.DB))
	})
	return _c
}

func (_c *QueryOption_Apply_Call) Return(_a0 *gorm.DB) *QueryOption_Apply_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *QueryOption_Apply_Call) RunAndReturn(run func(*gorm.DB) *gorm.DB) *QueryOption_Apply_Call {
	_c.Call.Return(run)
	return _c
}

// NewQueryOption creates a new instance of QueryOption. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewQueryOption(t interface {
	mock.TestingT
	Cleanup(func())
}) *QueryOption {
	mock := &QueryOption{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
