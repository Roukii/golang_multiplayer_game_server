package middleware

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type CustomerInfo struct {
	UUID   string
	Name   string
	Device string
}

type CustomClaims struct {
	*jwt.StandardClaims
	TokenType string
	CustomerInfo
}

func CreateToken(uuid string, username string) (string, error) {
	t := jwt.New(jwt.GetSigningMethod("HS512"))
	t.Claims = &CustomClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(),
		},
		"level1",
		CustomerInfo{uuid, username, "human"},
	}

	return t.SignedString(viper.GetString("secrets.jwt"))
}
