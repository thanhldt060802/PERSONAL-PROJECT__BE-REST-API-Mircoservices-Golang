package service

import (
	"context"
	"fmt"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/repository"
	"thanhldt060802/utils"
	"time"
)

type userService struct {
	userRepository repository.UserRepository
}

type UserService interface {
	IsUserExistedById(ctx context.Context, id int64) (bool, error)
	IsUserExistedByUsername(ctx context.Context, username string) (bool, error)
	IsUserExistedByEmail(ctx context.Context, email string) (bool, error)
	GetUsers(ctx context.Context, queryParam *dto.GetUsersRequestQueryParam) ([]model.User, error)
	GetUserById(ctx context.Context, id int64) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, reqDTO *dto.CreateUserRequest) error
	UpdateUserById(ctx context.Context, id int64, reqDTO *dto.UpdateUserRequest) error
	DeleteUserById(ctx context.Context, id int64) error
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

func (userService *userService) IsUserExistedById(ctx context.Context, id int64) (bool, error) {
	existed, err := userService.userRepository.ExistsById(ctx, id)
	if err != nil {
		return false, fmt.Errorf("get user by id failed: %w", err)
	}
	return existed, nil
}

func (userService *userService) IsUserExistedByUsername(ctx context.Context, username string) (bool, error) {
	existed, err := userService.userRepository.ExistsByUsername(ctx, username)
	if err != nil {
		return false, fmt.Errorf("get user by username failed: %w", err)
	}
	return existed, nil
}

func (userService *userService) IsUserExistedByEmail(ctx context.Context, email string) (bool, error) {
	existed, err := userService.userRepository.ExistsByEmail(ctx, email)
	if err != nil {
		return false, fmt.Errorf("get user by email failed: %w", err)
	}
	return existed, nil
}

func (userService *userService) GetUsers(ctx context.Context, queryParam *dto.GetUsersRequestQueryParam) ([]model.User, error) {
	sortFields := dto.ParseSortBy(queryParam.SortBy)

	users, err := userService.userRepository.Get(ctx, queryParam.Offset, queryParam.Limit, sortFields)
	if err != nil {
		return nil, fmt.Errorf("get users failed: %w", err)
	}

	return users, nil
}

func (userService *userService) GetUserById(ctx context.Context, id int64) (*model.User, error) {
	foundUser, err := userService.userRepository.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user by id failed: %w", err)
	}

	return foundUser, nil
}

func (userService *userService) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	foundUser, err := userService.userRepository.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get user by username failed: %w", err)
	}

	return foundUser, nil
}

func (userService *userService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	foundUser, err := userService.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get user by email failed: %w", err)
	}

	return foundUser, nil
}

func (userService *userService) CreateUser(ctx context.Context, reqDTO *dto.CreateUserRequest) error {
	existedByEmail, err := userService.IsUserExistedByEmail(ctx, reqDTO.Body.Email)
	if err != nil {
		return err
	}
	if existedByEmail {
		return fmt.Errorf("email of user is already exists")
	}

	existedByUsername, err := userService.IsUserExistedByUsername(ctx, reqDTO.Body.Username)
	if err != nil {
		return err
	}
	if existedByUsername {
		return fmt.Errorf("username of user is already exists")
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

	return userService.userRepository.Create(ctx, &newUser)
}

func (userService *userService) UpdateUserById(ctx context.Context, id int64, reqDTO *dto.UpdateUserRequest) error {
	foundUser, err := userService.GetUserById(ctx, id)
	if err != nil {
		return err
	}
	if foundUser == nil {
		return fmt.Errorf("id of user is not valid")
	}

	if reqDTO.Body.FullName != nil {
		foundUser.FullName = *reqDTO.Body.FullName
	}
	if reqDTO.Body.Email != nil {
		existedByEmail, err := userService.IsUserExistedByEmail(ctx, *reqDTO.Body.Email)
		if err != nil {
			return err
		}
		if existedByEmail {
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

	return userService.userRepository.UpdateById(ctx, id, foundUser)
}

func (userService *userService) DeleteUserById(ctx context.Context, id int64) error {
	existed, err := userService.IsUserExistedById(ctx, id)
	if err != nil {
		return err
	}
	if !existed {
		return fmt.Errorf("id of user is not valid")
	}

	return userService.userRepository.DeleteById(ctx, id)
}
