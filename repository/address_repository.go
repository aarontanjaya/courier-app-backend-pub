package repository

import (
	"courier-app/dto"
	"courier-app/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AddressRepositoryConfig struct {
	DB *gorm.DB
}

type addressRepositoryImpl struct {
	db *gorm.DB
}

type AddressRepository interface {
	GetAll(paging *dto.AddressPaginationRequest) (*[]entity.Address, error)
	GetCountAll(paging *dto.AddressPaginationRequest) (int64, error)
	UpdateAddress(address *entity.Address) error
	GetById(id uint) (*entity.Address, error)
	CreateAddress(address *entity.Address) (*entity.Address, error)
	DeleteAddress(address *entity.Address) (*entity.Address, error)
}

func NewAddressRepository(ad AddressRepositoryConfig) AddressRepository {
	return &addressRepositoryImpl{
		db: ad.DB,
	}
}

func (ad *addressRepositoryImpl) GetById(id uint) (*entity.Address, error) {
	var address *entity.Address
	err := ad.db.Where("id = ?", id).First(&address).Error
	return address, err
}

func (ad *addressRepositoryImpl) GetAll(paging *dto.AddressPaginationRequest) (*[]entity.Address, error) {
	var addresses *[]entity.Address
	offset := int((paging.Page - 1) * paging.PageSize)
	query := ad.db.Limit(int(paging.PageSize)).Offset(offset)
	if paging.UserId != nil && *paging.UserId != 0 {
		query = query.Where("user_id = ?", paging.UserId)
	}
	err := query.Where("((recipient_name ILIKE ?) OR (full_address ILIKE ?) OR (label ILIKE ?)) ", paging.Search, paging.Search, paging.Search).Find(&addresses).Error
	return addresses, err
}

func (ad *addressRepositoryImpl) GetCountByUserId(id uint, paging *dto.PaginationRequest) (int64, error) {
	var count int64
	err := ad.db.Model(&entity.Address{}).Where("user_id = ? AND ((recipient_name ILIKE ?) OR (full_address ILIKE ?) OR (label ILIKE ?)) ", id, paging.Search, paging.Search, paging.Search).Count(&count).Error
	return count, err
}

func (ad *addressRepositoryImpl) GetCountAll(paging *dto.AddressPaginationRequest) (int64, error) {
	var count int64
	query := ad.db.Model(&entity.Address{})
	if paging.UserId != nil && *paging.UserId != 0 {
		query = query.Where("user_id = ?", paging.UserId)
	}
	err := query.Where("((recipient_name ILIKE ?) OR (full_address ILIKE ?) OR (label ILIKE ?)) ", paging.Search, paging.Search, paging.Search).Count(&count).Error
	return count, err
}

func (ad *addressRepositoryImpl) UpdateAddress(address *entity.Address) error {
	err := ad.db.Model(&entity.Address{}).Where("id = ?", address.ID).Updates(address).Error
	return err
}

func (ad *addressRepositoryImpl) CreateAddress(address *entity.Address) (*entity.Address, error) {
	err := ad.db.Create(&address).Error
	return address, err
}

func (ad *addressRepositoryImpl) DeleteAddress(address *entity.Address) (*entity.Address, error) {
	err := ad.db.Clauses(clause.Returning{}).Delete(&address).Error
	return address, err
}
