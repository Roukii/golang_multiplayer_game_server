package middleware

import (
	"net/http"
	"strings"

	"github.com/Roukii/pock_multiplayer/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func CheckTokenJWT(c *gin.Context) {
	token := strings.Split(c.Request.Header["Authorization"][0], " ")[1]

	claims, err := jwt.VerifyToken(token)
	if err !=nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	c.Set("uuid", claims.UUID)
	c.Set("username", claims.Name)
	c.Set("device", claims.Device)
  c.Next()
}
