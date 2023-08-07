package handler

import (
	"courier-app/usecase"
)

type Handler struct {
	userUsecase        usecase.UserUsecase
	authUsecase        usecase.AuthUsecase
	addressUsecase     usecase.AddressUsecase
	shippingUsecase    usecase.ShippingUsecase
	paymentUsecase     usecase.PaymentUsecase
	transactionUsecase usecase.TransactionUsecase
	addOnUsecase       usecase.AddOnUsecase
	sizeUsecase        usecase.SizeUsecase
	categoryUsecase    usecase.CategoryUsecase
	promoUsecase       usecase.PromoUsecase
}

type HandlerConfig struct {
	UserUsecase        usecase.UserUsecase
	AuthUsecase        usecase.AuthUsecase
	AddressUsecase     usecase.AddressUsecase
	ShippingUsecase    usecase.ShippingUsecase
	PaymentUsecase     usecase.PaymentUsecase
	TransactionUsecase usecase.TransactionUsecase
	AddOnUsecase       usecase.AddOnUsecase
	SizeUsecase        usecase.SizeUsecase
	CategoryUsecase    usecase.CategoryUsecase
	PromoUsecase       usecase.PromoUsecase
}

func New(c HandlerConfig) *Handler {
	return &Handler{
		userUsecase:        c.UserUsecase,
		authUsecase:        c.AuthUsecase,
		addressUsecase:     c.AddressUsecase,
		shippingUsecase:    c.ShippingUsecase,
		paymentUsecase:     c.PaymentUsecase,
		transactionUsecase: c.TransactionUsecase,
		addOnUsecase:       c.AddOnUsecase,
		categoryUsecase:    c.CategoryUsecase,
		sizeUsecase:        c.SizeUsecase,
		promoUsecase:       c.PromoUsecase,
	}
}
