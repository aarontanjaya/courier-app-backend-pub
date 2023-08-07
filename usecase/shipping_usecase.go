package usecase

import (
	"courier-app/app_constant"
	"courier-app/dto"
	"courier-app/entity"
	"courier-app/repository"
	"courier-app/usecase/usecase_errors"
	"courier-app/util"
	"errors"

	"gorm.io/gorm"
)

type ShippingUsecase interface {
	CreateShipping(userId uint, req *dto.ShippingCreateRequest) (*entity.Shipping, error)
	GetByUserId(userId uint, paging *dto.ShippingTableRequest) (*[]entity.Shipping, *dto.Pagination, error)
	GetAll(paging *dto.ShippingTableRequest) (*[]entity.Shipping, *dto.Pagination, error)
	GetById(id uint) (*entity.Shipping, error)
	GetShippingStatuses() (*[]entity.ShippingStatus, error)
	UpdateStatus(shippingId uint, statusId uint) error
	Review(userId uint, req *dto.ShippingReviewRequest) error
}

type shippingUsecaseImpl struct {
	repository      repository.ShippingRepository
	categoryUsecase CategoryUsecase
	sizeUsecase     SizeUsecase
	addOnUsecase    AddOnUsecase
	addressUsecase  AddressUsecase
	userUsecase     UserUsecase
}

type ShippingUsecaseConfig struct {
	ShippingRepo    repository.ShippingRepository
	CategoryUsecase CategoryUsecase
	SizeUsecase     SizeUsecase
	AddOnUsecase    AddOnUsecase
	AddressUsease   AddressUsecase
	UserUsecase     UserUsecase
}

func NewShippingUsecase(s ShippingUsecaseConfig) ShippingUsecase {
	return &shippingUsecaseImpl{
		repository:      s.ShippingRepo,
		categoryUsecase: s.CategoryUsecase,
		sizeUsecase:     s.SizeUsecase,
		addOnUsecase:    s.AddOnUsecase,
		addressUsecase:  s.AddressUsease,
		userUsecase:     s.UserUsecase,
	}
}

func (s *shippingUsecaseImpl) GetById(id uint) (*entity.Shipping, error) {
	var shipping *entity.Shipping
	shipping, err := s.repository.GetById(id)
	if err == gorm.ErrRecordNotFound {
		return shipping, usecase_errors.ErrRecordNotExist
	}
	return shipping, err
}

func (s *shippingUsecaseImpl) GetByUserId(userId uint, paging *dto.ShippingTableRequest) (*[]entity.Shipping, *dto.Pagination, error) {
	var shippings *[]entity.Shipping
	var pagination *dto.Pagination
	if paging.PageSize <= 0 {
		paging.PageSize = util.DefaultPageSize
	}
	if paging.Page <= 0 {
		paging.Page = util.DefaultPage
	}
	shippings, err := s.repository.GetByUserId(userId, paging)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return shippings, pagination, nil
		}
		return shippings, pagination, err
	}

	count, err := s.repository.GetCountByUserId(userId, paging)
	if err != nil {
		return shippings, pagination, err
	}

	pagination = util.BuildPaginationObject(&dto.PaginationRequest{
		Search:   "",
		PageSize: paging.PageSize,
		Page:     paging.Page,
	}, count)
	return shippings, pagination, err
}

func (s *shippingUsecaseImpl) GetAll(paging *dto.ShippingTableRequest) (*[]entity.Shipping, *dto.Pagination, error) {
	var shippings *[]entity.Shipping
	var pagination *dto.Pagination
	if paging.PageSize <= 0 {
		paging.PageSize = util.DefaultPageSize
	}
	if paging.Page <= 0 {
		paging.Page = util.DefaultPage
	}
	shippings, err := s.repository.GetAll(paging)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return shippings, pagination, nil
		}
		return shippings, pagination, err
	}

	count, err := s.repository.GetCountAll(paging)
	if err != nil {
		return shippings, pagination, err
	}

	pagination = util.BuildPaginationObject(&dto.PaginationRequest{
		Search:   "",
		PageSize: paging.PageSize,
		Page:     paging.Page,
	}, count)
	return shippings, pagination, err
}

func (s *shippingUsecaseImpl) Review(userId uint, req *dto.ShippingReviewRequest) error {
	shipping, err := s.GetById(req.ShippingId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return usecase_errors.ErrRecordNotExist
		}
		return err
	}

	if shipping.UserId != userId {
		return usecase_errors.ErrUserNotAuthorized
	}

	if shipping.StatusId != app_constant.StatusShippingArrived {
		return usecase_errors.ErrInvalidReq
	}

	return s.repository.Review(req)
}

func (s *shippingUsecaseImpl) UpdateStatus(shippingId uint, statusId uint) error {
	shipping, err := s.GetById(shippingId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return usecase_errors.ErrRecordNotExist
		}
		return err
	}

	_, err = s.repository.GetShippingStatusById(statusId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return usecase_errors.ErrRecordNotExist
		}
		return err
	}

	err = s.repository.UpdateStatus(shippingId, statusId)
	if err != nil {
		return err
	}

	if shipping.StatusId != app_constant.StatusShippingArrived && statusId == app_constant.StatusShippingArrived {
		changeAmount := 1
		err = s.userUsecase.ChangeGachaQuota(shipping.UserId, changeAmount)
	}

	return err
}

func (s *shippingUsecaseImpl) CreateShipping(userId uint, req *dto.ShippingCreateRequest) (*entity.Shipping, error) {
	var shipping *entity.Shipping
	var shippingAddons []entity.AddOn
	category, err := s.categoryUsecase.GetById(req.CategoryId)
	if err != nil {
		return shipping, err
	}

	size, err := s.sizeUsecase.GetById(req.SizeId)
	if err != nil {
		return shipping, err
	}

	address, err := s.addressUsecase.GetById(req.AddressId)
	if err != nil {
		if address.UserId != userId {
			return shipping, usecase_errors.ErrUserNotAuthorized
		}
		return shipping, err
	}
	var addOnCost float64
	if len(req.AddOnIds) != 0 {
		addons, err := s.addOnUsecase.GetByIds(req.AddOnIds)
		if addons == nil {
			return shipping, usecase_errors.ErrRecordNotExist
		}
		if len(req.AddOnIds) != len(*addons) {
			return shipping, usecase_errors.ErrRecordNotExist
		}
		if err != nil {
			return shipping, err
		}
		addOnCost = calculateTotalAddOnCost(addons)
		shippingAddons = *addons
	}

	totalCost := addOnCost + size.Price + category.Price
	payment := entity.Payment{
		TotalCost: totalCost,
	}

	shipping, err = s.repository.CreateShipping(&entity.Shipping{
		SizeId:     req.SizeId,
		AddressId:  req.AddressId,
		CategoryId: req.CategoryId,
		StatusId:   app_constant.StatusShippingPayment,
		Payment:    payment,
		UserId:     userId,
		AddOn:      shippingAddons,
	})

	return shipping, err
}

func calculateTotalAddOnCost(addons *[]entity.AddOn) float64 {
	var totalCost float64
	for i := 0; i < len(*addons); i++ {
		totalCost += (*addons)[i].Price
	}
	return totalCost
}

func (s *shippingUsecaseImpl) GetShippingStatuses() (*[]entity.ShippingStatus, error) {
	var statuses *[]entity.ShippingStatus
	statuses, err := s.repository.GetShippingStatuses()
	if statuses == nil {
		return statuses, usecase_errors.ErrRecordNotExist
	}
	return statuses, err
}
