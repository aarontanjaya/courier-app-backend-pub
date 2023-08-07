package repository

import (
	"courier-app/entity"

	"gorm.io/gorm"
)

type TransactionRepositoryConfig struct {
	DB *gorm.DB
}

type transactionRepositoryImpl struct {
	db *gorm.DB
}

type TransactionRepository interface {
	GetTotalSpendingByUser(userId uint) (float64, error)
	CreateTransaction(txreq *entity.Transaction) (*entity.Transaction, error)
}

func NewTransactionRepository(s TransactionRepositoryConfig) TransactionRepository {
	return &transactionRepositoryImpl{
		db: s.DB,
	}
}

func (tr *transactionRepositoryImpl) GetTotalSpendingByUser(userId uint) (float64, error) {
	var total float64
	err := tr.db.Model(&entity.Transaction{}).Select("sum(payments.total_cost) as total").Joins("left join payments on payments.id = transactions.payment_id").Where("transactions.user_id = ? AND status = true", userId).Scan(&total).Error
	return total, err
}

func (tr *transactionRepositoryImpl) CreateTransaction(txreq *entity.Transaction) (*entity.Transaction, error) {
	err := tr.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&txreq).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(&entity.UserDetail{}).Where("user_id = ?", txreq.UserId).Update("balance", gorm.Expr("balance + ?", txreq.Amount)).Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	})
	return txreq, err
}
