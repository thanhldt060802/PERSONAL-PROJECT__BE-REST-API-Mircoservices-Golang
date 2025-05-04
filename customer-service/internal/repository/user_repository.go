package repository

import (
	"context"
	"fmt"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
	"thanhldt060802/utils"
)

type userRepository struct {
}

type UserRepository interface {
	Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField) ([]model.User, error)
	GetById(ctx context.Context, id int64) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, newUser *model.User) error
	UpdateById(ctx context.Context, id int64, updatedUser *model.User) error
	DeleteById(ctx context.Context, id int64) error
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (userRepository *userRepository) Get(ctx context.Context, offser int, limit int, sortFields []utils.SortField) ([]model.User, error) {
	var users []model.User
	query := infrastructure.DB.NewSelect().Model(&users).
		Offset(offser).
		Limit(limit)
	for _, sortField := range sortFields {
		query = query.Order(fmt.Sprintf("%s %s", sortField.Field, sortField.Direction))
	}
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (userRepository *userRepository) GetById(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := infrastructure.DB.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepository *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := infrastructure.DB.NewSelect().Model(&user).Where("username = ?", username).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepository *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := infrastructure.DB.NewSelect().Model(&user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepository *userRepository) Create(ctx context.Context, newUser *model.User) error {
	_, err := infrastructure.DB.NewInsert().Model(newUser).Exec(ctx)
	return err
}

func (userRepository *userRepository) UpdateById(ctx context.Context, id int64, updatedUser *model.User) error {
	_, err := infrastructure.DB.NewUpdate().Model(updatedUser).Where("id = ?", id).Exec(ctx)
	return err
}

func (userRepository *userRepository) DeleteById(ctx context.Context, id int64) error {
	_, err := infrastructure.DB.NewDelete().Model(&model.User{}).Where("id = ?", id).Exec(ctx)
	return err
}
