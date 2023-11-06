// Code generated by mockery v2.36.1. DO NOT EDIT.

package mock_gorm_provider

import (
	context "context"

	clause "gorm.io/gorm/clause"

	gorm_provider "example.com/the-boring-to-do-list-1/pkg/gorm_provider"

	mock "github.com/stretchr/testify/mock"
)

// AbstractRepository is an autogenerated mock type for the AbstractRepository type
type AbstractRepository[T gorm_provider.AbstractEntity] struct {
	mock.Mock
}

type AbstractRepository_Expecter[T gorm_provider.AbstractEntity] struct {
	mock *mock.Mock
}

func (_m *AbstractRepository[T]) EXPECT() *AbstractRepository_Expecter[T] {
	return &AbstractRepository_Expecter[T]{mock: &_m.Mock}
}

// Create provides a mock function with given fields: ctx, record, clauses
func (_m *AbstractRepository[T]) Create(ctx context.Context, record *T, clauses ...clause.Expression) error {
	_va := make([]interface{}, len(clauses))
	for _i := range clauses {
		_va[_i] = clauses[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, record)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *T, ...clause.Expression) error); ok {
		r0 = rf(ctx, record, clauses...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AbstractRepository_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type AbstractRepository_Create_Call[T gorm_provider.AbstractEntity] struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx context.Context
//   - record *T
//   - clauses ...clause.Expression
func (_e *AbstractRepository_Expecter[T]) Create(ctx interface{}, record interface{}, clauses ...interface{}) *AbstractRepository_Create_Call[T] {
	return &AbstractRepository_Create_Call[T]{Call: _e.mock.On("Create",
		append([]interface{}{ctx, record}, clauses...)...)}
}

func (_c *AbstractRepository_Create_Call[T]) Run(run func(ctx context.Context, record *T, clauses ...clause.Expression)) *AbstractRepository_Create_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]clause.Expression, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(clause.Expression)
			}
		}
		run(args[0].(context.Context), args[1].(*T), variadicArgs...)
	})
	return _c
}

func (_c *AbstractRepository_Create_Call[T]) Return(_a0 error) *AbstractRepository_Create_Call[T] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AbstractRepository_Create_Call[T]) RunAndReturn(run func(context.Context, *T, ...clause.Expression) error) *AbstractRepository_Create_Call[T] {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function with given fields: ctx, clauses
func (_m *AbstractRepository[T]) Delete(ctx context.Context, clauses ...clause.Expression) error {
	_va := make([]interface{}, len(clauses))
	for _i := range clauses {
		_va[_i] = clauses[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...clause.Expression) error); ok {
		r0 = rf(ctx, clauses...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AbstractRepository_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type AbstractRepository_Delete_Call[T gorm_provider.AbstractEntity] struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx context.Context
//   - clauses ...clause.Expression
func (_e *AbstractRepository_Expecter[T]) Delete(ctx interface{}, clauses ...interface{}) *AbstractRepository_Delete_Call[T] {
	return &AbstractRepository_Delete_Call[T]{Call: _e.mock.On("Delete",
		append([]interface{}{ctx}, clauses...)...)}
}

func (_c *AbstractRepository_Delete_Call[T]) Run(run func(ctx context.Context, clauses ...clause.Expression)) *AbstractRepository_Delete_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]clause.Expression, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(clause.Expression)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *AbstractRepository_Delete_Call[T]) Return(_a0 error) *AbstractRepository_Delete_Call[T] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AbstractRepository_Delete_Call[T]) RunAndReturn(run func(context.Context, ...clause.Expression) error) *AbstractRepository_Delete_Call[T] {
	_c.Call.Return(run)
	return _c
}

// Find provides a mock function with given fields: ctx, clauses
func (_m *AbstractRepository[T]) Find(ctx context.Context, clauses ...clause.Expression) ([]T, error) {
	_va := make([]interface{}, len(clauses))
	for _i := range clauses {
		_va[_i] = clauses[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []T
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...clause.Expression) ([]T, error)); ok {
		return rf(ctx, clauses...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...clause.Expression) []T); ok {
		r0 = rf(ctx, clauses...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]T)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...clause.Expression) error); ok {
		r1 = rf(ctx, clauses...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AbstractRepository_Find_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Find'
type AbstractRepository_Find_Call[T gorm_provider.AbstractEntity] struct {
	*mock.Call
}

// Find is a helper method to define mock.On call
//   - ctx context.Context
//   - clauses ...clause.Expression
func (_e *AbstractRepository_Expecter[T]) Find(ctx interface{}, clauses ...interface{}) *AbstractRepository_Find_Call[T] {
	return &AbstractRepository_Find_Call[T]{Call: _e.mock.On("Find",
		append([]interface{}{ctx}, clauses...)...)}
}

func (_c *AbstractRepository_Find_Call[T]) Run(run func(ctx context.Context, clauses ...clause.Expression)) *AbstractRepository_Find_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]clause.Expression, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(clause.Expression)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *AbstractRepository_Find_Call[T]) Return(_a0 []T, _a1 error) *AbstractRepository_Find_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AbstractRepository_Find_Call[T]) RunAndReturn(run func(context.Context, ...clause.Expression) ([]T, error)) *AbstractRepository_Find_Call[T] {
	_c.Call.Return(run)
	return _c
}

