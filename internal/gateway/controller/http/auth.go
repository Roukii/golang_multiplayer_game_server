package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Roukii/pock_multiplayer/internal/gateway/middleware"
	"github.com/Roukii/pock_multiplayer/internal/gateway/service"
	"github.com/Roukii/pock_multiplayer/internal/gateway/service/user"
	"github.com/Roukii/pock_multiplayer/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type authRoutes struct {
	services *service.Service
}

func newAuthRoutes(handler *gin.RouterGroup, services *service.Service) {
	r := &authRoutes{services}

	h := handler.Group("/auth")
	{
		h.POST("/login", r.login)
		h.POST("/reset_password", r.resetPassword)
		h.POST("/register", r.register)
	}

	auth := h.Group("")
	auth.Use(middleware.CheckTokenJWT)
	{
		auth.POST("/logout", r.logout)
		auth.POST("/refresh", r.refresh)
	}
}

func (ar *authRoutes) login(c *gin.Context) {
	var u user.UserInput
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		ar.services.Logger.Error("Invalid json for login " + err.Error())
		return
	}
	user, err := ar.services.UserService.Login(u.Email, u.Password, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		ar.services.Logger.Error("Login failed for ", u.Email, " with error : ", err.Error())
		return
	}

	accessToken, refreshToken, err := jwt.CreateToken(user.UUID, user.Email)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		ar.services.Logger.Error("Couldn't create token for ", u.Email, " with error : ", err.Error())
		return
	}
	ar.services.Logger.Error("Login succesful for ", u.Email)
	returnData := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	res, err := json.Marshal(returnData)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		ar.services.Logger.Error("Couldn't marshal token for ", u.Email, " with error : ", err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

// TODO
func (ar *authRoutes) logout(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

func (ar *authRoutes) refresh(c *gin.Context) {
	token := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
	if token == "" {
		c.JSON(http.StatusUnprocessableEntity, "token not found")
		ar.services.Logger.Error("Token not found")
		return
	}
	refreshedToken, err := jwt.RefreshToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		ar.services.Logger.Error("Couldn't refresh token for ", c.GetString("Email"), " with error : ", err.Error())
		return
	}
	ar.services.Logger.Error("Refresh token for ", c.GetString("Email"))
	c.JSON(http.StatusOK, refreshedToken)
}

// TODO
func (ar *authRoutes) resetPassword(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, "rip")
}

func (ar *authRoutes) register(c *gin.Context) {
	var u user.UserInput
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		ar.services.Logger.Error("Invalid json for register ", err.Error())
		return
	}
	_, err := ar.services.UserService.Register(u)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Please provide valid register details")
		ar.services.Logger.Error("Register failed for Email ", u.Email, " with error : ", err.Error())
		return
	}
	c.JSON(http.StatusOK, true)
	ar.services.Logger.Error("Register succesful for ", u.Email)
}
