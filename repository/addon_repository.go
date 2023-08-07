package repository

import (
	"courier-app/entity"

	"gorm.io/gorm"
)

type AddOnRepositoryConfig struct {
	DB *gorm.DB
}

type addOnRepositoryImpl struct {
	db *gorm.DB
}

type AddOnRepository interface {
	GetById(id uint) (*entity.AddOn, error)
	GetByIds(ids []uint) (*[]entity.AddOn, error)
	GetAll() (*[]entity.AddOn, error)
}

func NewAddOnRepository(ad AddOnRepositoryConfig) AddOnRepository {
	return &addOnRepositoryImpl{
		db: ad.DB,
	}
}

func (ad *addOnRepositoryImpl) GetById(id uint) (*entity.AddOn, error) {
	var addon *entity.AddOn
	err := ad.db.Where("id = ?", id).First(&addon).Error
	return addon, err
}

func (ad *addOnRepositoryImpl) GetByIds(ids []uint) (*[]entity.AddOn, error) {
	var addons *[]entity.AddOn
	err := ad.db.Where("id IN ?", ids).Find(&addons).Error
	return addons, err
}

func (ad *addOnRepositoryImpl) GetAll() (*[]entity.AddOn, error) {
	var addon *[]entity.AddOn
	err := ad.db.Find(&addon).Error
	return addon, err
}
