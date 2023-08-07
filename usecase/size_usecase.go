package usecase

import (
	"courier-app/entity"
	"courier-app/repository"
	"courier-app/usecase/usecase_errors"

	"gorm.io/gorm"
)

type SizeUsecase interface {
	GetById(id uint) (*entity.Size, error)
	GetAll() (*[]entity.Size, error)
}

type sizeUsecaseImpl struct {
	repository repository.SizeRepository
}

type SizeUsecaseConfig struct {
	SizeRepo repository.SizeRepository
}

func NewSizeUsecase(s SizeUsecaseConfig) SizeUsecase {
	return &sizeUsecaseImpl{
		repository: s.SizeRepo,
	}
}

func (s *sizeUsecaseImpl) GetById(id uint) (*entity.Size, error) {
	var size *entity.Size
	size, err := s.repository.GetById(id)
	if err == gorm.ErrRecordNotFound {
		return size, usecase_errors.ErrRecordNotExist
	}
	return size, err
}

func (s *sizeUsecaseImpl) GetAll() (*[]entity.Size, error) {
	var sizes *[]entity.Size
	sizes, err := s.repository.GetAll()
	if sizes == nil {
		return sizes, usecase_errors.ErrRecordNotExist
	}
	return sizes, err
}
