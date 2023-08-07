package repository

import (
	"courier-app/app_constant"
	"courier-app/dto"
	"courier-app/entity"

	"gorm.io/gorm"
)

type PaymentRepositoryConfig struct {
	DB *gorm.DB
}

type paymentRepositoryImpl struct {
	db *gorm.DB
}

type PaymentRepository interface {
	GetPaymentDetail(paymentId uint) (*entity.Payment, error)
	GetPaymentReport(req *dto.PeriodRequest) (*dto.PaymentReport, error)
	Pay(userId uint, payment *dto.PaymentRequest) (*entity.Payment, error)
}

func NewPaymentRepository(pa PaymentRepositoryConfig) PaymentRepository {
	return &paymentRepositoryImpl{
		db: pa.DB,
	}
}

func (pa *paymentRepositoryImpl) Pay(userId uint, payment *dto.PaymentRequest) (*entity.Payment, error) {
	transaction := &entity.Transaction{
		Description: payment.Description,
		Amount:      -payment.Amount,
		UserId:      userId,
		PaymentId:   &payment.PaymentId,
	}

	paymentRec := &entity.Payment{
		ID:                payment.PaymentId,
		VoucherId:         payment.VoucherId,
		TotalDiscount:     payment.TotalDiscount,
		TransactionRecord: *transaction,
		Status:            true,
	}

	err := pa.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Updates(paymentRec).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(&entity.UserDetail{}).Where("user_id = ?", userId).Update("balance", gorm.Expr("balance - ?", payment.Amount)).Error; err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(&entity.Shipping{}).Where("payment_id = ?", payment.PaymentId).Update("status_id", app_constant.StatusShippingProcessing).Error; err != nil {
			tx.Rollback()
			return err
		}

		if payment.VoucherId != nil {
			if err := tx.Delete(&entity.UserVoucher{
				ID: *payment.VoucherId,
			}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
		return nil
	})

	return paymentRec, err
}

func (pa *paymentRepositoryImpl) GetPaymentDetail(paymentId uint) (*entity.Payment, error) {
	var payment *entity.Payment
	err := pa.db.Where("id = ?", paymentId).First(&payment).Error
	return payment, err
}

func (pa *paymentRepositoryImpl) GetPaymentReport(req *dto.PeriodRequest) (*dto.PaymentReport, error) {
	query := pa.db.Model(&entity.Payment{}).Select("sum(total_cost) as revenue, sum(total_discount) as total_discount, Count(id) as count, AVG(total_cost) as avg_size").Where("status = ?", true)
	if !req.StartDate.IsZero() {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if !req.EndDate.IsZero() {
		query = query.Where("created_at <= ?", req.EndDate)
	}
	var res *dto.PaymentReport
	err := query.Scan(&res).Error
	return res, err
}