// FindInBatches provides a mock function with given fields: ctx, bacthSize, fn, clauses
func (_m *AbstractRepository[T]) FindInBatches(ctx context.Context, bacthSize int, fn func([]T) error, clauses ...clause.Expression) error {
	_va := make([]interface{}, len(clauses))
	for _i := range clauses {
		_va[_i] = clauses[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, bacthSize, fn)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, func([]T) error, ...clause.Expression) error); ok {
		r0 = rf(ctx, bacthSize, fn, clauses...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AbstractRepository_FindInBatches_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindInBatches'
type AbstractRepository_FindInBatches_Call[T gorm_provider.AbstractEntity] struct {
	*mock.Call
}

// FindInBatches is a helper method to define mock.On call
//   - ctx context.Context
//   - bacthSize int
//   - fn func([]T) error
//   - clauses ...clause.Expression
func (_e *AbstractRepository_Expecter[T]) FindInBatches(ctx interface{}, bacthSize interface{}, fn interface{}, clauses ...interface{}) *AbstractRepository_FindInBatches_Call[T] {
	return &AbstractRepository_FindInBatches_Call[T]{Call: _e.mock.On("FindInBatches",
		append([]interface{}{ctx, bacthSize, fn}, clauses...)...)}
}

func (_c *AbstractRepository_FindInBatches_Call[T]) Run(run func(ctx context.Context, bacthSize int, fn func([]T) error, clauses ...clause.Expression)) *AbstractRepository_FindInBatches_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]clause.Expression, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(clause.Expression)
			}
		}
		run(args[0].(context.Context), args[1].(int), args[2].(func([]T) error), variadicArgs...)
	})
	return _c
}

func (_c *AbstractRepository_FindInBatches_Call[T]) Return(_a0 error) *AbstractRepository_FindInBatches_Call[T] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AbstractRepository_FindInBatches_Call[T]) RunAndReturn(run func(context.Context, int, func([]T) error, ...clause.Expression) error) *AbstractRepository_FindInBatches_Call[T] {
	_c.Call.Return(run)
	return _c
}

