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

func (h *Handler) HandleCreateShipping(c *gin.Context) {
	var req *dto.ShippingCreateRequest
	claim, ok := c.Get("user")
	user, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
		return
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(httperror.BadRequestError("Bad Request", "BAD_REQUEST"))
		return
	}

	shipping, err := h.shippingUsecase.CreateShipping(user.UserId, req)
	if err != nil {
		if errors.Is(err, usecase_errors.ErrRecordNotExist) || errors.Is(err, usecase_errors.ErrUserNotAuthorized) {
			c.Error(httperror.BadRequestError("Bad Request", "BAD_REQUEST"))
			return
		}
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
		"data":        shipping,
	})
}

func (h *Handler) HandleGetShippingById(c *gin.Context) {
	claim, ok := c.Get("user")
	user, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
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

	shipping, err := h.shippingUsecase.GetById(uint(id))
	if shipping.UserId != user.UserId {
		c.Error(httperror.NotFoundError("Not Found"))
		return
	}
	if err != nil {
		if errors.Is(err, usecase_errors.ErrRecordNotExist) {
			c.JSON(http.StatusNotFound, httperror.NotFoundError("Not Found"))
			return
		}
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
		"data":        shipping,
	})

}

func (h *Handler) HandleGetShippings(c *gin.Context) {
	var reqParams *dto.ShippingTableRequest
	claim, ok := c.Get("user")
	user, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
		return
	}

	err := c.BindQuery(&reqParams)
	if err != nil {
		c.Error(httperror.BadRequestError("Bad Request", "BAD_REQUEST"))
		return
	}

	shippings, pagination, err := h.shippingUsecase.GetByUserId(user.UserId, reqParams)
	if err != nil {
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
		"data": dto.ShippingResponsePayload{
			Records:    shippings,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			PageCount:  pagination.PageCount,
			TotalCount: pagination.TotalCount,
		}})
}

func (h *Handler) HandleGetAllShippings(c *gin.Context) {
	var reqParams *dto.ShippingTableRequest
	claim, ok := c.Get("user")
	_, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
		return
	}

	err := c.BindQuery(&reqParams)
	if err != nil {
		c.Error(httperror.BadRequestError("Bad Request", "BAD_REQUEST"))
		return
	}

	shippings, pagination, err := h.shippingUsecase.GetAll(reqParams)
	if err != nil {
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
		"data": dto.ShippingResponsePayload{
			Records:    shippings,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			PageCount:  pagination.PageCount,
			TotalCount: pagination.TotalCount,
		}})
}

func (h *Handler) HandleGetShippingStatuses(c *gin.Context) {
	claim, ok := c.Get("user")
	_, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
		return
	}

	statuses, err := h.shippingUsecase.GetShippingStatuses()
	if err != nil && !errors.Is(err, usecase_errors.ErrRecordNotExist) {
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
		"data":        statuses,
	})
}

func (h *Handler) HandleUpdateShippingStatus(c *gin.Context) {
	var req *dto.ShippingStatusUpdateRequest
	claim, ok := c.Get("user")
	_, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
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

	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(httperror.BadRequestError("Bad Request", "BAD_REQUEST"))
		return
	}

	err = h.shippingUsecase.UpdateStatus(uint(id), req.StatusId)
	if err != nil {
		if errors.Is(err, usecase_errors.ErrRecordNotExist) {
			c.Error(httperror.BadRequestError("Bad Request", "BAD_REQUEST"))
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

func (h *Handler) HandleReviewShipping(c *gin.Context) {
	var req *dto.ShippingReviewRequest
	claim, ok := c.Get("user")
	user, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
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

	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(httperror.BadRequestError("Bad Request", "BAD_REQUEST"))
		return
	}

	req.ShippingId = uint(id)

	err = h.shippingUsecase.Review(user.UserId, req)
	if err != nil {
		if errors.Is(err, usecase_errors.ErrRecordNotExist) || errors.Is(err, usecase_errors.ErrUserNotAuthorized) || errors.Is(err, usecase_errors.ErrInvalidReq) {
			c.Error(httperror.BadRequestError("Bad Request", "BAD_REQUEST"))
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
