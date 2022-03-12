package middleware

import (
	"net/http"
	"strings"

	"github.com/Roukii/pock_multiplayer/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func CheckTokenJWT(c *gin.Context) {
	reqToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer")
	if len(splitToken) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	token := strings.TrimSpace(splitToken[1])

	claims, err := jwt.VerifyToken(token)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Set("uuid", claims.UUID)
	c.Set("username", claims.Name)
	c.Set("device", claims.Device)
	c.Next()
}
