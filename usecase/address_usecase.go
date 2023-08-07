package usecase

import (
	"courier-app/dto"
	"courier-app/entity"
	"courier-app/repository"
	"courier-app/usecase/usecase_errors"
	"courier-app/util"

	"gorm.io/gorm"
)

type AddressUsecase interface {
	GetAll(paging *dto.AddressPaginationRequest) (*[]entity.Address, *dto.Pagination, error)
	GetById(id uint) (*entity.Address, error)
	UpdateAddress(userId uint, address *dto.AddressUpdateRequest) error
	CreateAddress(userId uint, address *dto.AddressCreateRequest) (*entity.Address, error)
	DeleteAddress(recordId uint, userId uint) (*entity.Address, error)
}

type addressUsecaseImpl struct {
	repository repository.AddressRepository
}

type AddressUsecaseConfig struct {
	AddressRepo repository.AddressRepository
}

func NewAddressUsecase(c AddressUsecaseConfig) AddressUsecase {
	return &addressUsecaseImpl{
		repository: c.AddressRepo,
	}
}

func (ad *addressUsecaseImpl) GetById(id uint) (*entity.Address, error) {
	var address *entity.Address
	address, err := ad.repository.GetById(id)
	if err == gorm.ErrRecordNotFound {
		return address, usecase_errors.ErrRecordNotExist
	}

	return address, err
}

func (ad *addressUsecaseImpl) GetAll(p *dto.AddressPaginationRequest) (*[]entity.Address, *dto.Pagination, error) {
	var addresses *[]entity.Address
	var pagination *dto.Pagination
	paging := &dto.PaginationRequest{
		Search:   p.Search,
		Page:     p.Page,
		PageSize: p.PageSize,
	}
	util.ValidatePaginationRequest(paging)
	p.Search = paging.Search
	p.Page = paging.Page
	p.PageSize = paging.PageSize
	addresses, err := ad.repository.GetAll(p)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return addresses, pagination, nil
		}
		return addresses, pagination, err
	}

	count, err := ad.repository.GetCountAll(p)
	if err != nil {
		return addresses, pagination, err
	}

	pagination = util.BuildPaginationObject(paging, count)
	return addresses, pagination, err
}

func (ad *addressUsecaseImpl) UpdateAddress(userId uint, address *dto.AddressUpdateRequest) error {
	originalAddress, err := ad.repository.GetById(address.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return usecase_errors.ErrRecordNotExist
		}
		return err
	}
	if originalAddress.UserId != userId {
		return usecase_errors.ErrUserNotAuthorized
	}

	err = ad.repository.UpdateAddress(&entity.Address{
		ID:             address.ID,
		UserId:         userId,
		RecipientName:  address.RecipientName,
		RecipientPhone: address.RecipientPhone,
		FullAddress:    address.FullAddress,
		Label:          address.Label,
	})
	return err

}

func (ad *addressUsecaseImpl) CreateAddress(userId uint, address *dto.AddressCreateRequest) (*entity.Address, error) {
	var inputAddress *entity.Address
	inputAddress = &entity.Address{
		UserId:         userId,
		RecipientName:  address.RecipientName,
		Label:          address.Label,
		RecipientPhone: address.RecipientPhone,
		FullAddress:    address.FullAddress,
	}
	inputAddress, err := ad.repository.CreateAddress(inputAddress)
	return inputAddress, err
}

func (ad *addressUsecaseImpl) DeleteAddress(recordId uint, userId uint) (*entity.Address, error) {
	address, err := ad.repository.GetById(recordId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return address, usecase_errors.ErrRecordNotExist
		}
		return address, err
	}

	if address.UserId != userId {
		return address, usecase_errors.ErrUserNotAuthorized
	}

	inputAddress := &entity.Address{
		ID: recordId,
	}
	res, err := ad.repository.DeleteAddress(inputAddress)
	return res, err
}
