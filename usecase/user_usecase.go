package usecase

import (
	"bytes"
	"courier-app/app_constant"
	"courier-app/config"
	"courier-app/dto"
	"courier-app/entity"
	"courier-app/repository"
	"courier-app/usecase/usecase_errors"
	"encoding/base64"
	"fmt"
	"io"

	"gorm.io/gorm"
)

type UserUsecase interface {
	GetDetailsById(id uint) (*entity.UserDetail, error)
	GetById(id uint) (*dto.UserResponse, error)
	GetByEmail(email string) (*dto.UserResponse, error)
	GetByReferral(referral string) (*dto.UserResponse, error)
	GetUserBalance(id uint) (*dto.UserBalanceResponse, error)
	GetUserGachaQuota(id uint) (*dto.UserGachaQuotaResponse, error)
	Register(*entity.User) (*dto.RegisterResponse, error)
	ChangeGachaQuota(userId uint, amount int) error
	UpdateProfile(user *dto.UpdateProfileRequest, id uint) error
	HandleReferral(user *entity.UserDetail) (*[]entity.Transaction, error)
}

type userUsecaseImpl struct {
	repository         repository.UserRepository
	transactionUsecase TransactionUsecase
}

type UserUsecaseConfig struct {
	UserRepo           repository.UserRepository
	TransactionUsecase TransactionUsecase
}

func NewUserUsecase(c UserUsecaseConfig) UserUsecase {
	return &userUsecaseImpl{
		repository:         c.UserRepo,
		transactionUsecase: c.TransactionUsecase,
	}
}

func (u *userUsecaseImpl) GetById(id uint) (*dto.UserResponse, error) {
	var photo string
	user, err := u.repository.GetById(id)
	if len(user.Photo) > 0 {
		str := base64.StdEncoding.EncodeToString(user.Photo)
		photo = fmt.Sprintf("data:%s;base64, %s", user.PhotoFormat, str)
	}
	userDto := &dto.UserResponse{
		Name:         user.Name,
		Email:        user.Email,
		Password:     user.Password,
		ID:           user.ID,
		Phone:        user.Phone,
		Role:         user.Role,
		RoleName:     user.RoleDetail.Name,
		Balance:      user.UserDetail.Balance,
		ReferralCode: user.UserDetail.ReferralCode,
		Photo:        photo,
		GachaQuota:   user.UserDetail.GachaQuota,
	}
	if err == gorm.ErrRecordNotFound {
		return userDto, usecase_errors.ErrRecordNotExist
	}
	return userDto, err
}

func (u *userUsecaseImpl) GetByEmail(email string) (*dto.UserResponse, error) {
	var userDto *dto.UserResponse
	var photo string
	user, err := u.repository.GetByEmail(email)
	if err != nil {
		return &dto.UserResponse{}, err
	}

	if len(user.Photo) > 0 {
		str := base64.StdEncoding.EncodeToString(user.Photo)
		photo = fmt.Sprintf("data:%s;base64, %s", user.PhotoFormat, str)
	}

	userDto = &dto.UserResponse{
		Name:         user.Name,
		Email:        user.Email,
		Password:     user.Password,
		ID:           user.ID,
		Phone:        user.Phone,
		Role:         user.Role,
		RoleName:     user.RoleDetail.Name,
		Balance:      user.UserDetail.Balance,
		ReferralCode: user.UserDetail.ReferralCode,
		Photo:        photo,
		GachaQuota:   user.UserDetail.GachaQuota,
	}

	if err == gorm.ErrRecordNotFound {
		return userDto, usecase_errors.ErrRecordNotExist
	}
	return userDto, err
}

func (u *userUsecaseImpl) GetByReferral(referral string) (*dto.UserResponse, error) {
	var userDto *dto.UserResponse
	var photo string
	user, err := u.repository.GetByReferral(referral)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &dto.UserResponse{}, usecase_errors.ErrReferralInvalid
		}
		return &dto.UserResponse{}, err
	}

	if len(user.Photo) > 0 {
		str := base64.StdEncoding.EncodeToString(user.Photo)
		photo = fmt.Sprintf("data:%s;base64, %s", user.PhotoFormat, str)
	}

	userDto = &dto.UserResponse{
		Name:         user.Name,
		Email:        user.Email,
		Password:     user.Password,
		ID:           user.ID,
		Phone:        user.Phone,
		Role:         user.Role,
		RoleName:     user.RoleDetail.Name,
		Balance:      user.UserDetail.Balance,
		ReferralCode: user.UserDetail.ReferralCode,
		Photo:        photo,
		GachaQuota:   user.UserDetail.GachaQuota,
	}

	return userDto, err
}

