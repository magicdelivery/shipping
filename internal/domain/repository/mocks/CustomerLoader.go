// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"
	model "shipping/internal/domain/model"

	mock "github.com/stretchr/testify/mock"
)

// CustomerLoader is an autogenerated mock type for the CustomerLoader type
type CustomerLoader struct {
	mock.Mock
}

// LoadAllCustomers provides a mock function with given fields: ctx
func (_m *CustomerLoader) LoadAllCustomers(ctx context.Context) ([]*model.Customer, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for LoadAllCustomers")
	}

	var r0 []*model.Customer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*model.Customer, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*model.Customer); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Customer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadCustomerById provides a mock function with given fields: ctx, id
func (_m *CustomerLoader) LoadCustomerById(ctx context.Context, id string) (*model.Customer, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for LoadCustomerById")
	}

	var r0 *model.Customer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*model.Customer, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Customer); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Customer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCustomerLoader creates a new instance of CustomerLoader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCustomerLoader(t interface {
	mock.TestingT
	Cleanup(func())
}) *CustomerLoader {
	mock := &CustomerLoader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}