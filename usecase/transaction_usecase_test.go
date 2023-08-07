package usecase_test

import (
	mocks "courier-app/mocks/repository"
	"courier-app/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTotalSpendingByUser(t *testing.T) {
	t.Run("should return total spending", func(t *testing.T) {
		mockResponse := float64(10000)
		expectedResponse := float64(10000)
		mockRepo := &mocks.TransactionRepository{}
		mockUc := usecase.NewTransactionUsecase(usecase.TransactionUsecaseConfig{
			TransactionRepo: mockRepo,
		})

		mockRepo.On("GetTotalSpendingByUser", uint(1)).Return(mockResponse, nil)
		actualRes, err := mockUc.GetTotalSpendingByUser(1)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, actualRes)
	})
}
