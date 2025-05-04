package handler

import (
	"context"
	"net/http"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/middleware"
	"thanhldt060802/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type UserHandler struct {
	userService    service.UserService
	authMiddleware *middleware.AuthMiddleware
}

func NewUserHandler(api huma.API, userService service.UserService, authMiddleware *middleware.AuthMiddleware) *UserHandler {
	userHandler := &UserHandler{
		userService:    userService,
		authMiddleware: authMiddleware,
	}

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

	// Get user by username
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/users/username/{username}",
		Summary:     "/users/username/{username}",
		Description: "Get user by username.",
		Tags:        []string{"User"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, userHandler.GetUserById)

	// Get user by email
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/users/email/{email}",
		Summary:     "/users/email/{email}",
		Description: "Get user by email.",
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

	// Login user
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/login",
		Summary:     "/login",
		Description: "Login user",
		Tags:        []string{"Auth"},
	}, userHandler.LoginUser)

	// Register user
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/register",
		Summary:     "/register",
		Description: "Register user.",
		Tags:        []string{"User"},
	}, userHandler.Register)

	// Get user using account
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/my-account",
		Summary:     "/my-account",
		Description: "Get user using account",
		Tags:        []string{"User"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication},
	}, userHandler.GetUserUsingAccount)

	// Update user using account
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/my-account",
		Summary:     "/my-account",
		Description: "Update user using account",
		Tags:        []string{"User"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication},
	}, userHandler.UpdateUserUsingAccount)

	return userHandler
}

func (userHandler *UserHandler) GetUsers(ctx context.Context, reqDTO *dto.GetUsersWithQueryParamRequest) (*dto.PaginationBodyResponseList[dto.UserView], error) {
	users, err := userHandler.userService.GetUsers(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get users failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToListUserView(users)
	res := &dto.PaginationBodyResponseList[dto.UserView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get users successful"
	res.Body.Data = data
	res.Body.Total = len(data)
	return res, nil
}

func (userHandler *UserHandler) GetUserById(ctx context.Context, reqDTO *dto.GetUserByIdRequest) (*dto.BodyResponse[dto.UserView], error) {
	foundUser, err := userHandler.userService.GetUserById(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get user by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToUserView(foundUser)
	res := &dto.BodyResponse[dto.UserView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get user by id successful"
	res.Body.Data = *data
	return res, nil
}

func (userHandler *UserHandler) GetUserByUsername(ctx context.Context, reqDTO *dto.GetUserByUsernameRequest) (*dto.BodyResponse[dto.UserView], error) {
	foundUser, err := userHandler.userService.GetUserByUsername(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get user by username failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToUserView(foundUser)
	res := &dto.BodyResponse[dto.UserView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get user by username successful"
	res.Body.Data = *data
	return res, nil
}

func (userHandler *UserHandler) GetUserByEmail(ctx context.Context, reqDTO *dto.GetUserByEmailRequest) (*dto.BodyResponse[dto.UserView], error) {
	foundUser, err := userHandler.userService.GetUserByEmail(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get user by email failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToUserView(foundUser)
	res := &dto.BodyResponse[dto.UserView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get user by email successful"
	res.Body.Data = *data
	return res, nil
}

func (userHandler *UserHandler) CreateUser(ctx context.Context, reqDTO *dto.CreateUserRequest) (*dto.SuccessResponse, error) {
	if err := userHandler.userService.CreateUser(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Create user failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Create user successful"
	return res, nil
}

func (userHandler *UserHandler) UpdateUserById(ctx context.Context, reqDTO *dto.UpdateUserRequest) (*dto.SuccessResponse, error) {
	if err := userHandler.userService.UpdateUserById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update user failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Update user successful"
	return res, nil
}

func (userHandler *UserHandler) DeleteUserById(ctx context.Context, reqDTO *dto.DeleteUserRequest) (*dto.SuccessResponse, error) {
	if err := userHandler.userService.DeleteUserById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete user failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Delete user successful"
	return res, nil
}

func (userHandler *UserHandler) LoginUser(ctx context.Context, reqDTO *dto.LoginRequest) (*dto.BodyResponse[string], error) {
	token, err := userHandler.userService.LoginUser(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Login user failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.BodyResponse[string]{}
	res.Body.Code = "OK"
	res.Body.Message = "Login user successful"
	res.Body.Data = *token
	return res, nil
}

func (userHandler *UserHandler) Register(ctx context.Context, reqDTO *dto.RegisterRequest) (*dto.SuccessResponse, error) {
	convertReqDTO := &dto.CreateUserRequest{}
	convertReqDTO.Body.FullName = reqDTO.Body.FullName
	convertReqDTO.Body.Email = reqDTO.Body.Email
	convertReqDTO.Body.Username = reqDTO.Body.Username
	convertReqDTO.Body.Password = reqDTO.Body.Password
	convertReqDTO.Body.Address = reqDTO.Body.Address
	convertReqDTO.Body.RoleName = "CUSTOMER"

	if err := userHandler.userService.CreateUser(ctx, convertReqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Register user failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Register user successful"
	return res, nil
}

func (userHandler *UserHandler) GetUserUsingAccount(ctx context.Context, reqDTO *struct{}) (*dto.BodyResponse[dto.UserView], error) {
	userId := ctx.Value("user_id").(int64)

	convertReqDTO := &dto.GetUserByIdRequest{Id: userId}

	foundUser, err := userHandler.userService.GetUserById(ctx, convertReqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get user using account failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToUserView(foundUser)
	res := &dto.BodyResponse[dto.UserView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get user using account successful"
	res.Body.Data = *data
	return res, nil
}

func (userHandler *UserHandler) UpdateUserUsingAccount(ctx context.Context, reqDTO *dto.UpdateUserUsingAccountRequest) (*dto.SuccessResponse, error) {
	userId := ctx.Value("user_id").(int64)

	convertReqDTO := &dto.UpdateUserRequest{Id: userId}
	convertReqDTO.Body.FullName = reqDTO.Body.FullName
	convertReqDTO.Body.Email = reqDTO.Body.Password
	convertReqDTO.Body.Password = reqDTO.Body.Password
	convertReqDTO.Body.Address = reqDTO.Body.Address

	if err := userHandler.userService.UpdateUserById(ctx, convertReqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update account info failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Update account info successful"
	return res, nil
}