// FindOne provides a mock function with given fields: ctx, clauses
func (_m *AbstractRepository[T]) FindOne(ctx context.Context, clauses ...clause.Expression) (T, error) {
	_va := make([]interface{}, len(clauses))
	for _i := range clauses {
		_va[_i] = clauses[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 T
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...clause.Expression) (T, error)); ok {
		return rf(ctx, clauses...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...clause.Expression) T); ok {
		r0 = rf(ctx, clauses...)
	} else {
		r0 = ret.Get(0).(T)
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...clause.Expression) error); ok {
		r1 = rf(ctx, clauses...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AbstractRepository_FindOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindOne'
type AbstractRepository_FindOne_Call[T gorm_provider.AbstractEntity] struct {
	*mock.Call
}

// FindOne is a helper method to define mock.On call
//   - ctx context.Context
//   - clauses ...clause.Expression
func (_e *AbstractRepository_Expecter[T]) FindOne(ctx interface{}, clauses ...interface{}) *AbstractRepository_FindOne_Call[T] {
	return &AbstractRepository_FindOne_Call[T]{Call: _e.mock.On("FindOne",
		append([]interface{}{ctx}, clauses...)...)}
}

func (_c *AbstractRepository_FindOne_Call[T]) Run(run func(ctx context.Context, clauses ...clause.Expression)) *AbstractRepository_FindOne_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]clause.Expression, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(clause.Expression)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *AbstractRepository_FindOne_Call[T]) Return(_a0 T, _a1 error) *AbstractRepository_FindOne_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AbstractRepository_FindOne_Call[T]) RunAndReturn(run func(context.Context, ...clause.Expression) (T, error)) *AbstractRepository_FindOne_Call[T] {
	_c.Call.Return(run)
	return _c
}

// First provides a mock function with given fields: ctx, clauses
func (_m *AbstractRepository[T]) First(ctx context.Context, clauses ...clause.Expression) (T, bool, error) {
	_va := make([]interface{}, len(clauses))
	for _i := range clauses {
		_va[_i] = clauses[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 T
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, ...clause.Expression) (T, bool, error)); ok {
		return rf(ctx, clauses...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...clause.Expression) T); ok {
		r0 = rf(ctx, clauses...)
	} else {
		r0 = ret.Get(0).(T)
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...clause.Expression) bool); ok {
		r1 = rf(ctx, clauses...)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(context.Context, ...clause.Expression) error); ok {
		r2 = rf(ctx, clauses...)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// AbstractRepository_First_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'First'
type AbstractRepository_First_Call[T gorm_provider.AbstractEntity] struct {
	*mock.Call
}

// First is a helper method to define mock.On call
//   - ctx context.Context
//   - clauses ...clause.Expression
func (_e *AbstractRepository_Expecter[T]) First(ctx interface{}, clauses ...interface{}) *AbstractRepository_First_Call[T] {
	return &AbstractRepository_First_Call[T]{Call: _e.mock.On("First",
		append([]interface{}{ctx}, clauses...)...)}
}

func (_c *AbstractRepository_First_Call[T]) Run(run func(ctx context.Context, clauses ...clause.Expression)) *AbstractRepository_First_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]clause.Expression, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(clause.Expression)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *AbstractRepository_First_Call[T]) Return(_a0 T, _a1 bool, _a2 error) *AbstractRepository_First_Call[T] {
	_c.Call.Return(_a0, _a1, _a2)
	return _c
}

func (_c *AbstractRepository_First_Call[T]) RunAndReturn(run func(context.Context, ...clause.Expression) (T, bool, error)) *AbstractRepository_First_Call[T] {
	_c.Call.Return(run)
	return _c
}

// NewQuery provides a mock function with given fields: ctx, clauses
func (_m *AbstractRepository[T]) NewQuery(ctx context.Context, clauses ...clause.Expression) gorm_provider.Query[T] {
	_va := make([]interface{}, len(clauses))
	for _i := range clauses {
		_va[_i] = clauses[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 gorm_provider.Query[T]
	if rf, ok := ret.Get(0).(func(context.Context, ...clause.Expression) gorm_provider.Query[T]); ok {
		r0 = rf(ctx, clauses...)
	} else {
		r0 = ret.Get(0).(gorm_provider.Query[T])
	}

	return r0
}

// AbstractRepository_NewQuery_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewQuery'
type AbstractRepository_NewQuery_Call[T gorm_provider.AbstractEntity] struct {
	*mock.Call
}

// NewQuery is a helper method to define mock.On call
//   - ctx context.Context
//   - clauses ...clause.Expression
func (_e *AbstractRepository_Expecter[T]) NewQuery(ctx interface{}, clauses ...interface{}) *AbstractRepository_NewQuery_Call[T] {
	return &AbstractRepository_NewQuery_Call[T]{Call: _e.mock.On("NewQuery",
		append([]interface{}{ctx}, clauses...)...)}
}

func (_c *AbstractRepository_NewQuery_Call[T]) Run(run func(ctx context.Context, clauses ...clause.Expression)) *AbstractRepository_NewQuery_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]clause.Expression, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(clause.Expression)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *AbstractRepository_NewQuery_Call[T]) Return(_a0 gorm_provider.Query[T]) *AbstractRepository_NewQuery_Call[T] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AbstractRepository_NewQuery_Call[T]) RunAndReturn(run func(context.Context, ...clause.Expression) gorm_provider.Query[T]) *AbstractRepository_NewQuery_Call[T] {
	_c.Call.Return(run)
	return _c
}

