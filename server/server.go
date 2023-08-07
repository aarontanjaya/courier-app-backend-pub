package server

import (
	"courier-app/db"
	"courier-app/repository"
	"courier-app/usecase"
	"fmt"

	"github.com/gin-gonic/gin"
)

func initRouter() *gin.Engine {
	transactionRepo := repository.NewTransactionRepository(repository.TransactionRepositoryConfig{
		DB: db.Get(),
	})
	transactionUsecase := usecase.NewTransactionUsecase(usecase.TransactionUsecaseConfig{
		TransactionRepo: transactionRepo,
	})
	userRepo := repository.NewUserRepository(repository.UserRepositoryConfig{
		DB: db.Get(),
	})
	userUsecase := usecase.NewUserUsecase(usecase.UserUsecaseConfig{
		UserRepo:           userRepo,
		TransactionUsecase: transactionUsecase,
	})
	authUsecase := usecase.NewAuthUsecase(usecase.AuthUsecaseConfig{
		UserUsecase: userUsecase,
	})
	addressRepo := repository.NewAddressRepository(repository.AddressRepositoryConfig{
		DB: db.Get(),
	})
	addressUsecase := usecase.NewAddressUsecase(usecase.AddressUsecaseConfig{
		AddressRepo: addressRepo,
	})
	categoryRepo := repository.NewCategoryRepository(repository.CategoryRepositoryConfig{
		DB: db.Get(),
	})
	categoryUsecase := usecase.NewCategoryUsecase(usecase.CategoryUsecaseConfig{
		CategoryRepo: categoryRepo,
	})
	sizeRepo := repository.NewSizeRepository(repository.SizeRepositoryConfig{
		DB: db.Get(),
	})
	sizeUsecase := usecase.NewSizeUsecase(usecase.SizeUsecaseConfig{
		SizeRepo: sizeRepo,
	})
	addOnRepo := repository.NewAddOnRepository(repository.AddOnRepositoryConfig{
		DB: db.Get(),
	})
	addOnUsecase := usecase.NewAddOnUsecase(usecase.AddOnUsecaseConfig{
		AddOnRepo: addOnRepo,
	})
	shippingRepo := repository.NewShippingRepository(repository.ShippingRepositoryConfig{
		DB: db.Get(),
	})
	shippingUsecase := usecase.NewShippingUsecase(usecase.ShippingUsecaseConfig{
		ShippingRepo:    shippingRepo,
		CategoryUsecase: categoryUsecase,
		SizeUsecase:     sizeUsecase,
		AddOnUsecase:    addOnUsecase,
		AddressUsease:   addressUsecase,
		UserUsecase:     userUsecase,
	})
	promoRepo := repository.NewPromoRepository(repository.PromoRepositoryConfig{
		DB: db.Get(),
	})
	promoUsecase := usecase.NewPromoUsecase(usecase.PromoUsecaseConfig{
		PromoRepo:   promoRepo,
		UserUsecase: userUsecase,
	})
	paymentRepo := repository.NewPaymentRepository(repository.PaymentRepositoryConfig{
		DB: db.Get(),
	})
	paymentUsecase := usecase.NewPaymentUsecase(usecase.PaymentUsecaseConfig{
		PaymentRepo:  paymentRepo,
		PromoUsecase: promoUsecase,
		UserUsecase:  userUsecase,
	})
	r := NewRouter(RouterConfig{
		UserUsecase:        userUsecase,
		AuthUsecase:        authUsecase,
		AddressUsecase:     addressUsecase,
		ShippingUsecase:    shippingUsecase,
		PaymentUsecase:     paymentUsecase,
		TransactionUsecase: transactionUsecase,
		AddOnUsecase:       addOnUsecase,
		SizeUsecase:        sizeUsecase,
		CategoryUsecase:    categoryUsecase,
		PromoUsecase:       promoUsecase,
	})
	fmt.Println("available routes", r.Routes())
	return r
}

func Init() {
	r := initRouter()
	err := r.Run()
	if err != nil {
		fmt.Println("error while running server", err)
		return
	}
}
