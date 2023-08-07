package usecase_test

import (
	"courier-app/dto"
	"courier-app/entity"
	mocks "courier-app/mocks/repository"
	"courier-app/usecase"
	"courier-app/usecase/usecase_errors"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetUserBalance(t *testing.T) {
	t.Run("should return user balance", func(t *testing.T) {
		mockRepositoryResponse := &dto.UserBalanceResponse{
			UserId:  1,
			Balance: 10000,
		}
		expectedResponse := &dto.UserBalanceResponse{
			UserId:  1,
			Balance: 10000,
		}
		userRepo := &mocks.UserRepository{}
		userUc := usecase.NewUserUsecase(usecase.UserUsecaseConfig{
			UserRepo: userRepo,
		})

		userRepo.On("GetUserBalance", uint(1)).Return(mockRepositoryResponse, nil)
		response, err := userUc.GetUserBalance(1)

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, response)
	})

	t.Run("should return error record not exist if userid invalid", func(t *testing.T) {
		userRepo := &mocks.UserRepository{}
		userUc := usecase.NewUserUsecase(usecase.UserUsecaseConfig{
			UserRepo: userRepo,
		})

		userRepo.On("GetUserBalance", uint(1)).Return(nil, gorm.ErrRecordNotFound)
		response, err := userUc.GetUserBalance(1)

		assert.EqualError(t, err, usecase_errors.ErrRecordNotExist.Error())
		assert.Nil(t, response)
	})
}

func TestGetUserGachaQuota(t *testing.T) {
	t.Run("should return gacha quota", func(t *testing.T) {
		mockRepoResponse := &dto.UserGachaQuotaResponse{
			GachaQuota: 10,
			UserId:     1,
		}
		expectedResponse := &dto.UserGachaQuotaResponse{
			GachaQuota: 10,
			UserId:     1,
		}
		userRepo := &mocks.UserRepository{}
		userUc := usecase.NewUserUsecase(usecase.UserUsecaseConfig{
			UserRepo: userRepo,
		})

		userRepo.On("GetUserGachaQuota", uint(1)).Return(mockRepoResponse, nil)

		actualRes, err := userUc.GetUserGachaQuota(1)

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, actualRes)
	})

	t.Run("should return error record not exist if not found", func(t *testing.T) {
		userRepo := &mocks.UserRepository{}
		userUc := usecase.NewUserUsecase(usecase.UserUsecaseConfig{
			UserRepo: userRepo,
		})

		userRepo.On("GetUserGachaQuota", uint(1)).Return(nil, gorm.ErrRecordNotFound)

		actualRes, err := userUc.GetUserGachaQuota(1)

		assert.Nil(t, actualRes)
		assert.EqualError(t, err, usecase_errors.ErrRecordNotExist.Error())
	})
}

func TestGetById(t *testing.T) {
	t.Run("should return user", func(t *testing.T) {
		mockRoleDetail := &entity.Role{
			Name: "user",
		}
		mockUserDetail := &entity.UserDetail{
			Balance:      123,
			ReferralCode: "123123",
			GachaQuota:   123,
		}
		mockRepoResponse := &entity.User{
			Name:       "tes",
			Email:      "tes",
			Password:   "tes",
			ID:         1,
			Phone:      "1",
			Role:       1,
			UserDetail: *mockUserDetail,
			RoleDetail: *mockRoleDetail,
		}
		expectedResponse := &dto.UserResponse{
			Name:         "tes",
			Email:        "tes",
			Password:     "tes",
			ID:           1,
			Phone:        "1",
			Role:         1,
			RoleName:     "user",
			Balance:      123,
			ReferralCode: "123123",
			Photo:        "",
			GachaQuota:   123,
		}
		userRepo := &mocks.UserRepository{}
		userUc := usecase.NewUserUsecase(usecase.UserUsecaseConfig{
			UserRepo: userRepo,
		})

		userRepo.On("GetById", uint(1)).Return(mockRepoResponse, nil)
		actualRes, err := userUc.GetById(1)

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, actualRes)
	})

	t.Run("should return user with photo", func(t *testing.T) {
		mockPhoto := []byte{'1', '2', '3'}
		mockResPhoto := base64.StdEncoding.EncodeToString(mockPhoto)
		mockFinPhoto := fmt.Sprintf("data:%s;base64, %s", "img/png", mockResPhoto)
		mockRoleDetail := &entity.Role{
			Name: "user",
		}
		mockUserDetail := &entity.UserDetail{
			Balance:      123,
			ReferralCode: "123123",
			GachaQuota:   123,
		}
		mockRepoResponse := &entity.User{
			Name:        "tes",
			Email:       "tes",
			Password:    "tes",
			ID:          1,
			Photo:       mockPhoto,
			Phone:       "1",
			Role:        1,
			UserDetail:  *mockUserDetail,
			RoleDetail:  *mockRoleDetail,
			PhotoFormat: "img/png",
		}
		expectedResponse := &dto.UserResponse{
			Name:         "tes",
			Email:        "tes",
			Password:     "tes",
			ID:           1,
			Phone:        "1",
			Photo:        mockFinPhoto,
			Role:         1,
			RoleName:     "user",
			Balance:      123,
			ReferralCode: "123123",
			GachaQuota:   123,
		}
		userRepo := &mocks.UserRepository{}
		userUc := usecase.NewUserUsecase(usecase.UserUsecaseConfig{
			UserRepo: userRepo,
		})

		userRepo.On("GetById", uint(1)).Return(mockRepoResponse, nil)
		actualRes, err := userUc.GetById(1)

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, actualRes)
	})
}

