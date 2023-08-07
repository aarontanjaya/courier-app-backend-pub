package usecase

import (
	"courier-app/entity"
	"courier-app/repository"
	"courier-app/usecase/usecase_errors"

	"gorm.io/gorm"
)

type CategoryUsecase interface {
	GetById(id uint) (*entity.Category, error)
	GetAll() (*[]entity.Category, error)
}

type categoryUsecaseImpl struct {
	repository repository.CategoryRepository
}

type CategoryUsecaseConfig struct {
	CategoryRepo repository.CategoryRepository
}

func NewCategoryUsecase(c CategoryUsecaseConfig) CategoryUsecase {
	return &categoryUsecaseImpl{
		repository: c.CategoryRepo,
	}
}

func (c *categoryUsecaseImpl) GetById(id uint) (*entity.Category, error) {
	var category *entity.Category
	category, err := c.repository.GetById(id)
	if err == gorm.ErrRecordNotFound {
		return category, usecase_errors.ErrRecordNotExist
	}
	return category, err
}

func (c *categoryUsecaseImpl) GetAll() (*[]entity.Category, error) {
	var categories *[]entity.Category
	categories, err := c.repository.GetAll()
	if categories == nil {
		return categories, usecase_errors.ErrRecordNotExist
	}
	return categories, err
}
