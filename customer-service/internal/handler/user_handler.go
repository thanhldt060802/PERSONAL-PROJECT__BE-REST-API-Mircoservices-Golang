package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"thanhldt060802/config"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/middleware"
	"thanhldt060802/internal/service"
	"thanhldt060802/utils"

	"github.com/danielgtaylor/huma/v2"
	"github.com/redis/go-redis/v9"
)

type UserHandler struct {
	userService    service.UserService
	authMiddleware *middleware.AuthMiddleware
	redisClient    *redis.Client
}

func NewUserHandler(api huma.API, userService service.UserService, authMiddleware *middleware.AuthMiddleware, redisClient *redis.Client) *UserHandler {
	userHandler := &UserHandler{
		userService:    userService,
		authMiddleware: authMiddleware,
		redisClient:    redisClient,
	}

	// Testing for creating user
	// ################################################################################
	huma.Register(api, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/testcreate",
		Summary: "/testcreate",
		Tags:    []string{"Test"},
	}, userHandler.CreateUser)
	// ################################################################################

	// Get users
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/users",
		Summary:     "/users",
		Description: "Get users.",
		Tags:        []string{"User"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, userHandler.GetUsers)

	// Get user by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/users/id/{id}",
		Summary:     "/users/id/{id}",
		Description: "Get user by id.",
		Tags:        []string{"User"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, userHandler.GetUserById)

	// Create user
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/users",
		Summary:     "/users",
		Description: "Create user.",
		Tags:        []string{"User"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, userHandler.CreateUser)

	// Update user by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/users/id/{id}",
		Summary:     "/users/id/{id}",
		Description: "Update user by id.",
		Tags:        []string{"User"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, userHandler.UpdateUserById)

	// Delete user by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/users/id/{id}",
		Summary:     "/users/id/{id}",
		Description: "Delete user by id.",
		Tags:        []string{"User"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, userHandler.DeleteUserById)

	// Login account
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "/login",
		Description: "Login account",
		Tags:        []string{"Auth"},
	}, userHandler.Login)

	// Register account
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/register",
		Summary:     "/register",
		Description: "Register account.",
		Tags:        []string{"User"},
	}, userHandler.Register)

	// Get account info
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/account-info",
		Summary:     "/account-info",
		Description: "Get account info.",
		Tags:        []string{"User"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication},
	}, userHandler.GetAccountInfo)

	// Update account info
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/account-info",
		Summary:     "/account-info",
		Description: "Update account info.",
		Tags:        []string{"User"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication},
	}, userHandler.UpdateAccountInfo)

	return userHandler
}

func (userHandler *UserHandler) GetUsers(ctx context.Context, queryParam *dto.GetUsersRequestQueryParam) (*dto.SuccessResponse[[]dto.UserDTO], error) {
	users, err := userHandler.userService.GetUsers(ctx, queryParam)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Message = "Get users failed"
		res.Error_ = "Internal Server Error"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToUserDTOs(users)
	res := &dto.SuccessResponse[[]dto.UserDTO]{}
	res.Body.Status = http.StatusOK
	res.Body.Message = "Get users successful"
	res.Body.Data = data
	res.Body.Total = len(data)
	return res, nil
}

func (userHandler *UserHandler) GetUserById(ctx context.Context, reqDTO *dto.GetUserByIdRequest) (*dto.SuccessResponse[*dto.UserDTO], error) {
	id := reqDTO.Id

	foundUser, err := userHandler.userService.GetUserById(ctx, id)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Message = "Get user by id failed"
		res.Error_ = "Bad Request"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToUserDTO(foundUser)
	res := &dto.SuccessResponse[*dto.UserDTO]{}
	res.Body.Status = http.StatusOK
	res.Body.Message = "Get user by id successful"
	res.Body.Data = data
	return res, nil
}

