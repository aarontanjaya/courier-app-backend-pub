package handler

import (
	"courier-app/dto"
	"courier-app/httperror"
	"courier-app/usecase/usecase_errors"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleGetAllActiveUserVouchers(c *gin.Context) {
	claim, ok := c.Get("user")
	user, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
		return
	}

	vouchers, err := h.promoUsecase.GetAllActiveUserVouchers(user.UserId)
	if err != nil && !errors.Is(err, usecase_errors.ErrRecordNotExist) {
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
		"data":        vouchers,
	})
}

func (h *Handler) HandleCreatePromo(c *gin.Context) {
	var req *dto.PromoRequest
	claim, ok := c.Get("user")
	_, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
		return
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(httperror.BadRequestError("Bad Request", "BAD_REQUEST"))
		return
	}

	promo, err := h.promoUsecase.CreatePromo(req)
	if err != nil {
		if errors.Is(err, usecase_errors.ErrPromoDiscountNegative) || errors.Is(err, usecase_errors.ErrPromoDiscountMax) || errors.Is(err, usecase_errors.ErrPromoMinFeeNegative) || errors.Is(err, usecase_errors.ErrMaxDiscountNegative) || errors.Is(err, usecase_errors.ErrQuotaNegative) {
			c.Error(httperror.BadRequestError(err.Error(), "BAD_REQUEST"))
			return
		}
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
		"data":        promo,
	})
}

func (h *Handler) HandleUpdatePromo(c *gin.Context) {
	var req *dto.PromoRequest
	claim, ok := c.Get("user")
	_, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
		return
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(httperror.BadRequestError("Bad Request", "BAD_REQUEST"))
		return
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

	err = h.promoUsecase.UpdatePromo(uint(id), req)

	if err != nil {
		if errors.Is(err, usecase_errors.ErrRecordNotExist) {
			c.Error(httperror.NotFoundError("Not Found"))
			return
		}
		if errors.Is(err, usecase_errors.ErrPromoDiscountNegative) || errors.Is(err, usecase_errors.ErrPromoDiscountMax) || errors.Is(err, usecase_errors.ErrPromoMinFeeNegative) || errors.Is(err, usecase_errors.ErrMaxDiscountNegative) || errors.Is(err, usecase_errors.ErrQuotaNegative) {
			c.Error(httperror.BadRequestError(err.Error(), "BAD_REQUEST"))
			return
		}
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
	})
}

func (h *Handler) HandleGetAllPromo(c *gin.Context) {
	claim, ok := c.Get("user")
	_, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.JSON(http.StatusUnauthorized, httperror.UnauthorizedError())
		return
	}
	search := c.DefaultQuery("search", "")
	sortField := c.DefaultQuery("sort_field", "exp_date")
	sortDirection := c.DefaultQuery("sort_direction", "desc")
	pageSize, err := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 64)
	if err != nil {
		c.Error(httperror.BadRequestError("Bad Request, page_size invalid", "BAD_REQUEST"))
		return
	}

	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		c.Error(httperror.BadRequestError("Bad Request, page invalid", "BAD_REQUEST"))
		return
	}

	paging := dto.PaginationSortRequest{
		PageSize:      pageSize,
		Search:        search,
		SortDirection: sortDirection,
		SortField:     sortField,
		Page:          page,
	}
	promos, pagination, err := h.promoUsecase.GetAllPromo(&paging)
	if err != nil {
		if errors.Is(err, usecase_errors.ErrPaginationInvalid) {
			c.Error(httperror.BadRequestError("Bad Request, pagination invalid", "BAD_REQUEST"))
			return
		}
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
		"data": dto.PromoResponsePayload{
			Records:    promos,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			PageCount:  pagination.PageCount,
			TotalCount: pagination.TotalCount,
		},
	})
}

func (h *Handler) HandlePlayGacha(c *gin.Context) {
	claim, ok := c.Get("user")
	user, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
		return
	}

	voucher, err := h.promoUsecase.IssueGachaVoucher(user.UserId)
	if err != nil {
		if errors.Is(err, usecase_errors.ErrRecordNotExist) {
			c.Error(httperror.BadRequestError("Bad Request", "BAD_REQUEST"))
			return
		}
		if errors.Is(err, usecase_errors.ErrInsufficientGachaQuota) {
			c.Error(httperror.BadRequestError("Insufficient Quota", "BAD_REQUEST"))
			return
		}
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
		"data":        voucher,
	})

}
