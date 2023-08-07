package repository

import (
	"courier-app/dto"
	"courier-app/entity"
	"fmt"

	"gorm.io/gorm"
)

type PromoRepositoryConfig struct {
	DB *gorm.DB
}

type promoRepositoryImpl struct {
	db *gorm.DB
}

type PromoRepository interface {
	GetAllPromo(p *dto.PaginationSortRequest) (*[]entity.Promo, error)
	GetCountAllPromo(p *dto.PaginationSortRequest) (int64, error)
	GetAllActiveUserVouchers(userId uint) (*[]entity.UserVoucher, error)
	UseVoucher(id uint) error
	CreatePromo(promo *entity.Promo) (*entity.Promo, error)
	UpdatePromo(promo *entity.Promo) error
	GetPromoById(promoId uint) (*entity.Promo, error)
	GetRandomActivePromo() (*entity.Promo, error)
	IssueVoucher(userId uint, promoId uint) (*entity.UserVoucher, error)
	GetVoucherById(voucherId uint) (*entity.UserVoucher, error)
}

func NewPromoRepository(pr PromoRepositoryConfig) PromoRepository {
	return &promoRepositoryImpl{
		db: pr.DB,
	}
}

func (pr *promoRepositoryImpl) GetAllPromo(p *dto.PaginationSortRequest) (*[]entity.Promo, error) {
	var promos *[]entity.Promo
	orderString := fmt.Sprintf("%s %s", p.SortField, p.SortDirection)
	err := pr.db.Where("name ILIKE ?", p.Search).Order(orderString).Limit(int(p.PageSize)).Find(&promos).Error
	return promos, err
}

func (pr *promoRepositoryImpl) GetCountAllPromo(p *dto.PaginationSortRequest) (int64, error) {
	var count int64
	err := pr.db.Model(&entity.Promo{}).Where("name ILIKE ?", p.Search).Count(&count).Error
	return count, err
}

func (pr *promoRepositoryImpl) GetPromoById(promoId uint) (*entity.Promo, error) {
	var promo *entity.Promo
	err := pr.db.Where("id = ?", promoId).First(&promo).Error
	return promo, err
}

func (pr *promoRepositoryImpl) GetRandomActivePromo() (*entity.Promo, error) {
	var promo *entity.Promo
	err := pr.db.Where("exp_date > NOW() AND (quota >0 OR limited = false)").Order("random()").First(&promo).Error
	return promo, err
}

func (pr *promoRepositoryImpl) CreatePromo(promo *entity.Promo) (*entity.Promo, error) {
	err := pr.db.Create(&promo).Error
	return promo, err
}

func (pr *promoRepositoryImpl) UpdatePromo(promo *entity.Promo) error {
	err := pr.db.Model(&entity.Promo{}).Where("id = ?", promo.ID).Updates(entity.Promo{
		Name:        promo.Name,
		Discount:    promo.Discount,
		MinFee:      promo.MinFee,
		MaxDiscount: promo.MaxDiscount,
		Quota:       promo.Quota,
		ExpDate:     promo.ExpDate,
		Limited:     promo.Limited,
	}).Error
	return err
}

func (pr *promoRepositoryImpl) IssueVoucher(userId uint, promoId uint) (*entity.UserVoucher, error) {
	var promo *entity.Promo
	voucher := &entity.UserVoucher{
		PromoId: promoId,
		UserId:  userId,
	}
	err := pr.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", promoId).First(&promo).Error; err != nil {
			return err
		}

		voucher.ExpDate = promo.ExpDate
		if err := tx.Create(&voucher).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(&entity.Promo{}).Where("id = ?", promoId).Update("quota", gorm.Expr("quota - 1")).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(&entity.UserDetail{}).Where("user_id = ? AND gacha_quota > 0", userId).Update("gacha_quota", gorm.Expr("gacha_quota - 1")).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})

	return voucher, err
}

func (pr *promoRepositoryImpl) GetAllActiveUserVouchers(userId uint) (*[]entity.UserVoucher, error) {
	var vouchers *[]entity.UserVoucher
	err := pr.db.Model(&entity.UserVoucher{}).Preload("Promo").Where("user_id = ? AND exp_date >= NOW()", userId).Find(&vouchers).Error
	return vouchers, err
}

func (pr *promoRepositoryImpl) GetVoucherById(voucherId uint) (*entity.UserVoucher, error) {
	var voucher *entity.UserVoucher
	err := pr.db.Preload("Promo").Where("id = ? AND exp_date >= NOW()", voucherId).First(&voucher).Error
	return voucher, err
}

func (pr *promoRepositoryImpl) UseVoucher(id uint) error {
	return pr.db.Delete(&entity.UserVoucher{
		ID: id,
	}).Error
}
