package usecase

import (
	"courier-app/app_constant"
	"courier-app/config"
	"courier-app/dto"
	"courier-app/entity"
	"courier-app/usecase/usecase_errors"
	"courier-app/util"
	crypto_rand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type idTokenClaims struct {
	jwt.RegisteredClaims
	User dto.UserClaims `json:"user"`
}

type AuthUsecase interface {
	Register(creds *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(creds *dto.LoginRequest) (*dto.LoginResponse, error)
}

type authUsecaseImpl struct {
	userUsecase UserUsecase
}

type AuthUsecaseConfig struct {
	UserUsecase UserUsecase
}

func NewAuthUsecase(c AuthUsecaseConfig) AuthUsecase {
	return &authUsecaseImpl{
		userUsecase: c.UserUsecase,
	}
}

func (au *authUsecaseImpl) Register(creds *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	var res *dto.RegisterResponse
	referralStatus := 1
	if creds.ClaimedReferral != "" {
		_, err := au.userUsecase.GetByReferral(creds.ClaimedReferral)
		if err != nil {
			if err == usecase_errors.ErrRecordNotExist {
				return res, usecase_errors.ErrReferralInvalid
			}
			return res, err
		}
		referralStatus = 2
	}

	var referral string
	for {
		var err error
		referral, err = generateReferral()
		if err != nil {
			return res, err
		}

		_, err = au.userUsecase.GetByReferral(referral)
		if err == usecase_errors.ErrReferralInvalid {
			break
		}

		if err != nil {
			return res, err
		}

	}

	password, err := util.HashAndSalt(creds.Password)
	if err != nil {
		return res, usecase_errors.ErrGeneratePassword
	}

	userDetail := entity.UserDetail{
		ReferralCode:    referral,
		ClaimedReferral: creds.ClaimedReferral,
		ReferralStatus:  uint(referralStatus),
	}

	res, err = au.userUsecase.Register(&entity.User{
		Name:       creds.Name,
		Email:      creds.Email,
		Password:   password,
		UserDetail: userDetail,
		Role:       uint(app_constant.RoleUser),
		Phone:      creds.Phone,
	})

	return res, err

}

func (au *authUsecaseImpl) Login(creds *dto.LoginRequest) (*dto.LoginResponse, error) {
	token, err := au.AuthenticateUser(creds)
	return token, err
}

func (au *authUsecaseImpl) AuthenticateUser(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	var token string
	var res *dto.LoginResponse
	user, err := au.userUsecase.GetByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return res, usecase_errors.ErrCredsInvalid
		}
		return res, err
	}
	authValid := ComparePassword(user.Password, req.Password)
	if !authValid {
		return res, usecase_errors.ErrCredsInvalid
	}

	token, err = GenerateAccessToken(&dto.UserClaims{
		UserId: user.ID,
		Scope:  user.RoleName,
	})
	res = &dto.LoginResponse{
		AccessToken: token,
		UserId:      user.ID,
		Scope:       user.RoleName,
	}
	return res, err
}

func GenerateAccessToken(user *dto.UserClaims) (string, error) {
	var tokenString string
	authConfig := config.Config.AuthConfig
	now := time.Now()
	duration, err := strconv.Atoi(authConfig.TokenDur)
	if err != nil {
		return tokenString, err
	}
	end := now.Add(time.Second * time.Duration(duration))
	claims := &idTokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(end),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		*user,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	hmacSampleSecret := config.Config.AuthConfig.Secret
	tokenString, err = token.SignedString([]byte(hmacSampleSecret))
	return tokenString, err
}

func ComparePassword(hashedPwd string, inputPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(inputPwd))
	return err == nil
}

func generateReferral() (string, error) {
	referralLength := 16
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())

	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		return "", usecase_errors.ErrReferralFailed
	}

	res := make([]byte, referralLength)
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
	for i := range res {
		res[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(res), nil
}
