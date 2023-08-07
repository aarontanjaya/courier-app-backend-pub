package handler

import (
	"courier-app/config"
	"courier-app/dto"
	"courier-app/httperror"
	"courier-app/usecase/usecase_errors"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (h *Handler) HandleTopUp(c *gin.Context) {
	var req *dto.TopUpRequestBody

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

	p := message.NewPrinter(language.Indonesian)
	topupConfig := config.Config.TopupConfig
	res, err := h.transactionUsecase.TopUp(user.UserId, req)
	if err != nil {
		if errors.Is(err, usecase_errors.ErrTopupAmountInvalid) {
			c.Error(httperror.BadRequestError(p.Sprintf("Amount invalid, minimum IDR %.2f, maximum IDR %.2f", topupConfig.TopupMin, topupConfig.TopupMax), "BAD_REQUEST"))
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
