package usecase_test

import (
	"courier-app/entity"
	mocks "courier-app/mocks/repository"
	"courier-app/usecase"
	"courier-app/usecase/usecase_errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSizeGetById(t *testing.T) {
	t.Run("should return size", func(t *testing.T) {
		userId := uint(1)
		mockResponse := &entity.Size{
			Name:        "tes",
			Description: "tes",
			Price:       10000,
			ID:          userId,
		}
		mockRepo := &mocks.SizeRepository{}
		mockUc := usecase.NewSizeUsecase(usecase.SizeUsecaseConfig{
			SizeRepo: mockRepo,
		})
		mockRepo.On("GetById", userId).Return(mockResponse, nil)

		actualRes, err := mockUc.GetById(userId)

		assert.Nil(t, err)
		assert.Equal(t, mockResponse, actualRes)
	})

	t.Run("should return error if not exist", func(t *testing.T) {
		userId := uint(1)
		mockRepo := &mocks.SizeRepository{}
		mockUc := usecase.NewSizeUsecase(usecase.SizeUsecaseConfig{
			SizeRepo: mockRepo,
		})
		mockRepo.On("GetById", userId).Return(nil, gorm.ErrRecordNotFound)

		actualRes, err := mockUc.GetById(userId)

		assert.EqualError(t, err, usecase_errors.ErrRecordNotExist.Error())
		assert.Nil(t, actualRes)
	})
}

func TestSizeGetAll(t *testing.T) {
	t.Run("should return all size", func(t *testing.T) {
		userId := uint(1)
		mockResponse := &[]entity.Size{{
			Name:        "tes",
			Description: "tes",
			Price:       10000,
			ID:          userId}, {
			Name:        "tes",
			Description: "tes",
			Price:       10000,
			ID:          userId + 1},
		}
		mockRepo := &mocks.SizeRepository{}
		mockUc := usecase.NewSizeUsecase(usecase.SizeUsecaseConfig{
			SizeRepo: mockRepo,
		})

		mockRepo.On("GetAll").Return(mockResponse, nil)
		actualRes, err := mockUc.GetAll()

		assert.Nil(t, err)
		assert.Equal(t, mockResponse, actualRes)
	})

	t.Run("should return error if nothing returned", func(t *testing.T) {
		mockRepo := &mocks.SizeRepository{}
		mockUc := usecase.NewSizeUsecase(usecase.SizeUsecaseConfig{
			SizeRepo: mockRepo,
		})

		mockRepo.On("GetAll").Return(nil, nil)
		actualRes, err := mockUc.GetAll()

		assert.EqualError(t, err, usecase_errors.ErrRecordNotExist.Error())
		assert.Nil(t, actualRes)
	})
}
