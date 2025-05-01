package utils

import (
	"fmt"
	"thanhldt060802/config"
	"thanhldt060802/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user *model.User) (string, error) {
	expireDuration, err := config.AppConfig.GetTokenExpireMinutes()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id":   user.Id,
		"role_name": user.RoleName,
		"exp":       time.Now().Add(*expireDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("generate token failed")
	}

	return tokenStr, nil
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
