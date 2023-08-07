package server

import (
	"courier-app/handler"
	"courier-app/middleware"
	"courier-app/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouterConfig struct {
	UserUsecase        usecase.UserUsecase
	AuthUsecase        usecase.AuthUsecase
	AddressUsecase     usecase.AddressUsecase
	ShippingUsecase    usecase.ShippingUsecase
	PaymentUsecase     usecase.PaymentUsecase
	TransactionUsecase usecase.TransactionUsecase
	AddOnUsecase       usecase.AddOnUsecase
	CategoryUsecase    usecase.CategoryUsecase
	SizeUsecase        usecase.SizeUsecase
	PromoUsecase       usecase.PromoUsecase
}

func NewRouter(c RouterConfig) *gin.Engine {
	r := gin.Default()
	AppHandler := handler.New(handler.HandlerConfig{
		UserUsecase:        c.UserUsecase,
		AuthUsecase:        c.AuthUsecase,
		AddressUsecase:     c.AddressUsecase,
		ShippingUsecase:    c.ShippingUsecase,
		PaymentUsecase:     c.PaymentUsecase,
		TransactionUsecase: c.TransactionUsecase,
		AddOnUsecase:       c.AddOnUsecase,
		SizeUsecase:        c.SizeUsecase,
		CategoryUsecase:    c.CategoryUsecase,
		PromoUsecase:       c.PromoUsecase,
	})
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.ErrorHandler)
	v1 := r.Group("/v1")

	v1.POST("/login", AppHandler.HandleLogin)
	v1.POST("/register", AppHandler.HandleRegister)
	v1.POST("/logout", AppHandler.HandleLogout)
	v1.Use(middleware.SetHeaderFromCookie)
	v1.Use(middleware.Authorize)
	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	v1.GET("/whoami", AppHandler.HandleWhoAmI)
	v1.GET("/profile", AppHandler.HandleGetProfile)
	v1.POST("/profile", AppHandler.HandleUpdateProfile)
	v1.GET("/categories", AppHandler.HandleGetAllCategories)
	v1.GET("/add-ons", AppHandler.HandleGetAllAddOns)
	v1.GET("/sizes", AppHandler.HandleGetAllSizes)
	v1.GET("/shipping-statuses", AppHandler.HandleGetShippingStatuses)

	ug := v1.Group("/user")
	ug.Use(middleware.AuthorizeUser)
	ug.GET("/addresses", AppHandler.HandleGetUserAddresses)
	ug.PUT("/addresses/:id", AppHandler.HandleUpdateAddress)
	ug.POST("/addresses", AppHandler.HandleCreateAddress)
	ug.DELETE("/addresses/:id", AppHandler.HandleDeleteAdress)
	ug.GET("/shippings/:id", AppHandler.HandleGetShippingById)
	ug.GET("/shippings", AppHandler.HandleGetShippings)
	ug.POST("/shippings/:id/review", AppHandler.HandleReviewShipping)
	ug.POST("/shippings", AppHandler.HandleCreateShipping)
	ug.POST("/topup", AppHandler.HandleTopUp)
	ug.POST("/pay/:id", AppHandler.HandlePayment)
	ug.GET("/gacha", AppHandler.HandleGetGachaQuota)
	ug.POST("/gacha", AppHandler.HandlePlayGacha)
	ug.GET("/vouchers", AppHandler.HandleGetAllActiveUserVouchers)
	ug.GET("/balance", AppHandler.HandleGetBalance)

	ag := v1.Group("/admin")
	ag.Use(middleware.AuthorizeAdmin)
	ag.GET("/addresses", AppHandler.HandleGetAdminAddresses)
	ag.PATCH("/shippings/:id/status", AppHandler.HandleUpdateShippingStatus)
	ag.GET("/shippings", AppHandler.HandleGetAllShippings)
	ag.GET("/payment/report", AppHandler.HandleGetPaymentReport)
	ag.GET("/promos", AppHandler.HandleGetAllPromo)
	ag.POST("/promos", AppHandler.HandleCreatePromo)
	ag.PUT("/promos/:id", AppHandler.HandleUpdatePromo)

	return r
}
