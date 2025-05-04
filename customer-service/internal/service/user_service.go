package service

import (
	"context"
	"encoding/json"
	"fmt"
	"thanhldt060802/config"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/repository"
	"thanhldt060802/utils"
	"time"
)

type userService struct {
	userRepository repository.UserRepository
	cartRepository repository.CartRepository
}

type UserService interface {
	GetUsers(ctx context.Context, reqDTO *dto.GetUsersWithQueryParamRequest) ([]model.User, error)
	GetUserById(ctx context.Context, reqDTO *dto.GetUserByIdRequest) (*model.User, error)
	GetUserByUsername(ctx context.Context, reqDTO *dto.GetUserByUsernameRequest) (*model.User, error)
	GetUserByEmail(ctx context.Context, reqDTO *dto.GetUserByEmailRequest) (*model.User, error)
	CreateUser(ctx context.Context, reqDTO *dto.CreateUserRequest) error
	UpdateUserById(ctx context.Context, reqDTO *dto.UpdateUserRequest) error
	DeleteUserById(ctx context.Context, reqDTO *dto.DeleteUserRequest) error

	LoginUser(ctx context.Context, reqDTO *dto.LoginRequest) (*string, error)
}

func NewUserService(userRepository repository.UserRepository, cartRepository repository.CartRepository) UserService {
	return &userService{
		userRepository: userRepository,
		cartRepository: cartRepository,
	}
}

func (userService *userService) GetUsers(ctx context.Context, reqDTO *dto.GetUsersWithQueryParamRequest) ([]model.User, error) {
	sortFields := utils.ParseSortBy(reqDTO.SortBy)

	users, err := userService.userRepository.Get(ctx, reqDTO.Offset, reqDTO.Limit, sortFields)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (userService *userService) GetUserById(ctx context.Context, reqDTO *dto.GetUserByIdRequest) (*model.User, error) {
	foundUser, err := userService.userRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return nil, err
	}

	return foundUser, nil
}

func (userService *userService) GetUserByUsername(ctx context.Context, reqDTO *dto.GetUserByUsernameRequest) (*model.User, error) {
	foundUser, err := userService.userRepository.GetByUsername(ctx, reqDTO.Username)
	if err != nil {
		return nil, err
	}

	return foundUser, nil
}

func (userService *userService) GetUserByEmail(ctx context.Context, reqDTO *dto.GetUserByEmailRequest) (*model.User, error) {
	foundUser, err := userService.userRepository.GetByEmail(ctx, reqDTO.Email)
	if err != nil {
		return nil, err
	}

	return foundUser, nil
}

func (userService *userService) CreateUser(ctx context.Context, reqDTO *dto.CreateUserRequest) error {
	if _, err := userService.userRepository.GetByUsername(ctx, reqDTO.Body.Username); err == nil {
		return fmt.Errorf("username of user is already exists")
	}
	if _, err := userService.userRepository.GetByEmail(ctx, reqDTO.Body.Email); err == nil {
		return fmt.Errorf("email of user is already exists")
	}

	hashedPassword, err := utils.HashPassword(reqDTO.Body.Password)
	if err != nil {
		return fmt.Errorf("hash password failed")
	}

	newUser := model.User{
		FullName:       reqDTO.Body.FullName,
		Email:          reqDTO.Body.Email,
		Username:       reqDTO.Body.Username,
		HashedPassword: hashedPassword,
		Address:        reqDTO.Body.Address,
		RoleName:       reqDTO.Body.RoleName,
	}
	if err := userService.userRepository.Create(ctx, &newUser); err != nil {
		return err
	}

	newCart := model.Cart{
		UserId: newUser.Id,
	}
	if err := userService.cartRepository.Create(ctx, &newCart); err != nil {
		return err
	}

	return nil
}

func (userService *userService) UpdateUserById(ctx context.Context, reqDTO *dto.UpdateUserRequest) error {
	foundUser, err := userService.userRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return fmt.Errorf("id of user is not valid")
	}

	if reqDTO.Body.FullName != nil {
		foundUser.FullName = *reqDTO.Body.FullName
	}
	if reqDTO.Body.Email != nil && reqDTO.Body.Email != &foundUser.Email {
		if _, err = userService.userRepository.GetByEmail(ctx, *reqDTO.Body.Email); err == nil {
			return fmt.Errorf("email of user is already exists")
		}
		foundUser.Email = *reqDTO.Body.Email
	}
	if reqDTO.Body.Password != nil {
		hashedPassword, err := utils.HashPassword(*reqDTO.Body.Password)
		if err != nil {
			return fmt.Errorf("hash password failed")
		}
		foundUser.HashedPassword = hashedPassword
	}
	if reqDTO.Body.Address != nil {
		foundUser.Address = *reqDTO.Body.Address
	}
	if reqDTO.Body.RoleName != nil {
		foundUser.RoleName = *reqDTO.Body.RoleName
	}
	foundUser.UpdatedAt = time.Now().UTC()

	if err := userService.userRepository.UpdateById(ctx, reqDTO.Id, foundUser); err != nil {
		return err
	}

	return nil
}

func (userService *userService) DeleteUserById(ctx context.Context, reqDTO *dto.DeleteUserRequest) error {
	foundCart, err := userService.cartRepository.GetByUserId(ctx, reqDTO.Id)
	if err != nil {
		return fmt.Errorf("user id of cart is not valid")
	}

	if _, err := userService.userRepository.GetById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("id of user is not valid")
	}

	if err := userService.cartRepository.DeleteById(ctx, foundCart.Id); err != nil {
		return err
	}

	if err := userService.userRepository.DeleteById(ctx, reqDTO.Id); err != nil {
		return err
	}

	return nil
}

func (userService *userService) LoginUser(ctx context.Context, reqDTO *dto.LoginRequest) (*string, error) {
	foundUser, err := userService.userRepository.GetByUsername(ctx, reqDTO.Body.Username)
	if err != nil {
		return nil, err
	} else if utils.CheckPassword(foundUser.HashedPassword, reqDTO.Body.Password) != nil {
		return nil, fmt.Errorf("password does not match")
	}

	foundCart, err := userService.cartRepository.GetByUserId(ctx, foundUser.Id)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(foundUser.Id, foundUser.RoleName, foundCart.Id)
	if err != nil {
		return nil, fmt.Errorf("generate token failed")
	}

	expireDuration := config.AppConfig.GetTokenExpireMinutes()
	if expireDuration == nil {
		return nil, fmt.Errorf("convert expire failed")
	}

	redisKey := fmt.Sprintf("token:%s", *token)
	userData := map[string]interface{}{
		"user_id":   foundUser.Id,
		"role_name": foundUser.RoleName,
		"cart_id":   foundCart.Id,
	}
	userDataBytes, _ := json.Marshal(userData)
	infrastructure.RedisClient.SetEx(ctx, redisKey, userDataBytes, *expireDuration)

	return token, nil
}
