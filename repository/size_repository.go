package repository

import (
	"courier-app/entity"

	"gorm.io/gorm"
)

type SizeRepositoryConfig struct {
	DB *gorm.DB
}

type sizeRepositoryImpl struct {
	db *gorm.DB
}

type SizeRepository interface {
	GetById(id uint) (*entity.Size, error)
	GetAll() (*[]entity.Size, error)
}

func NewSizeRepository(s SizeRepositoryConfig) SizeRepository {
	return &sizeRepositoryImpl{
		db: s.DB,
	}
}

func (s *sizeRepositoryImpl) GetById(id uint) (*entity.Size, error) {
	var size *entity.Size
	err := s.db.Where("id = ?", id).First(&size).Error
	return size, err
}

func (s *sizeRepositoryImpl) GetAll() (*[]entity.Size, error) {
	var sizes *[]entity.Size
	err := s.db.Find(&sizes).Error
	return sizes, err
}
