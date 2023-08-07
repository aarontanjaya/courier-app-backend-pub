// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	entity "courier-app/entity"

	mock "github.com/stretchr/testify/mock"
)

// CategoryUsecase is an autogenerated mock type for the CategoryUsecase type
type CategoryUsecase struct {
	mock.Mock
}

// GetAll provides a mock function with given fields:
func (_m *CategoryUsecase) GetAll() (*[]entity.Category, error) {
	ret := _m.Called()

	var r0 *[]entity.Category
	if rf, ok := ret.Get(0).(func() *[]entity.Category); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]entity.Category)
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
func (_m *CategoryUsecase) GetById(id uint) (*entity.Category, error) {
	ret := _m.Called(id)

	var r0 *entity.Category
	if rf, ok := ret.Get(0).(func(uint) *entity.Category); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Category)
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

type mockConstructorTestingTNewCategoryUsecase interface {
	mock.TestingT
	Cleanup(func())
}

// NewCategoryUsecase creates a new instance of CategoryUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCategoryUsecase(t mockConstructorTestingTNewCategoryUsecase) *CategoryUsecase {
	mock := &CategoryUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}