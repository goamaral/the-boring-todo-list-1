// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import (
	context "context"

	gorm "gorm.io/gorm"

	gormprovider "example.com/the-boring-to-do-list-1/pkg/gormprovider"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository[T interface{}] struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, entity
func (_m *Repository[T]) Create(ctx context.Context, entity *T) error {
	ret := _m.Called(ctx, entity)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *T) error); ok {
		r0 = rf(ctx, entity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, opts
func (_m *Repository[T]) Delete(ctx context.Context, opts ...gormprovider.QueryOption) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...gormprovider.QueryOption) error); ok {
		r0 = rf(ctx, opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, opts
func (_m *Repository[T]) Get(ctx context.Context, opts ...gormprovider.QueryOption) (T, bool, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 T
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, ...gormprovider.QueryOption) (T, bool, error)); ok {
		return rf(ctx, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...gormprovider.QueryOption) T); ok {
		r0 = rf(ctx, opts...)
	} else {
		r0 = ret.Get(0).(T)
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...gormprovider.QueryOption) bool); ok {
		r1 = rf(ctx, opts...)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(context.Context, ...gormprovider.QueryOption) error); ok {
		r2 = rf(ctx, opts...)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// List provides a mock function with given fields: ctx, opts
func (_m *Repository[T]) List(ctx context.Context, opts ...gormprovider.QueryOption) ([]T, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []T
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...gormprovider.QueryOption) ([]T, error)); ok {
		return rf(ctx, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...gormprovider.QueryOption) []T); ok {
		r0 = rf(ctx, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]T)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...gormprovider.QueryOption) error); ok {
		r1 = rf(ctx, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewQuery provides a mock function with given fields: ctx
func (_m *Repository[T]) NewQuery(ctx context.Context) *gorm.DB {
	ret := _m.Called(ctx)

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func(context.Context) *gorm.DB); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// NewQueryWithOpts provides a mock function with given fields: ctx, opts
func (_m *Repository[T]) NewQueryWithOpts(ctx context.Context, opts ...gormprovider.QueryOption) *gorm.DB {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func(context.Context, ...gormprovider.QueryOption) *gorm.DB); ok {
		r0 = rf(ctx, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// Patch provides a mock function with given fields: ctx, patch, opts
func (_m *Repository[T]) Patch(ctx context.Context, patch *T, opts ...gormprovider.QueryOption) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, patch)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *T, ...gormprovider.QueryOption) error); ok {
		r0 = rf(ctx, patch, opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TableName provides a mock function with given fields:
func (_m *Repository[T]) TableName() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, update, opts
func (_m *Repository[T]) Update(ctx context.Context, update *T, opts ...gormprovider.QueryOption) error {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, update)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *T, ...gormprovider.QueryOption) error); ok {
		r0 = rf(ctx, update, opts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository[T interface{}](t mockConstructorTestingTNewRepository) *Repository[T] {
	mock := &Repository[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
