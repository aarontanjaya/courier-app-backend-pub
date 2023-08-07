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

func (h *Handler) HandleGetUserAddresses(c *gin.Context) {
	claim, ok := c.Get("user")
	user, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
		return
	}
	search := c.DefaultQuery("search", "")
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
	paging := dto.AddressPaginationRequest{
		Search:   search,
		PageSize: pageSize,
		Page:     page,
		UserId:   &user.UserId,
	}
	addresses, pagination, err := h.addressUsecase.GetAll(&paging)
	if err != nil {
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
		"data": dto.AddressResponsePayload{
			Records:    addresses,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			PageCount:  pagination.PageCount,
			TotalCount: pagination.TotalCount,
		},
	})
}

func (h *Handler) HandleGetAdminAddresses(c *gin.Context) {
	claim, ok := c.Get("user")
	_, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
		return
	}

	search := c.DefaultQuery("search", "")
	userIdStr := c.DefaultQuery("user_id", "0")
	if userIdStr == "" {
		userIdStr = "0"
	}
	userIdInt, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.Error(httperror.BadRequestError("Bad Request, user_id invalid", "BAD_REQUEST"))
		return
	}
	userId := uint(userIdInt)

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
	paging := dto.AddressPaginationRequest{
		Search:   search,
		PageSize: pageSize,
		Page:     page,
		UserId:   &userId,
	}
	addresses, pagination, err := h.addressUsecase.GetAll(&paging)
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
		"data": dto.AddressResponsePayload{
			Records:    addresses,
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			PageCount:  pagination.PageCount,
			TotalCount: pagination.TotalCount,
		},
	})
}

func (h *Handler) HandleUpdateAddress(c *gin.Context) {
	var req *dto.AddressUpdateRequest
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
	req.ID = uint(id)
	err = h.addressUsecase.UpdateAddress(user.UserId, req)
	if err != nil {
		if errors.Is(err, usecase_errors.ErrRecordNotExist) {
			c.Error(httperror.NotFoundError("Resource not found"))
			return
		}
		if errors.Is(err, usecase_errors.ErrUserNotAuthorized) {
			c.Error(httperror.UnauthorizedError())
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

func (h *Handler) HandleCreateAddress(c *gin.Context) {
	var req *dto.AddressCreateRequest
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

	address, err := h.addressUsecase.CreateAddress(user.UserId, req)
	if err != nil {
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
		"data":        address,
	})
}

func (h *Handler) HandleDeleteAdress(c *gin.Context) {
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

	res, err := h.addressUsecase.DeleteAddress(uint(id), user.UserId)
	if err != nil {
		if errors.Is(err, usecase_errors.ErrRecordNotExist) {
			c.Error(httperror.NotFoundError("Resource not found"))
			return
		}
		if errors.Is(err, usecase_errors.ErrUserNotAuthorized) {
			c.Error(httperror.UnauthorizedError())
			return
		}
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
		"data":        res,
	})
}
