package handler

import (
	"courier-app/dto"
	"courier-app/httperror"
	"courier-app/usecase/usecase_errors"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleGetAllAddOns(c *gin.Context) {
	claim, ok := c.Get("user")
	_, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
		return
	}

	addons, err := h.addOnUsecase.GetAll()
	if err != nil && !errors.Is(err, usecase_errors.ErrRecordNotExist) {
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
		"data":        addons,
	})
}
