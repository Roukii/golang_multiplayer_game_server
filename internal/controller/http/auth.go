package http

import (
	"net/http"

	"github.com/Roukii/pock_multiplayer/internal/middleware"
	"github.com/Roukii/pock_multiplayer/internal/service"
	"github.com/Roukii/pock_multiplayer/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type authRoutes struct {
	services *service.Service
}

func newAuthRoutes(handler *gin.RouterGroup,  services *service.Service) {
	r := &authRoutes{services}

	h := handler.Group("/auth")
	{
		h.POST("/login", r.login)
		h.POST("/logout", r.logout)
		h.POST("/refresh", r.refresh)
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

type user struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (ar *authRoutes) login(c *gin.Context) {
  var u user
  if err := c.ShouldBindJSON(&u); err != nil {
     c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
     return
  }
	if u.Username != u.Username || u.Password != u.Password {
     c.JSON(http.StatusUnauthorized, "Please provide valid login details")
     return
  }
  token, err := jwt.CreateToken(u.UUID, u.Username)
  if err != nil {
     c.JSON(http.StatusUnprocessableEntity, err.Error())
     return
  }
  c.JSON(http.StatusOK, token)
}

func (ar *authRoutes) logout(c *gin.Context) {

}

func (ar *authRoutes) refresh(c *gin.Context) {

}

func (ar *authRoutes) resetPassword(c *gin.Context) {

}

func (ar *authRoutes) register(c *gin.Context) {

}
