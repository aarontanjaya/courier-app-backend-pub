package usecase

import (
	"courier-app/entity"
	"courier-app/repository"
	"courier-app/usecase/usecase_errors"

	"gorm.io/gorm"
)

type AddOnUsecase interface {
	GetById(id uint) (*entity.AddOn, error)
	GetByIds(ids []uint) (*[]entity.AddOn, error)
	GetAll() (*[]entity.AddOn, error)
}

type addOnUsecaseImpl struct {
	repository repository.AddOnRepository
}

type AddOnUsecaseConfig struct {
	AddOnRepo repository.AddOnRepository
}

func NewAddOnUsecase(ad AddOnUsecaseConfig) AddOnUsecase {
	return &addOnUsecaseImpl{
		repository: ad.AddOnRepo,
	}
}

func (ad *addOnUsecaseImpl) GetById(id uint) (*entity.AddOn, error) {
	var addon *entity.AddOn
	addon, err := ad.repository.GetById(id)
	if err == gorm.ErrRecordNotFound {
		return addon, usecase_errors.ErrRecordNotExist
	}
	return addon, err
}

func (ad *addOnUsecaseImpl) GetByIds(ids []uint) (*[]entity.AddOn, error) {
	addons, err := ad.repository.GetByIds(ids)
	if addons == nil {
		return addons, usecase_errors.ErrRecordNotExist
	}
	return addons, err
}

func (ad *addOnUsecaseImpl) GetAll() (*[]entity.AddOn, error) {
	var addon *[]entity.AddOn
	addon, err := ad.repository.GetAll()
	if addon == nil {
		return addon, usecase_errors.ErrRecordNotExist
	}
	return addon, err
}
