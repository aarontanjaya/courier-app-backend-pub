package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"courier-app/config"
	"courier-app/dto"
	"courier-app/httperror"
	"courier-app/usecase"
	"courier-app/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func SetHeaderFromCookie(c *gin.Context) {
	token, err := c.Request.Cookie("token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httperror.UnauthorizedError())
		return
	}
	c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.Value))
}

func Authorize(c *gin.Context) {
	if config.Config.ENV == "testing" {
		fmt.Println("testing mode, authorization off")
		return
	}

	authHeader := c.GetHeader("Authorization")
	s := strings.Split(authHeader, "Bearer ")
	if len(s) < 2 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httperror.UnauthorizedError())
		return
	}
	encodedToken := s[1]

	token, err := util.ValidateToken(encodedToken)
	if err != nil || !token.Valid {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httperror.UnauthorizedError())
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httperror.UnauthorizedError())
		return
	}

	userJson, _ := json.Marshal(claims["user"])
	var user dto.UserClaims
	err = json.Unmarshal(userJson, &user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httperror.UnauthorizedError())
		return
	}
	newToken, err := usecase.GenerateAccessToken(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, httperror.InternalServerError("Internal Server Error"))
	}

	authConfig := config.Config.AuthConfig
	duration, err := strconv.Atoi(authConfig.TokenDur)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, httperror.InternalServerError("Internal Server Error"))
	}
	//c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("token", newToken, duration, "/", "", false, true)

	c.Set("user", user)
	fmt.Printf("loggedin user %+v", user)
}

func AuthorizeUser(c *gin.Context) {
	claim, ok := c.Get("user")
	user, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 || !strings.Contains(user.Scope, "user") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httperror.UnauthorizedError())
	}
}

func AuthorizeAdmin(c *gin.Context) {
	claim, ok := c.Get("user")
	user, ok1 := claim.(dto.UserClaims)
	if !ok || !ok1 || !strings.Contains(user.Scope, "admin") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httperror.UnauthorizedError())
	}
}
