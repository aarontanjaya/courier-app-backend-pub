package usecase

import (
	"courier-app/dto"
	"courier-app/entity"
	"courier-app/repository"
	"courier-app/usecase/usecase_errors"
	"courier-app/util"
	"errors"

	"gorm.io/gorm"
)

type PromoUsecase interface {
	GetAllPromo(*dto.PaginationSortRequest) (*[]entity.Promo, *dto.Pagination, error)
	CreatePromo(promo *dto.PromoRequest) (*entity.Promo, error)
	UpdatePromo(promoId uint, p *dto.PromoRequest) error
	IssueGachaVoucher(userId uint) (*entity.Promo, error)
	GetAllActiveUserVouchers(userId uint) (*[]entity.UserVoucher, error)
	GetVoucherById(voucherId uint) (*entity.UserVoucher, error)
}

type promoUsecaseImpl struct {
	repository  repository.PromoRepository
	userUsecase UserUsecase
}

type PromoUsecaseConfig struct {
	PromoRepo   repository.PromoRepository
	UserUsecase UserUsecase
}

func NewPromoUsecase(pr PromoUsecaseConfig) PromoUsecase {
	return &promoUsecaseImpl{
		repository:  pr.PromoRepo,
		userUsecase: pr.UserUsecase,
	}
}

func (pr *promoUsecaseImpl) GetAllPromo(p *dto.PaginationSortRequest) (*[]entity.Promo, *dto.Pagination, error) {
	var promos *[]entity.Promo
	var pagination *dto.Pagination
	err := validateAndFormatPagingFields(p)
	if err != nil {
		return promos, pagination, err
	}

	promos, err = pr.repository.GetAllPromo(p)
	if err != nil {
		return promos, pagination, err
	}

	count, err := pr.repository.GetCountAllPromo(p)
	if err != nil {
		return promos, pagination, err
	}

	paging := &dto.PaginationRequest{
		Search:   p.Search,
		Page:     p.Page,
		PageSize: p.PageSize,
	}
	pagination = util.BuildPaginationObject(paging, count)

	return promos, pagination, err
}

func validateAndFormatPagingFields(paging *dto.PaginationSortRequest) error {
	if paging.SortDirection == "" {
		paging.SortDirection = "desc"
	}
	err := validateSortDirection(paging.SortDirection)
	if err != nil {
		return err
	}

	if paging.SortField == "" {
		paging.SortField = "exp_date"
	}
	err = validateSortField(paging.SortField)
	if err != nil {
		return err
	}
	paging.Search = "%" + paging.Search + "%"
	return err
}

func validateSortDirection(direction string) error {
	switch direction {
	case "asc", "desc":
		return nil
	default:
		return usecase_errors.ErrPaginationInvalid
	}
}

func validateSortField(field string) error {
	switch field {
	case "quota", "exp_date":
		return nil
	default:
		return usecase_errors.ErrPaginationInvalid
	}
}

func (pr *promoUsecaseImpl) CreatePromo(promo *dto.PromoRequest) (*entity.Promo, error) {
	var resPromo *entity.Promo
	if err := validatePromoRequest(promo); err != nil {
		return resPromo, err
	}

	resPromo = &entity.Promo{
		Name:        promo.Name,
		MinFee:      promo.MinFee,
		Discount:    promo.Discount,
		MaxDiscount: promo.MaxDiscount,
		Quota:       promo.Quota,
		Limited:     *promo.Limited,
		ExpDate:     promo.ExpDate,
	}

	return pr.repository.CreatePromo(resPromo)
}

func (pr *promoUsecaseImpl) UpdatePromo(promoId uint, p *dto.PromoRequest) error {
	if err := validatePromoRequest(p); err != nil {
		return err
	}
	_, err := pr.repository.GetPromoById(promoId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return usecase_errors.ErrRecordNotExist
		}
	}
	promo := &entity.Promo{
		ID:          promoId,
		Name:        p.Name,
		MinFee:      p.MinFee,
		Discount:    p.Discount,
		MaxDiscount: p.MaxDiscount,
		Quota:       p.Quota,
		Limited:     *p.Limited,
		ExpDate:     p.ExpDate,
	}

	return pr.repository.UpdatePromo(promo)
}

func (pr *promoUsecaseImpl) IssueGachaVoucher(userId uint) (*entity.Promo, error) {
	var voucher *entity.Promo
	user, err := pr.userUsecase.GetById(userId)
	if err != nil {
		return voucher, err
	}
	if user.GachaQuota < 1 {
		return voucher, usecase_errors.ErrInsufficientGachaQuota
	}

	voucher, err = pr.repository.GetRandomActivePromo()
	if err != nil {
		return voucher, err
	}

	_, err = pr.repository.IssueVoucher(userId, voucher.ID)
	return voucher, err
}

func (pr *promoUsecaseImpl) GetAllActiveUserVouchers(userId uint) (*[]entity.UserVoucher, error) {
	return pr.repository.GetAllActiveUserVouchers(userId)
}

func (pr *promoUsecaseImpl) GetVoucherById(voucherId uint) (*entity.UserVoucher, error) {
	var voucher *entity.UserVoucher
	voucher, err := pr.repository.GetVoucherById(voucherId)
	if err == gorm.ErrRecordNotFound {
		return voucher, usecase_errors.ErrRecordNotExist
	}
	return voucher, err
}

func validatePromoRequest(promo *dto.PromoRequest) error {
	if promo.Discount < 0 {
		return usecase_errors.ErrPromoDiscountNegative
	}
	if promo.Discount > 1 {
		return usecase_errors.ErrPromoDiscountMax
	}
	if promo.MinFee < 0 {
		return usecase_errors.ErrPromoMinFeeNegative
	}
	if promo.MaxDiscount < 0 {
		return usecase_errors.ErrMaxDiscountNegative
	}
	if promo.Quota < 0 && *promo.Limited {
		return usecase_errors.ErrQuotaNegative
	}
	return nil
}