func (u *userUsecaseImpl) GetUserBalance(id uint) (*dto.UserBalanceResponse, error) {
	res, err := u.repository.GetUserBalance(id)
	if err == gorm.ErrRecordNotFound {
		return res, usecase_errors.ErrRecordNotExist
	}
	return res, err
}

func (u *userUsecaseImpl) GetUserGachaQuota(id uint) (*dto.UserGachaQuotaResponse, error) {
	res, err := u.repository.GetUserGachaQuota(id)
	if err == gorm.ErrRecordNotFound {
		return res, usecase_errors.ErrRecordNotExist
	}
	return res, err
}

func (u *userUsecaseImpl) Register(user *entity.User) (*dto.RegisterResponse, error) {
	var res *dto.RegisterResponse
	_, err := u.GetByEmail(user.Email)
	if err == nil {
		return res, usecase_errors.ErrEmailAlreadyExist
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return res, err
	}

	user, err = u.repository.Register(user)
	if err != nil {
		return res, err
	}

	res = &dto.RegisterResponse{
		Name:  user.Name,
		Email: user.Email,
		Id:    user.ID,
	}
	return res, err
}

func (u *userUsecaseImpl) UpdateProfile(user *dto.UpdateProfileRequest, id uint) error {
	var photo []byte
	existingUser, err := u.GetByEmail(user.Email)
	if err == nil && existingUser.ID != id {
		return usecase_errors.ErrEmailAlreadyExist
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if user.Photo != nil {
		photoFile, err := user.Photo.Open()
		defer func() {
			err = photoFile.Close()
		}()
		if err != nil {
			return err
		}
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, photoFile); err != nil {
			return err
		}
		photo = buf.Bytes()
	}

	err = u.repository.UpdateProfile(&dto.UserProfile{
		ID:    id,
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Photo: photo,
	})
	return err
}

func (u *userUsecaseImpl) ChangeGachaQuota(userId uint, amount int) error {
	_, err := u.GetDetailsById(userId)
	if err != nil {
		return err
	}

	return u.repository.ChangeGachaQuota(userId, amount)
}

func (u *userUsecaseImpl) GetDetailsById(id uint) (*entity.UserDetail, error) {
	detail, err := u.repository.GetDetailsById(id)
	if err == gorm.ErrRecordNotFound {
		return detail, usecase_errors.ErrRecordNotExist
	}
	return detail, err
}

func (u *userUsecaseImpl) HandleReferral(user *entity.UserDetail) (*[]entity.Transaction, error) {
	var ts []entity.Transaction
	referralConfig := config.Config.ReferralConfig
	if user.ReferralStatus == uint(app_constant.ReferralStatusNone) || user.ReferralStatus == uint(app_constant.ReferralStatusClaimedFull) {
		return &ts, nil
	}

	total, err := u.transactionUsecase.GetTotalSpendingByUser(user.UserID)
	if err != nil {
		return &ts, err
	}
	if total >= referralConfig.ClaimerTreshold && user.ReferralStatus == uint(app_constant.ReferralStatusUnclaimed) {
		tx1 := &entity.Transaction{
			Amount:      referralConfig.ClaimerReward,
			UserId:      user.UserID,
			Description: app_constant.ReferralDescription,
		}
		tx1, err := u.repository.RewardReferralClaimer(tx1)

		if err != nil {
			return &ts, err
		}
		ts = append(ts, *tx1)
	}

	if total >= referralConfig.OwnerTreshold && (user.ReferralStatus == uint(app_constant.ReferralStatusUnclaimed) || user.ReferralStatus == uint(app_constant.ReferralStatusClaimedUser)) {
		owner, err := u.GetByReferral(user.ClaimedReferral)
		if err != nil {
			return &ts, err
		}
		tx2 := &entity.Transaction{
			Amount:      referralConfig.OwnerReward,
			UserId:      owner.ID,
			Description: "Referral Reward",
		}
		tx2, err = u.repository.RewardReferralOwner(user.UserID, user.ClaimedReferral, tx2)
		if err != nil {
			return &ts, err
		}
		ts = append(ts, *tx2)
	}

	return &ts, err
}
