package jwt

import (
	"errors"
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

			ExpiresAt: time.Now().Add(time.Minute * 120).Unix(),
		},
		"level1",
		CustomerInfo{uuid, username, "human"},
	}

	return t.SignedString(viper.GetString("secrets.jwt"))
}

func VerifyToken(tokenString string) (*CustomerInfo, error){
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return viper.GetString("secrets.jwt"), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}
	return &claims.CustomerInfo, nil
}

func RefreshToken(tokenString string) (string, error) {
	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return viper.GetString("secrets.jwt"), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", jwt.ErrInvalidKey
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		return "", errors.New("token still valid")
	}

	expirationTime := time.Now().Add(120 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(viper.GetString("secrets.jwt"))
	if err != nil {
		return "", err
	}
	return tokenString, nil
} 