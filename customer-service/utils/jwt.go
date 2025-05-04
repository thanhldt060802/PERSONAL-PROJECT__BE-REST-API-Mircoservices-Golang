package utils

import (
	"fmt"
	"thanhldt060802/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userId int64, roleName string, cartId int64) (*string, error) {
	expireDuration := config.AppConfig.GetTokenExpireMinutes()
	if expireDuration == nil {
		return nil, fmt.Errorf("convert expire failed")
	}

	claims := jwt.MapClaims{
		"user_id":   userId,
		"role_name": roleName,
		"cart_id":   cartId,
		"exp":       time.Now().Add(*expireDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return nil, fmt.Errorf("generate token failed")
	}

	return &tokenStr, nil
}

// func ValidateToken(tokenStr string) (jwt.MapClaims, error) {
// 	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, jwt.ErrSignatureInvalid
// 		}
// 		return []byte(config.AppConfig.JWTSecret), nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		return claims, nil
// 	}
// 	return nil, jwt.ErrSignatureInvalid
// }
