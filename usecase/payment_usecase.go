package usecase

import (
	"courier-app/app_constant"
	"courier-app/dto"
	"courier-app/entity"
	"courier-app/repository"
	"courier-app/usecase/usecase_errors"
	"errors"
	"time"

	"gorm.io/gorm"
)

type PaymentUsecase interface {
	GetPaymentDetail(paymentId uint) (*entity.Payment, error)
	GetPaymentReport(req *dto.PeriodRequest) (*dto.PaymentReport, error)
	Pay(userId uint, payment *dto.PaymentRequest) (*entity.Payment, error)
}

type paymentUsecaseImpl struct {
	repository   repository.PaymentRepository
	promoUsecase PromoUsecase
	userUsecase  UserUsecase
}

type PaymentUsecaseConfig struct {
	PaymentRepo  repository.PaymentRepository
	PromoUsecase PromoUsecase
	UserUsecase  UserUsecase
}

func NewPaymentUsecase(c PaymentUsecaseConfig) PaymentUsecase {
	return &paymentUsecaseImpl{
		repository:   c.PaymentRepo,
		promoUsecase: c.PromoUsecase,
		userUsecase:  c.UserUsecase,
	}
}

func (pa *paymentUsecaseImpl) GetPaymentDetail(paymentId uint) (*entity.Payment, error) {
	payment, err := pa.repository.GetPaymentDetail(paymentId)
	if err == gorm.ErrRecordNotFound {
		return payment, usecase_errors.ErrRecordNotExist
	}

	return payment, err
}

func (pa *paymentUsecaseImpl) Pay(userId uint, payment *dto.PaymentRequest) (*entity.Payment, error) {
	var paymentRecord *entity.Payment
	var totalDiscount float64
	user, err := pa.userUsecase.GetById(userId)
	if err != nil {
		return paymentRecord, err
	}

	paymentDetails, err := pa.repository.GetPaymentDetail(payment.PaymentId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return paymentRecord, usecase_errors.ErrPaymentNotExist
		}
		return paymentRecord, err
	}
	if paymentDetails.Status {
		return paymentRecord, usecase_errors.ErrPaymentAlreadyPaid
	}

	if payment.VoucherId != nil {
		voucher, err := pa.promoUsecase.GetVoucherById(*payment.VoucherId)
		if err != nil {
			if err == usecase_errors.ErrRecordNotExist {
				return paymentRecord, usecase_errors.ErrVoucherInvalid
			}
			return paymentRecord, err
		}
		if voucher.UserId != userId {
			return paymentRecord, usecase_errors.ErrVoucherInvalid
		}
		if time.Now().After(voucher.ExpDate) {
			return paymentRecord, usecase_errors.ErrVoucherExpired
		}
		if paymentDetails.TotalCost < voucher.Promo.MinFee {
			return paymentRecord, usecase_errors.ErrMinFeeNotReached
		}
		totalDiscount = paymentDetails.TotalCost * voucher.Promo.Discount
		if totalDiscount > voucher.Promo.MaxDiscount {
			totalDiscount = voucher.Promo.MaxDiscount
		}
		payment.TotalDiscount = totalDiscount
	}
	payment.Amount = paymentDetails.TotalCost - payment.TotalDiscount
	if payment.Amount > user.Balance {
		return paymentRecord, usecase_errors.ErrInsufficientBalance
	}
	payment.Description = app_constant.PaymentDescription

	return pa.repository.Pay(userId, payment)
}

func (pa *paymentUsecaseImpl) GetPaymentReport(req *dto.PeriodRequest) (*dto.PaymentReport, error) {
	res, err := pa.repository.GetPaymentReport(req)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return res, nil
	}
	return res, err
}
