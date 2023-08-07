package repository

import (
	"courier-app/dto"
	"courier-app/entity"

	"gorm.io/gorm"
)

type ShippingRepositoryConfig struct {
	DB *gorm.DB
}

type shippingRepositoryImpl struct {
	db *gorm.DB
}

type ShippingRepository interface {
	CreateShipping(*entity.Shipping) (*entity.Shipping, error)
	UpdateStatus(shippingId uint, statusId uint) error
	Review(req *dto.ShippingReviewRequest) error
	GetByUserId(userId uint, paging *dto.ShippingTableRequest) (*[]entity.Shipping, error)
	GetAll(paging *dto.ShippingTableRequest) (*[]entity.Shipping, error)
	GetCountByUserId(userId uint, paging *dto.ShippingTableRequest) (int64, error)
	GetCountAll(paging *dto.ShippingTableRequest) (int64, error)
	GetById(id uint) (*entity.Shipping, error)
	GetShippingStatusById(statusId uint) (*entity.ShippingStatus, error)
	GetShippingStatuses() (*[]entity.ShippingStatus, error)
}

func NewShippingRepository(sp ShippingRepositoryConfig) ShippingRepository {
	return &shippingRepositoryImpl{
		db: sp.DB,
	}
}

func (sp *shippingRepositoryImpl) CreateShipping(shipping *entity.Shipping) (*entity.Shipping, error) {
	err := sp.db.Create(&shipping).Error
	return shipping, err
}

func (sp *shippingRepositoryImpl) UpdateStatus(shippingId uint, statusId uint) error {
	err := sp.db.Model(&entity.Shipping{}).Where("id = ?", shippingId).Update("status_id", statusId).Error
	return err
}

func (sp *shippingRepositoryImpl) Review(req *dto.ShippingReviewRequest) error {
	return sp.db.Model(&entity.Shipping{}).Where("id = ?", req.ShippingId).Updates(entity.Shipping{
		ReviewComment: req.ReviewComment,
		ReviewRating:  int(req.ReviewRating),
	}).Error
}

func (sp *shippingRepositoryImpl) GetById(id uint) (*entity.Shipping, error) {
	var shipping *entity.Shipping
	err := sp.db.Preload("Payment").Preload("AddOn").Preload("Category").Preload("Address", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("Size").Preload("ShippingStatus").Where("id = ?", id).First(&shipping).Error
	return shipping, err
}

func (sp *shippingRepositoryImpl) GetByUserId(userId uint, paging *dto.ShippingTableRequest) (*[]entity.Shipping, error) {
	var shippings *[]entity.Shipping
	offset := int((paging.Page - 1) * paging.PageSize)
	query := sp.db.Preload("Payment").Preload("AddOn").Preload("Category").Preload("Payment.TransactionRecord").Preload("Address", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("Size").Preload("ShippingStatus").Where("user_id = ?", userId).Limit(int(paging.PageSize)).Offset(offset)
	if len(paging.SizeIds) != 0 && paging.SizeIds[0] != 0 {
		query = query.Where("size_id IN ?", paging.SizeIds)
	}
	if len(paging.CategoryIds) != 0 && paging.CategoryIds[0] != 0 {
		query = query.Where("category_id IN ?", paging.CategoryIds)
	}
	if len(paging.StatusIds) != 0 && paging.StatusIds[0] != 0 {
		query = query.Where("status_id IN ?", paging.StatusIds)
	}
	err := query.Order("status_id asc, updated_at desc, created_at desc").Find(&shippings).Error
	return shippings, err
}

func (sp *shippingRepositoryImpl) GetAll(paging *dto.ShippingTableRequest) (*[]entity.Shipping, error) {
	var shippings *[]entity.Shipping
	offset := int((paging.Page - 1) * paging.PageSize)
	query := sp.db.Preload("Payment").Preload("AddOn").Preload("Category").Preload("Payment.TransactionRecord").Preload("Address", func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}).Preload("Size").Preload("ShippingStatus").Limit(int(paging.PageSize)).Offset(offset)
	if len(paging.SizeIds) != 0 && paging.SizeIds[0] != 0 {
		query = query.Where("size_id IN ?", paging.SizeIds)
	}
	if len(paging.CategoryIds) != 0 && paging.CategoryIds[0] != 0 {
		query = query.Where("category_id IN ?", paging.CategoryIds)
	}
	if len(paging.StatusIds) != 0 && paging.StatusIds[0] != 0 {
		query = query.Where("status_id IN ?", paging.StatusIds)
	}
	err := query.Order("status_id asc, updated_at desc, created_at desc").Find(&shippings).Error
	return shippings, err
}

func (sp *shippingRepositoryImpl) GetCountByUserId(userId uint, paging *dto.ShippingTableRequest) (int64, error) {
	var count int64
	query := sp.db.Model(&entity.Shipping{}).Where("user_id = ?", userId)
	if len(paging.SizeIds) != 0 && paging.SizeIds[0] != 0 {
		query = query.Where("size_id IN ?", paging.SizeIds)
	}
	if len(paging.CategoryIds) != 0 && paging.CategoryIds[0] != 0 {
		query = query.Where("category_id IN ?", paging.CategoryIds)
	}
	if len(paging.StatusIds) != 0 && paging.StatusIds[0] != 0 {
		query = query.Where("status_id IN ?", paging.StatusIds)
	}
	err := query.Count(&count).Error
	return count, err
}

func (sp *shippingRepositoryImpl) GetCountAll(paging *dto.ShippingTableRequest) (int64, error) {
	var count int64
	query := sp.db.Model(&entity.Shipping{})
	if len(paging.SizeIds) != 0 && paging.SizeIds[0] != 0 {
		query = query.Where("size_id IN ?", paging.SizeIds)
	}
	if len(paging.CategoryIds) != 0 && paging.CategoryIds[0] != 0 {
		query = query.Where("category_id IN ?", paging.CategoryIds)
	}
	if len(paging.StatusIds) != 0 && paging.StatusIds[0] != 0 {
		query = query.Where("status_id IN ?", paging.StatusIds)
	}
	err := query.Count(&count).Error
	return count, err
}

func (sp *shippingRepositoryImpl) GetShippingStatuses() (*[]entity.ShippingStatus, error) {
	var statuses *[]entity.ShippingStatus
	err := sp.db.Find(&statuses).Error
	return statuses, err
}

func (sp *shippingRepositoryImpl) GetShippingStatusById(statusId uint) (*entity.ShippingStatus, error) {
	var status *entity.ShippingStatus
	err := sp.db.Where("id = ?", statusId).First(&status).Error
	return status, err
}
