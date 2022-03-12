package http

import (
	"fmt"
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
		return
	}
	user, err := ar.services.UserService.Login(u.Username, u.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		fmt.Println(err)
		ar.services.Logger.Error("Login failed for username " + u.Username)
		return
	}

	token, err := jwt.CreateToken(user.UUID, user.Username)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusOK, token)
}

// TODO
func (ar *authRoutes) logout(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

func (ar *authRoutes) refresh(c *gin.Context) {
	token := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
	if token == "" {
		c.JSON(http.StatusUnprocessableEntity, "token not found")
		return
	}
	refreshedToken, err := jwt.RefreshToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
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
		return
	}
	_, err := ar.services.UserService.Register(u)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Please provide valid register details")
		ar.services.Logger.Error("Login failed for username " + u.Username)
		return
	}
	c.JSON(http.StatusOK, true)
}
