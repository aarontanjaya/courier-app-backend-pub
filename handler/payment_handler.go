package handler

import (
	"courier-app/app_constant"
	"courier-app/dto"
	"courier-app/httperror"
	"courier-app/usecase/usecase_errors"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandlePayment(c *gin.Context) {
	var req dto.PaymentRequestBody
	var voucherId *uint
	claim, ok := c.Get("user")
	user, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 || !strings.Contains(user.Scope, "user") {
		c.Error(httperror.UnauthorizedError())
		return
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(httperror.BadRequestError("Bad Request", "BAD_REQUEST"))
		return
	}

	if req.VoucherId != 0 {
		voucherId = &req.VoucherId
	}

	idStr := c.Param("id")
	if idStr == "" {
		c.Error(httperror.NotFoundError("Resource id required"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.Error(httperror.NotFoundError("Not Found"))
		return
	}

	details, err := h.userUsecase.GetDetailsById(user.UserId)
	if err != nil {
		if errors.Is(err, usecase_errors.ErrRecordNotExist) {
			c.Error(httperror.UnauthorizedError())
			return
		}
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}

	res, err := h.paymentUsecase.Pay(details.UserID, &dto.PaymentRequest{
		VoucherId: voucherId,
		PaymentId: uint(id),
	})
	if err != nil {
		if errors.Is(err, usecase_errors.ErrPaymentNotExist) {
			c.Error(httperror.NotFoundError("Payment not found"))
			return
		}
		if errors.Is(err, usecase_errors.ErrPaymentAlreadyPaid) || errors.Is(err, usecase_errors.ErrVoucherInvalid) || errors.Is(err, usecase_errors.ErrVoucherExpired) || errors.Is(err, usecase_errors.ErrMinFeeNotReached) || errors.Is(err, usecase_errors.ErrInsufficientBalance) {
			handlePaymentBadRequestError(err, c)
			return
		}
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}

	referralNotFullyClaimed := !(details.ReferralStatus == uint(app_constant.ReferralStatusClaimedFull))
	if referralNotFullyClaimed {
		_, err = h.userUsecase.HandleReferral(details)
		if err != nil {
			c.Error(httperror.InternalServerError("Internal Server Error"))
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
		"data":        res,
	})
}

func (h *Handler) HandleGetPaymentReport(c *gin.Context) {
	var req *dto.PeriodRequest

	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.Error(httperror.BadRequestError("Bad Request", "BAD_REQUEST"))
		return
	}

	res, err := h.paymentUsecase.GetPaymentReport(req)
	if err != nil {
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
		"data":        res,
	})
}

func handlePaymentBadRequestError(err error, c *gin.Context) {
	if errors.Is(err, usecase_errors.ErrPaymentAlreadyPaid) {
		c.Error(httperror.BadRequestError("Bill already paid", "PAYMENT_PAID"))
		return
	}
	if errors.Is(err, usecase_errors.ErrVoucherInvalid) {
		c.Error(httperror.BadRequestError("Voucher not found", "VOUCHER_INVALID"))
		return
	}
	if errors.Is(err, usecase_errors.ErrVoucherExpired) {
		c.Error(httperror.BadRequestError("Voucher expired", "VOUCHER_EXPIRED"))
		return
	}
	if errors.Is(err, usecase_errors.ErrMinFeeNotReached) {
		c.Error(httperror.BadRequestError("Cost below minimum voucher fee", "VOUCHER_MINFEE"))
		return
	}
	if errors.Is(err, usecase_errors.ErrInsufficientBalance) {
		c.Error(httperror.BadRequestError("Balance insufficient", "BALANCE_INSUFFICIENT"))
		return
	}
}
