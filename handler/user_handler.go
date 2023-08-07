package handler

import (
	"courier-app/dto"
	"courier-app/httperror"
	"courier-app/usecase/usecase_errors"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (h *Handler) HandleGetProfile(c *gin.Context) {
	claim, ok := c.Get("user")
	user, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
		return
	}

	details, err := h.userUsecase.GetById(user.UserId)
	if err != nil {
		if errors.Is(err, usecase_errors.ErrRecordNotExist) {
			c.Error(httperror.NotFoundError("Not Found"))
			return
		}
		c.Error(httperror.InternalServerError("Internal Server Error"))
		return
	}
	response := dto.ProfileDetails{
		Name:         details.Name,
		Email:        details.Email,
		Phone:        details.Phone,
		Photo:        details.Photo,
		Role:         details.RoleName,
		ReferralCode: details.ReferralCode,
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"code":        "OK",
		"data":        response,
	})
}

func (h *Handler) HandleGetBalance(c *gin.Context) {
	claim, ok := c.Get("user")
	user, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
		return
	}

	res, err := h.userUsecase.GetUserBalance(user.UserId)
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
		"data":        res,
	})
}

func (h *Handler) HandleGetGachaQuota(c *gin.Context) {
	claim, ok := c.Get("user")
	user, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
		return
	}

	res, err := h.userUsecase.GetUserGachaQuota(user.UserId)
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
		"data":        res,
	})
}

func (h *Handler) HandleUpdateProfile(c *gin.Context) {
	var req *dto.UpdateProfileRequest
	var reqText *dto.UpdateProfileContentRequest
	claim, ok := c.Get("user")
	user, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 {
		c.Error(httperror.UnauthorizedError())
		return
	}
	p, _ := c.FormFile("photo")
	if p == nil {
		validate := validator.New()
		name := c.PostForm("name")
		email := c.PostForm("email")
		phone := c.PostForm("phone")
		reqText = &dto.UpdateProfileContentRequest{
			Name:  name,
			Email: email,
			Phone: phone,
		}
		err := validate.Struct(reqText)
		if err != nil {
			validationErrors, ok := err.(validator.ValidationErrors)
			if !ok {
				c.Error(httperror.BadRequestError("Bad Request", "BAD_REQUEST"))
				return
			}
			c.Error(httperror.BadRequestError(fmt.Sprintf("Field %v invalid", validationErrors[0].Tag()), "BAD_REQUEST"))
			return

		}
		req = &dto.UpdateProfileRequest{
			Name:  reqText.Name,
			Email: reqText.Email,
			Phone: reqText.Phone,
		}
	} else {
		err := c.ShouldBind(&req)
		if err != nil {
			c.Error(httperror.BadRequestError("Bad Request", "BAD_REQUEST"))
			return
		}

	}
	err := h.userUsecase.UpdateProfile(&dto.UpdateProfileRequest{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
		Photo: req.Photo,
	}, user.UserId)
	if err != nil {
		if errors.Is(err, usecase_errors.ErrEmailAlreadyExist) {
			c.Error(httperror.BadRequestError("Email already exist", "BAD_REQUEST"))
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
