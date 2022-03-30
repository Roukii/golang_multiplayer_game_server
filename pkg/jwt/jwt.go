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

func CreateToken(uuid string, email string) (string, string, error) {
	accessToken := jwt.New(jwt.GetSigningMethod("HS512"))
	accessToken.Claims = &CustomClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		},
		"access_token",
		CustomerInfo{uuid, email, "test"},
	}

	refreshToken := jwt.New(jwt.GetSigningMethod("HS512"))
	refreshToken.Claims = &CustomClaims{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(),
		},
		"refresh_token",
		CustomerInfo{uuid, email, "test"},
	}
	accessTokenValue, err := accessToken.SignedString([]byte(viper.GetString("secrets.jwt")))
	if err != nil {
		return "", "", err
	}
	refreshTokenValue, err := refreshToken.SignedString([]byte(viper.GetString("secrets.jwt")))
	if err != nil {
		return "", "", err
	}
	return accessTokenValue, refreshTokenValue, err
}

func VerifyToken(tokenString string) (*CustomerInfo, error){
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("secrets.jwt")), nil
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
		return []byte(viper.GetString("secrets.jwt")), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", jwt.ErrInvalidKey
	}

	if claims.TokenType != "refresh_token" {
		return "", errors.New("token not refreshToken")
	}

	expirationTime := time.Now().Add(time.Minute * 60)
	claims.ExpiresAt = expirationTime.Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(viper.GetString("secrets.jwt")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
} 