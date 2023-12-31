// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	entity "courier-app/entity"

	mock "github.com/stretchr/testify/mock"
)

// SizeRepository is an autogenerated mock type for the SizeRepository type
type SizeRepository struct {
	mock.Mock
}

// GetAll provides a mock function with given fields:
func (_m *SizeRepository) GetAll() (*[]entity.Size, error) {
	ret := _m.Called()

	var r0 *[]entity.Size
	if rf, ok := ret.Get(0).(func() *[]entity.Size); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]entity.Size)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetById provides a mock function with given fields: id
func (_m *SizeRepository) GetById(id uint) (*entity.Size, error) {
	ret := _m.Called(id)

	var r0 *entity.Size
	if rf, ok := ret.Get(0).(func(uint) *entity.Size); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Size)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewSizeRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewSizeRepository creates a new instance of SizeRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSizeRepository(t mockConstructorTestingTNewSizeRepository) *SizeRepository {
	mock := &SizeRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