// NewTransaction provides a mock function with given fields: ctx, fc
func (_m *AbstractRepository[T]) NewTransaction(ctx context.Context, fc func(context.Context) error) error {
	ret := _m.Called(ctx, fc)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, func(context.Context) error) error); ok {
		r0 = rf(ctx, fc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AbstractRepository_NewTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NewTransaction'
type AbstractRepository_NewTransaction_Call[T gorm_provider.AbstractEntity] struct {
	*mock.Call
}

// NewTransaction is a helper method to define mock.On call
//   - ctx context.Context
//   - fc func(context.Context) error
func (_e *AbstractRepository_Expecter[T]) NewTransaction(ctx interface{}, fc interface{}) *AbstractRepository_NewTransaction_Call[T] {
	return &AbstractRepository_NewTransaction_Call[T]{Call: _e.mock.On("NewTransaction", ctx, fc)}
}

func (_c *AbstractRepository_NewTransaction_Call[T]) Run(run func(ctx context.Context, fc func(context.Context) error)) *AbstractRepository_NewTransaction_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(func(context.Context) error))
	})
	return _c
}

func (_c *AbstractRepository_NewTransaction_Call[T]) Return(_a0 error) *AbstractRepository_NewTransaction_Call[T] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AbstractRepository_NewTransaction_Call[T]) RunAndReturn(run func(context.Context, func(context.Context) error) error) *AbstractRepository_NewTransaction_Call[T] {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, update, clauses
func (_m *AbstractRepository[T]) Update(ctx context.Context, update interface{}, clauses ...clause.Expression) error {
	_va := make([]interface{}, len(clauses))
	for _i := range clauses {
		_va[_i] = clauses[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, update)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...clause.Expression) error); ok {
		r0 = rf(ctx, update, clauses...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AbstractRepository_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type AbstractRepository_Update_Call[T gorm_provider.AbstractEntity] struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - update interface{}
//   - clauses ...clause.Expression
func (_e *AbstractRepository_Expecter[T]) Update(ctx interface{}, update interface{}, clauses ...interface{}) *AbstractRepository_Update_Call[T] {
	return &AbstractRepository_Update_Call[T]{Call: _e.mock.On("Update",
		append([]interface{}{ctx, update}, clauses...)...)}
}

func (_c *AbstractRepository_Update_Call[T]) Run(run func(ctx context.Context, update interface{}, clauses ...clause.Expression)) *AbstractRepository_Update_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]clause.Expression, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(clause.Expression)
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *AbstractRepository_Update_Call[T]) Return(_a0 error) *AbstractRepository_Update_Call[T] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *AbstractRepository_Update_Call[T]) RunAndReturn(run func(context.Context, interface{}, ...clause.Expression) error) *AbstractRepository_Update_Call[T] {
	_c.Call.Return(run)
	return _c
}

// NewAbstractRepository creates a new instance of AbstractRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAbstractRepository[T gorm_provider.AbstractEntity](t interface {
	mock.TestingT
	Cleanup(func())
}) *AbstractRepository[T] {
	mock := &AbstractRepository[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
