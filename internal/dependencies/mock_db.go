// Code generated by mockery v2.20.0. DO NOT EDIT.

package dependencies

import (
	context "context"

	pgx "github.com/jackc/pgx/v4"
	mock "github.com/stretchr/testify/mock"
)

// MockCoreStorageInterface is an autogenerated mock type for the CoreStorageInterface type
type MockCoreStorageInterface struct {
	mock.Mock
}

type MockCoreStorageInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *MockCoreStorageInterface) EXPECT() *MockCoreStorageInterface_Expecter {
	return &MockCoreStorageInterface_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with given fields:
func (_m *MockCoreStorageInterface) Close() {
	_m.Called()
}

// MockCoreStorageInterface_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type MockCoreStorageInterface_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *MockCoreStorageInterface_Expecter) Close() *MockCoreStorageInterface_Close_Call {
	return &MockCoreStorageInterface_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *MockCoreStorageInterface_Close_Call) Run(run func()) *MockCoreStorageInterface_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockCoreStorageInterface_Close_Call) Return() *MockCoreStorageInterface_Close_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockCoreStorageInterface_Close_Call) RunAndReturn(run func()) *MockCoreStorageInterface_Close_Call {
	_c.Call.Return(run)
	return _c
}

// Query provides a mock function with given fields: ctx, sql, args
func (_m *MockCoreStorageInterface) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, sql)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 pgx.Rows
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) (pgx.Rows, error)); ok {
		return rf(ctx, sql, args...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) pgx.Rows); ok {
		r0 = rf(ctx, sql, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pgx.Rows)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(ctx, sql, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockCoreStorageInterface_Query_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Query'
type MockCoreStorageInterface_Query_Call struct {
	*mock.Call
}

// Query is a helper method to define mock.On call
//  - ctx context.Context
//  - sql string
//  - args ...interface{}
func (_e *MockCoreStorageInterface_Expecter) Query(ctx interface{}, sql interface{}, args ...interface{}) *MockCoreStorageInterface_Query_Call {
	return &MockCoreStorageInterface_Query_Call{Call: _e.mock.On("Query",
		append([]interface{}{ctx, sql}, args...)...)}
}

func (_c *MockCoreStorageInterface_Query_Call) Run(run func(ctx context.Context, sql string, args ...interface{})) *MockCoreStorageInterface_Query_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockCoreStorageInterface_Query_Call) Return(_a0 pgx.Rows, _a1 error) *MockCoreStorageInterface_Query_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockCoreStorageInterface_Query_Call) RunAndReturn(run func(context.Context, string, ...interface{}) (pgx.Rows, error)) *MockCoreStorageInterface_Query_Call {
	_c.Call.Return(run)
	return _c
}

// QueryRow provides a mock function with given fields: ctx, sql, args
func (_m *MockCoreStorageInterface) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	var _ca []interface{}
	_ca = append(_ca, ctx, sql)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 pgx.Row
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) pgx.Row); ok {
		r0 = rf(ctx, sql, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(pgx.Row)
		}
	}

	return r0
}

// MockCoreStorageInterface_QueryRow_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'QueryRow'
type MockCoreStorageInterface_QueryRow_Call struct {
	*mock.Call
}

// QueryRow is a helper method to define mock.On call
//  - ctx context.Context
//  - sql string
//  - args ...interface{}
func (_e *MockCoreStorageInterface_Expecter) QueryRow(ctx interface{}, sql interface{}, args ...interface{}) *MockCoreStorageInterface_QueryRow_Call {
	return &MockCoreStorageInterface_QueryRow_Call{Call: _e.mock.On("QueryRow",
		append([]interface{}{ctx, sql}, args...)...)}
}

func (_c *MockCoreStorageInterface_QueryRow_Call) Run(run func(ctx context.Context, sql string, args ...interface{})) *MockCoreStorageInterface_QueryRow_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockCoreStorageInterface_QueryRow_Call) Return(_a0 pgx.Row) *MockCoreStorageInterface_QueryRow_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockCoreStorageInterface_QueryRow_Call) RunAndReturn(run func(context.Context, string, ...interface{}) pgx.Row) *MockCoreStorageInterface_QueryRow_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockCoreStorageInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockCoreStorageInterface creates a new instance of MockCoreStorageInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockCoreStorageInterface(t mockConstructorTestingTNewMockCoreStorageInterface) *MockCoreStorageInterface {
	mock := &MockCoreStorageInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
