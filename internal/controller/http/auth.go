package http

import (
	"net/http"

	"github.com/Roukii/pock_multiplayer/pkg/logger"
	"github.com/Roukii/pock_multiplayer/pkg/jwt"
	"github.com/gin-gonic/gin"

)

type authRoutes struct {
	l logger.Interface
}

func newAuthRoutes(handler *gin.Engine, l logger.Interface) {
	r := &authRoutes{l}

	h := handler.Group("/auth")
	{
		h.POST("/login", r.login)
		h.POST("/logout", r.logout)
		h.POST("/refresh", r.refresh)
		h.POST("/reset_password", r.resetPassword)
		h.POST("/register", r.register)
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
  token, err := CreateToken(u.UUID)
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