func (userHandler *UserHandler) CreateUser(ctx context.Context, reqDTO *dto.CreateUserRequest) (*dto.SuccessResponse[any], error) {
	if err := userHandler.userService.CreateUser(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Message = "Create user failed"
		res.Error_ = "Bad Request"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse[any]{}
	res.Body.Status = http.StatusOK
	res.Body.Message = "Create user successful"
	return res, nil
}

func (userHandler *UserHandler) UpdateUserById(ctx context.Context, reqDTO *dto.UpdateUserRequest) (*dto.SuccessResponse[any], error) {
	id := reqDTO.Id

	if err := userHandler.userService.UpdateUserById(ctx, id, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Message = "Update user failed"
		res.Error_ = "Bad Request"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse[any]{}
	res.Body.Status = http.StatusOK
	res.Body.Message = "Update user successful"
	return res, nil
}

func (userHandler *UserHandler) DeleteUserById(ctx context.Context, reqDTO *dto.DeleteUserRequest) (*dto.SuccessResponse[any], error) {
	id := reqDTO.Id

	if err := userHandler.userService.DeleteUserById(ctx, id); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Message = "Delete user failed"
		res.Error_ = "Bad Request"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse[any]{}
	res.Body.Status = http.StatusOK
	res.Body.Message = "Delete user successful"
	return res, nil
}

func (userHandler *UserHandler) Login(ctx context.Context, reqDTO *dto.LoginRequest) (*dto.SuccessResponse[string], error) {
	foundUser, err := userHandler.userService.GetUserByUsername(ctx, reqDTO.Body.Username)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Message = "Login failed"
		res.Error_ = "Bad Request"
		res.Details = []string{err.Error()}
		return nil, res
	}

	if utils.CheckPassword(foundUser.HashedPassword, reqDTO.Body.Password) != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Message = "Password not matching"
		res.Error_ = "Bad Request"
		return nil, res
	}

	token, err := utils.GenerateToken(foundUser)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Message = "Generate token failed"
		res.Error_ = "Internal Server Error"
		res.Details = []string{err.Error()}
		return nil, res
	}

	expireDuration, err := config.AppConfig.GetTokenExpireMinutes()
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Message = "Cannot parse expire duration for token"
		res.Error_ = "Internal Server Error"
		res.Details = []string{err.Error()}
		return nil, res
	}

	redisKey := fmt.Sprintf("token:%s", token)
	userData := map[string]interface{}{
		"user_id":   foundUser.Id,
		"role_name": foundUser.RoleName,
	}
	userDataBytes, _ := json.Marshal(userData)
	userHandler.redisClient.SetEx(ctx, redisKey, userDataBytes, *expireDuration)

	res := &dto.SuccessResponse[string]{}
	res.Body.Status = http.StatusOK
	res.Body.Message = "Login successful"
	res.Body.Data = token
	return res, nil
}

func (userHandler *UserHandler) Register(ctx context.Context, reqDTO *dto.RegisterRequest) (*dto.SuccessResponse[any], error) {
	createUserReqDTO := &dto.CreateUserRequest{}
	createUserReqDTO.Body.FullName = reqDTO.Body.FullName
	createUserReqDTO.Body.Email = reqDTO.Body.Email
	createUserReqDTO.Body.Username = reqDTO.Body.Username
	createUserReqDTO.Body.Password = reqDTO.Body.Password
	createUserReqDTO.Body.Address = reqDTO.Body.Address
	createUserReqDTO.Body.RoleName = "CUSTOMER"

	if err := userHandler.userService.CreateUser(ctx, createUserReqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Message = "Register failed"
		res.Error_ = "Bad Request"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse[any]{}
	res.Body.Status = http.StatusOK
	res.Body.Message = "Register successful"
	return res, nil
}

func (userHandler *UserHandler) GetAccountInfo(ctx context.Context, reqDTO *struct{}) (*dto.SuccessResponse[*dto.UserDTO], error) {
	id := ctx.Value("user_id")

	foundUser, err := userHandler.userService.GetUserById(ctx, id.(int64))
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Message = "Get account info failed"
		res.Error_ = "Bad Request"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToUserDTO(foundUser)
	res := &dto.SuccessResponse[*dto.UserDTO]{}
	res.Body.Status = http.StatusOK
	res.Body.Message = "Get account info successful"
	res.Body.Data = data
	return res, nil
}

func (userHandler *UserHandler) UpdateAccountInfo(ctx context.Context, reqDTO *dto.UpdateAccountInfoRequest) (*dto.SuccessResponse[any], error) {
	id := ctx.Value("user_id")

	updateUserReqDTO := &dto.UpdateUserRequest{}
	updateUserReqDTO.Body.FullName = reqDTO.Body.FullName
	updateUserReqDTO.Body.Email = reqDTO.Body.Email
	updateUserReqDTO.Body.Password = reqDTO.Body.Password
	updateUserReqDTO.Body.Address = reqDTO.Body.Address

	if err := userHandler.userService.UpdateUserById(ctx, id.(int64), updateUserReqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Message = "Update account info failed"
		res.Error_ = "Bad Request"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse[any]{}
	res.Body.Status = http.StatusOK
	res.Body.Message = "Update account info successful"
	return res, nil
}
