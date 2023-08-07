package usecase_test

import (
	"courier-app/dto"
	"courier-app/entity"
	mocks "courier-app/mocks/repository"
	"courier-app/usecase"
	"courier-app/usecase/usecase_errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAddressGetById(t *testing.T) {
	t.Run("should return address", func(t *testing.T) {
		addressId := uint(1)
		mockResponse := &entity.Address{
			ID:             uint(addressId),
			RecipientName:  "joko",
			FullAddress:    "bubububu",
			RecipientPhone: "123123123",
			UserId:         1,
			Label:          "rumah",
		}
		mockRepo := &mocks.AddressRepository{}
		mockUc := usecase.NewAddressUsecase(usecase.AddressUsecaseConfig{
			AddressRepo: mockRepo,
		})

		mockRepo.On("GetById", addressId).Return(mockResponse, nil)
		actualRes, err := mockUc.GetById(uint(addressId))

		assert.Nil(t, err)
		assert.Equal(t, actualRes, mockResponse)
	})

	t.Run("should return error when not exist", func(t *testing.T) {
		addressId := uint(1)
		mockRepo := &mocks.AddressRepository{}
		mockUc := usecase.NewAddressUsecase(usecase.AddressUsecaseConfig{
			AddressRepo: mockRepo,
		})

		mockRepo.On("GetById", addressId).Return(nil, gorm.ErrRecordNotFound)
		actualRes, err := mockUc.GetById(uint(addressId))

		assert.EqualError(t, err, usecase_errors.ErrRecordNotExist.Error())
		assert.Nil(t, actualRes)
	})
}

func TestAddressGetAll(t *testing.T) {
	t.Run("should return array of addresses", func(t *testing.T) {
		userId := uint(1)
		expectedPagi := &dto.Pagination{
			PageSize:   1,
			Page:       1,
			TotalCount: 2,
			PageCount:  2,
		}
		mockPagination := &dto.AddressPaginationRequest{
			UserId:   &userId,
			Search:   "123",
			PageSize: 1,
			Page:     1,
		}
		mockResponse := &[]entity.Address{{
			ID:             uint(userId),
			RecipientName:  "joko",
			FullAddress:    "bubububu",
			RecipientPhone: "123123123",
			UserId:         1,
			Label:          "rumah",
		}, {
			ID:             uint(userId),
			RecipientName:  "joko",
			FullAddress:    "bubububu",
			RecipientPhone: "123123123",
			UserId:         1,
			Label:          "rumah",
		}}
		mockRepo := &mocks.AddressRepository{}
		mockUc := usecase.NewAddressUsecase(usecase.AddressUsecaseConfig{
			AddressRepo: mockRepo,
		})

		mockRepo.On("GetAll", mockPagination).Return(mockResponse, nil)
		mockRepo.On("GetCountAll", mockPagination).Return(int64(2), nil)
		addresses, pagination, err := mockUc.GetAll(mockPagination)

		assert.Nil(t, err)
		assert.Equal(t, expectedPagi, pagination)
		assert.Equal(t, mockResponse, addresses)
	})
}
