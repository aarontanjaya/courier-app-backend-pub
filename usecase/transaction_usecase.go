package usecase

import (
	"courier-app/app_constant"
	"courier-app/config"
	"courier-app/dto"
	"courier-app/entity"
	"courier-app/repository"
	"courier-app/usecase/usecase_errors"
)

type TransactionUsecase interface {
	GetTotalSpendingByUser(userId uint) (float64, error)
	TopUp(userId uint, req *dto.TopUpRequestBody) (*entity.Transaction, error)
}

type transactionUsecaseImpl struct {
	repository repository.TransactionRepository
}

type TransactionUsecaseConfig struct {
	TransactionRepo repository.TransactionRepository
}

func NewTransactionUsecase(tr TransactionUsecaseConfig) TransactionUsecase {
	return &transactionUsecaseImpl{
		repository: tr.TransactionRepo,
	}
}

func (tr *transactionUsecaseImpl) GetTotalSpendingByUser(userId uint) (float64, error) {
	total, err := tr.repository.GetTotalSpendingByUser(userId)
	return total, err
}

func (tr *transactionUsecaseImpl) TopUp(userId uint, req *dto.TopUpRequestBody) (*entity.Transaction, error) {
	var res *entity.Transaction
	tx := &entity.Transaction{
		Description: app_constant.TopUpDescription,
		Amount:      *req.Amount,
		UserId:      userId,
	}
	topupConfig := config.Config.TopupConfig
	if *req.Amount > topupConfig.TopupMax || *req.Amount < topupConfig.TopupMin {
		return res, usecase_errors.ErrTopupAmountInvalid
	}
	res, err := tr.repository.CreateTransaction(tx)
	return res, err
}
