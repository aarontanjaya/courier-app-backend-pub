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

func TestCategoryGetById(t *testing.T) {
	t.Run("should return category", func(t *testing.T) {
		userId := uint(1)
		mockResponse := &entity.Category{
			Name:        "tes",
			Description: "tes",
			Price:       10000,
			ID:          userId,
		}
		mockRepo := &mocks.CategoryRepository{}
		mockUc := usecase.NewCategoryUsecase(usecase.CategoryUsecaseConfig{
			CategoryRepo: mockRepo,
		})
		mockRepo.On("GetById", userId).Return(mockResponse, nil)

		actualRes, err := mockUc.GetById(userId)

		assert.Nil(t, err)
		assert.Equal(t, mockResponse, actualRes)
	})

	t.Run("should return error if not exist", func(t *testing.T) {
		userId := uint(1)
		mockRepo := &mocks.CategoryRepository{}
		mockUc := usecase.NewCategoryUsecase(usecase.CategoryUsecaseConfig{
			CategoryRepo: mockRepo,
		})
		mockRepo.On("GetById", userId).Return(nil, gorm.ErrRecordNotFound)

		actualRes, err := mockUc.GetById(userId)

		assert.EqualError(t, err, usecase_errors.ErrRecordNotExist.Error())
		assert.Nil(t, actualRes)
	})

}

func TestCategoryGetAll(t *testing.T) {
	t.Run("should return all category", func(t *testing.T) {
		userId := uint(1)
		mockResponse := &[]entity.Category{{
			Name:        "tes",
			Description: "tes",
			Price:       10000,
			ID:          userId}, {
			Name:        "tes",
			Description: "tes",
			Price:       10000,
			ID:          userId + 1},
		}
		mockRepo := &mocks.CategoryRepository{}
		mockUc := usecase.NewCategoryUsecase(usecase.CategoryUsecaseConfig{
			CategoryRepo: mockRepo,
		})

		mockRepo.On("GetAll").Return(mockResponse, nil)
		actualRes, err := mockUc.GetAll()

		assert.Nil(t, err)
		assert.Equal(t, mockResponse, actualRes)
	})

	t.Run("should return error if nothing returned", func(t *testing.T) {
		mockRepo := &mocks.CategoryRepository{}
		mockUc := usecase.NewCategoryUsecase(usecase.CategoryUsecaseConfig{
			CategoryRepo: mockRepo,
		})

		mockRepo.On("GetAll").Return(nil, nil)
		actualRes, err := mockUc.GetAll()

		assert.EqualError(t, err, usecase_errors.ErrRecordNotExist.Error())
		assert.Nil(t, actualRes)
	})
}
