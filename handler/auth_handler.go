package handler

import (
	"courier-app/config"
	"courier-app/dto"
	"courier-app/httperror"
	"courier-app/usecase/usecase_errors"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) HandleRegister(c *gin.Context) {
	var req *dto.RegisterRequest
	err := c.ShouldBind(&req)
	if err != nil {
		c.Error(httperror.BadRequestError("Bad Request", "BAD_REQUEST"))
		return
	}

	res, err := h.authUsecase.Register(req)
	if err != nil {
		if errors.Is(err, usecase_errors.ErrEmailAlreadyExist) {
			c.Error(httperror.BadRequestError("Email already exist", "BAD_REQUEST"))
			return
		}
		if errors.Is(err, usecase_errors.ErrReferralInvalid) {
			c.Error(httperror.BadRequestError("Referral code invalid", "BAD_REQUEST"))
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

func (h *Handler) HandleLogin(c *gin.Context) {
	var req *dto.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(httperror.BadRequestError("Bad Request", "BAD_REQUEST"))
		return
	}

	res, err := h.authUsecase.Login(req)
	if err != nil {
		if errors.Is(err, usecase_errors.ErrCredsInvalid) {
			c.Error(httperror.UnauthorizedError())
			return
		}
		c.Error(err)
		return
	}

	authConfig := config.Config.AuthConfig
	duration, err := strconv.Atoi(authConfig.TokenDur)
	if err != nil {
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}
	//c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("token", res.AccessToken, duration, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
		"data":        res,
	})
}

func (h *Handler) HandleLogout(c *gin.Context) {
	//c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("token", "", 1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
	})
}

func (h *Handler) HandleWhoAmI(c *gin.Context) {
	claim, ok := c.Get("user")
	user, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
		"data": dto.LoginResponse{
			UserId: user.UserId,
			Scope:  user.Scope,
		},
	})
}
