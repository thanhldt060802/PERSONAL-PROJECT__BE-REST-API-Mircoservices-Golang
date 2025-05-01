package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/redis/go-redis/v9"
)

type AuthMiddleware struct {
	API         huma.API
	redisClient *redis.Client
}

func NewAuthMiddleware(api huma.API, redisClient *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{
		API:         api,
		redisClient: redisClient,
	}
}

func (authMiddleware *AuthMiddleware) Authentication(ctx huma.Context, next func(huma.Context)) {
	authHeader := ctx.Header("Authorization")
	if authHeader == "" {
		huma.WriteErr(authMiddleware.API, ctx, http.StatusUnauthorized, "Authorization header missing")
		return
	}

	// Simple way
	// ################################################################################
	// parts := strings.Split(authHeader, " ")
	// if len(parts) != 2 || parts[0] != "Bearer" {
	// 	huma.WriteErr(authMiddleware.API, ctx, http.StatusUnauthorized, "Token missing")
	// 	return
	// }

	// token := parts[1]
	// claims, err := utils.ValidateToken(token)
	// if err != nil {
	// 	huma.WriteErr(authMiddleware.API, ctx, http.StatusUnauthorized, "Invalid token")
	// 	return
	// }

	// userId, ok := claims["user_id"].(float64)
	// if !ok {
	// 	huma.WriteErr(authMiddleware.API, ctx, http.StatusUnauthorized, "Field user_id is not valid in token")
	// 	return
	// }
	// ctx = huma.WithValue(ctx, "user_id", int64(userId))

	// roleName, ok := claims["role_name"].(string)
	// if !ok {
	// 	huma.WriteErr(authMiddleware.API, ctx, http.StatusUnauthorized, "Field role_name is not valid in token")
	// 	return
	// }
	// ctx = huma.WithValue(ctx, "role_name", roleName)
	// ################################################################################

	// Redis way
	// ################################################################################
	token := strings.TrimPrefix(authHeader, "Bearer ")
	redisKey := fmt.Sprintf("token:%s", token)

	userDataJson, err := authMiddleware.redisClient.Get(ctx.Context(), redisKey).Result()
	if err == redis.Nil {
		huma.WriteErr(authMiddleware.API, ctx, http.StatusUnauthorized, "Token not found or expired")
		return
	} else if err != nil {
		huma.WriteErr(authMiddleware.API, ctx, http.StatusUnauthorized, "Failed to check token in Redis")
		return
	}

	var userData struct {
		UserId   int64  `json:"user_id"`
		RoleName string `json:"role_name"`
	}

	if err := json.Unmarshal([]byte(userDataJson), &userData); err != nil {
		huma.WriteErr(authMiddleware.API, ctx, http.StatusUnauthorized, "Invalid user data in token")
		return
	}

	ctx = huma.WithValue(ctx, "user_id", userData.UserId)
	ctx = huma.WithValue(ctx, "role_name", userData.RoleName)
	// ################################################################################

	next(ctx)
}

func (authMiddleware *AuthMiddleware) RequireAdmin(ctx huma.Context, next func(huma.Context)) {
	roleName, ok := ctx.Context().Value("role_name").(string)
	if !ok || roleName != "ADMIN" {
		huma.WriteErr(authMiddleware.API, ctx, http.StatusUnauthorized, "No permission (only admin)")
		return
	}

	next(ctx)
}
