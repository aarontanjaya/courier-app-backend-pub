package repository

import (
	"courier-app/entity"

	"gorm.io/gorm"
)

type CategoryRepositoryConfig struct {
	DB *gorm.DB
}

type categoryRepositoryImpl struct {
	db *gorm.DB
}

type CategoryRepository interface {
	GetById(id uint) (*entity.Category, error)
	GetAll() (*[]entity.Category, error)
}

func NewCategoryRepository(cat CategoryRepositoryConfig) CategoryRepository {
	return &categoryRepositoryImpl{
		db: cat.DB,
	}
}

func (cat *categoryRepositoryImpl) GetById(id uint) (*entity.Category, error) {
	var category *entity.Category
	err := cat.db.Where("id = ?", id).First(&category).Error
	return category, err
}

func (cat *categoryRepositoryImpl) GetAll() (*[]entity.Category, error) {
	var categories *[]entity.Category
	err := cat.db.Find(&categories).Error
	return categories, err
}
