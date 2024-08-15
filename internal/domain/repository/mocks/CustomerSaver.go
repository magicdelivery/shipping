// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"
	model "shipping/internal/domain/model"

	mock "github.com/stretchr/testify/mock"
)

// CustomerSaver is an autogenerated mock type for the CustomerSaver type
type CustomerSaver struct {
	mock.Mock
}

// SaveCustomer provides a mock function with given fields: ctx, shipping
func (_m *CustomerSaver) SaveCustomer(ctx context.Context, shipping model.Customer) error {
	ret := _m.Called(ctx, shipping)

	if len(ret) == 0 {
		panic("no return value specified for SaveCustomer")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.Customer) error); ok {
		r0 = rf(ctx, shipping)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewCustomerSaver creates a new instance of CustomerSaver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCustomerSaver(t interface {
	mock.TestingT
	Cleanup(func())
}) *CustomerSaver {
	mock := &CustomerSaver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}