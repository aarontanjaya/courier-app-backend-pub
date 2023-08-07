// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	dto "courier-app/dto"
	entity "courier-app/entity"

	mock "github.com/stretchr/testify/mock"
)

// PaymentRepository is an autogenerated mock type for the PaymentRepository type
type PaymentRepository struct {
	mock.Mock
}

// GetPaymentDetail provides a mock function with given fields: paymentId
func (_m *PaymentRepository) GetPaymentDetail(paymentId uint) (*entity.Payment, error) {
	ret := _m.Called(paymentId)

	var r0 *entity.Payment
	if rf, ok := ret.Get(0).(func(uint) *entity.Payment); ok {
		r0 = rf(paymentId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Payment)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(paymentId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPaymentReport provides a mock function with given fields: req
func (_m *PaymentRepository) GetPaymentReport(req *dto.PeriodRequest) (*dto.PaymentReport, error) {
	ret := _m.Called(req)

	var r0 *dto.PaymentReport
	if rf, ok := ret.Get(0).(func(*dto.PeriodRequest) *dto.PaymentReport); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.PaymentReport)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*dto.PeriodRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Pay provides a mock function with given fields: userId, payment
func (_m *PaymentRepository) Pay(userId uint, payment *dto.PaymentRequest) (*entity.Payment, error) {
	ret := _m.Called(userId, payment)

	var r0 *entity.Payment
	if rf, ok := ret.Get(0).(func(uint, *dto.PaymentRequest) *entity.Payment); ok {
		r0 = rf(userId, payment)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Payment)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint, *dto.PaymentRequest) error); ok {
		r1 = rf(userId, payment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewPaymentRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewPaymentRepository creates a new instance of PaymentRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPaymentRepository(t mockConstructorTestingTNewPaymentRepository) *PaymentRepository {
	mock := &PaymentRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}