func TestGetByEmail(t *testing.T) {
	t.Run("should return user", func(t *testing.T) {
		mockRoleDetail := &entity.Role{
			Name: "user",
		}
		mockUserDetail := &entity.UserDetail{
			Balance:      123,
			ReferralCode: "123123",
			GachaQuota:   123,
		}
		mockRepoResponse := &entity.User{
			Name:       "tes",
			Email:      "tes",
			Password:   "tes",
			ID:         1,
			Phone:      "1",
			Role:       1,
			UserDetail: *mockUserDetail,
			RoleDetail: *mockRoleDetail,
		}
		expectedResponse := &dto.UserResponse{
			Name:         "tes",
			Email:        "tes",
			Password:     "tes",
			ID:           1,
			Phone:        "1",
			Role:         1,
			RoleName:     "user",
			Balance:      123,
			ReferralCode: "123123",
			Photo:        "",
			GachaQuota:   123,
		}
		userRepo := &mocks.UserRepository{}
		userUc := usecase.NewUserUsecase(usecase.UserUsecaseConfig{
			UserRepo: userRepo,
		})

		userRepo.On("GetByEmail", "tes").Return(mockRepoResponse, nil)
		actualRes, err := userUc.GetByEmail("tes")

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, actualRes)
	})

	t.Run("should return user with photo", func(t *testing.T) {
		mockPhoto := []byte{'1', '2', '3'}
		mockResPhoto := base64.StdEncoding.EncodeToString(mockPhoto)
		mockFinPhoto := fmt.Sprintf("data:%s;base64, %s", "img/png", mockResPhoto)
		mockRoleDetail := &entity.Role{
			Name: "user",
		}
		mockUserDetail := &entity.UserDetail{
			Balance:      123,
			ReferralCode: "123123",
			GachaQuota:   123,
		}
		mockRepoResponse := &entity.User{
			Name:        "tes",
			Email:       "tes",
			Password:    "tes",
			ID:          1,
			Photo:       mockPhoto,
			Phone:       "1",
			Role:        1,
			UserDetail:  *mockUserDetail,
			RoleDetail:  *mockRoleDetail,
			PhotoFormat: "img/png",
		}
		expectedResponse := &dto.UserResponse{
			Name:         "tes",
			Email:        "tes",
			Password:     "tes",
			ID:           1,
			Phone:        "1",
			Photo:        mockFinPhoto,
			Role:         1,
			RoleName:     "user",
			Balance:      123,
			ReferralCode: "123123",
			GachaQuota:   123,
		}
		userRepo := &mocks.UserRepository{}
		userUc := usecase.NewUserUsecase(usecase.UserUsecaseConfig{
			UserRepo: userRepo,
		})

		userRepo.On("GetByEmail", "tes").Return(mockRepoResponse, nil)
		actualRes, err := userUc.GetByEmail("tes")

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, actualRes)
	})
}

func TestGetByReferral(t *testing.T) {
	t.Run("should return user", func(t *testing.T) {
		mockRoleDetail := &entity.Role{
			Name: "user",
		}
		mockUserDetail := &entity.UserDetail{
			Balance:      123,
			ReferralCode: "123123",
			GachaQuota:   123,
		}
		mockRepoResponse := &entity.User{
			Name:       "tes",
			Email:      "tes",
			Password:   "tes",
			ID:         1,
			Phone:      "1",
			Role:       1,
			UserDetail: *mockUserDetail,
			RoleDetail: *mockRoleDetail,
		}
		expectedResponse := &dto.UserResponse{
			Name:         "tes",
			Email:        "tes",
			Password:     "tes",
			ID:           1,
			Phone:        "1",
			Role:         1,
			RoleName:     "user",
			Balance:      123,
			ReferralCode: "123123",
			Photo:        "",
			GachaQuota:   123,
		}
		userRepo := &mocks.UserRepository{}
		userUc := usecase.NewUserUsecase(usecase.UserUsecaseConfig{
			UserRepo: userRepo,
		})

		userRepo.On("GetByReferral", "123123").Return(mockRepoResponse, nil)
		actualRes, err := userUc.GetByReferral("123123")

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, actualRes)
	})

	t.Run("should return user with photo", func(t *testing.T) {
		mockPhoto := []byte{'1', '2', '3'}
		mockResPhoto := base64.StdEncoding.EncodeToString(mockPhoto)
		mockFinPhoto := fmt.Sprintf("data:%s;base64, %s", "img/png", mockResPhoto)
		mockRoleDetail := &entity.Role{
			Name: "user",
		}
		mockUserDetail := &entity.UserDetail{
			Balance:      123,
			ReferralCode: "123123",
			GachaQuota:   123,
		}
		mockRepoResponse := &entity.User{
			Name:        "tes",
			Email:       "tes",
			Password:    "tes",
			ID:          1,
			Photo:       mockPhoto,
			Phone:       "1",
			Role:        1,
			UserDetail:  *mockUserDetail,
			RoleDetail:  *mockRoleDetail,
			PhotoFormat: "img/png",
		}
		expectedResponse := &dto.UserResponse{
			Name:         "tes",
			Email:        "tes",
			Password:     "tes",
			ID:           1,
			Phone:        "1",
			Photo:        mockFinPhoto,
			Role:         1,
			RoleName:     "user",
			Balance:      123,
			ReferralCode: "123123",
			GachaQuota:   123,
		}
		userRepo := &mocks.UserRepository{}
		userUc := usecase.NewUserUsecase(usecase.UserUsecaseConfig{
			UserRepo: userRepo,
		})

		userRepo.On("GetByReferral", "123123").Return(mockRepoResponse, nil)
		actualRes, err := userUc.GetByReferral("123123")

		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, actualRes)
	})
}
