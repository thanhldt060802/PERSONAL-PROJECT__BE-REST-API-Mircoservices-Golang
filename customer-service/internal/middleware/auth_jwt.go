package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"thanhldt060802/infrastructure"

	"github.com/danielgtaylor/huma/v2"
	"github.com/redis/go-redis/v9"
)

type AuthMiddleware struct {
	API huma.API
}

func NewAuthMiddleware(api huma.API) *AuthMiddleware {
	return &AuthMiddleware{
		API: api,
	}
}

func (authMiddleware *AuthMiddleware) Authentication(ctx huma.Context, next func(huma.Context)) {
	authHeader := ctx.Header("Authorization")
	if authHeader == "" {
		CustomerHumaWriteErr(ctx, http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Authorization header missing", []string{"invalid credentials"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	redisKey := fmt.Sprintf("token:%s", token)

	userDataJson, err := infrastructure.RedisClient.Get(ctx.Context(), redisKey).Result()
	if err == redis.Nil {
		CustomerHumaWriteErr(ctx, http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Token not found or expired", []string{"invalid token"})
		return
	} else if err != nil {
		CustomerHumaWriteErr(ctx, http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Failed to check token in Redis", []string{"some thing wrong in redis"})
		return
	}

	var userData struct {
		UserId   int64  `json:"user_id"`
		RoleName string `json:"role_name"`
		CartId   int64  `json:"cart_id"`
	}

	if err := json.Unmarshal([]byte(userDataJson), &userData); err != nil {
		CustomerHumaWriteErr(ctx, http.StatusUnauthorized, "ERR_UNAUTHORIZED", "Invalid user data in token", []string{"invalid token"})
		return
	}

	ctx = huma.WithValue(ctx, "user_id", userData.UserId)
	ctx = huma.WithValue(ctx, "role_name", userData.RoleName)
	ctx = huma.WithValue(ctx, "cart_id", userData.CartId)

	next(ctx)
}

func (authMiddleware *AuthMiddleware) RequireAdmin(ctx huma.Context, next func(huma.Context)) {
	if roleName, _ := ctx.Context().Value("role_name").(string); roleName != "ADMIN" {
		CustomerHumaWriteErr(ctx, http.StatusForbidden, "ERR_FORBIDDEN", "Access denied", []string{"no permission"})
		return
	}

	next(ctx)
}
