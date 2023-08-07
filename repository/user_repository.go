package repository

import (
	"courier-app/app_constant"
	"courier-app/dto"
	"courier-app/entity"
	"fmt"

	"gorm.io/gorm"
)

type UserRepositoryConfig struct {
	DB *gorm.DB
}

type userRepositoryImpl struct {
	db *gorm.DB
}
type UserRepository interface {
	GetDetailsById(id uint) (*entity.UserDetail, error)
	GetById(id uint) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	GetByReferral(referral string) (*entity.User, error)
	GetUserBalance(id uint) (*dto.UserBalanceResponse, error)
	GetUserGachaQuota(id uint) (*dto.UserGachaQuotaResponse, error)
	Register(*entity.User) (*entity.User, error)
	ChangeGachaQuota(userId uint, amount int) error
	UpdateProfile(*dto.UserProfile) error
	RewardReferralClaimer(txreq *entity.Transaction) (*entity.Transaction, error)
	RewardReferralOwner(claimerId uint, referralCode string, txreq *entity.Transaction) (*entity.Transaction, error)
}

func NewUserRepository(u UserRepositoryConfig) UserRepository {
	return &userRepositoryImpl{
		db: u.DB,
	}
}

func (u *userRepositoryImpl) GetById(id uint) (*entity.User, error) {
	var user *entity.User
	err := u.db.Preload("UserDetail").Preload("RoleDetail").Where("id = ?", id).First(&user).Error
	return user, err
}

func (u *userRepositoryImpl) GetByEmail(email string) (*entity.User, error) {
	var user *entity.User
	err := u.db.Preload("RoleDetail").Preload("UserDetail").Where("email = ?", email).First(&user).Error
	return user, err
}

func (u *userRepositoryImpl) Register(user *entity.User) (*entity.User, error) {
	err := u.db.Create(&user).Error
	return user, err
}

func (u *userRepositoryImpl) GetByReferral(referral string) (*entity.User, error) {
	var user *entity.User
	err := u.db.Transaction(func(tx *gorm.DB) error {
		var userDetail *entity.UserDetail
		if err := tx.Model(&entity.UserDetail{}).Where("referral_code= ?", referral).First(&userDetail).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&entity.User{}).Preload("RoleDetail").Preload("UserDetail").Where("id = ?", userDetail.UserID).First(&user).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
	return user, err
}

func (u *userRepositoryImpl) GetDetailsById(id uint) (*entity.UserDetail, error) {
	var detail *entity.UserDetail
	err := u.db.Where("user_id = ?", id).First(&detail).Error
	return detail, err
}

func (u *userRepositoryImpl) GetUserBalance(id uint) (*dto.UserBalanceResponse, error) {
	var res *dto.UserBalanceResponse
	err := u.db.Model(&entity.UserDetail{}).Select("balance, user_id").Where("user_id = ?", id).First(&res).Error
	return res, err
}

func (u *userRepositoryImpl) GetUserGachaQuota(id uint) (*dto.UserGachaQuotaResponse, error) {
	var res *dto.UserGachaQuotaResponse
	err := u.db.Model(&entity.UserDetail{}).Select("gacha_quota, user_id").Where("user_id = ?", id).First(&res).Error
	return res, err
}

func (u *userRepositoryImpl) UpdateProfile(user *dto.UserProfile) error {
	err := u.db.Model(&entity.User{}).Where("id = ?", user.ID).Updates(entity.User{
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Photo: user.Photo,
	}).Error
	return err
}

func (u *userRepositoryImpl) ChangeGachaQuota(userId uint, amount int) error {
	return u.db.Model(&entity.UserDetail{}).Where("user_id = ?", userId).Update("gacha_quota", gorm.Expr("gacha_quota + ?", amount)).Error
}
func (u *userRepositoryImpl) RewardReferralClaimer(txreq *entity.Transaction) (*entity.Transaction, error) {
	fmt.Println("jalan referral")
	err := u.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&txreq).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(&entity.UserDetail{}).Where("user_id = ?", txreq.UserId).Update("balance", gorm.Expr("balance + ?", txreq.Amount)).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(&entity.UserDetail{}).Where("user_id = ?", txreq.UserId).Update("referral_status", app_constant.ReferralStatusClaimedUser).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
	return txreq, err
}

func (u *userRepositoryImpl) RewardReferralOwner(claimerId uint, referralCode string, txreq *entity.Transaction) (*entity.Transaction, error) {
	err := u.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&txreq).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(&entity.UserDetail{}).Where("user_id = ?", txreq.UserId).Update("balance", gorm.Expr("balance + ?", txreq.Amount)).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(&entity.UserDetail{}).Where("user_id = ?", claimerId).Update("referral_status", app_constant.ReferralStatusClaimedFull).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
	return txreq, err
}
