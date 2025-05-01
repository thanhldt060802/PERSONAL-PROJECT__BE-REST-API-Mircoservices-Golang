package repository

import (
	"context"
	"fmt"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/model"

	"github.com/uptrace/bun"
)

type userRepository struct {
	db *bun.DB
}

type UserRepository interface {
	ExistsById(ctx context.Context, id int64) (bool, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Get(ctx context.Context, offset int, limit int, sortFields []dto.SortField) ([]model.User, error)
	GetById(ctx context.Context, id int64) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, newUser *model.User) error
	UpdateById(ctx context.Context, id int64, updatedUser *model.User) error
	DeleteById(ctx context.Context, id int64) error
}

func NewUserRepository(db *bun.DB) UserRepository {
	return &userRepository{db: db}
}

func (userRepository *userRepository) ExistsById(ctx context.Context, id int64) (bool, error) {
	count, err := userRepository.db.NewSelect().Model(&model.User{}).Where("id = ?", id).Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (userRepository *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	count, err := userRepository.db.NewSelect().Model(&model.User{}).Where("username = ?", username).Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (userRepository *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	count, err := userRepository.db.NewSelect().Model(&model.User{}).Where("email = ?", email).Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (userRepository *userRepository) Get(ctx context.Context, offser int, limit int, sortFields []dto.SortField) ([]model.User, error) {
	var users []model.User
	query := userRepository.db.NewSelect().Model(&users).
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
	err := userRepository.db.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepository *userRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := userRepository.db.NewSelect().Model(&user).Where("username = ?", username).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepository *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := userRepository.db.NewSelect().Model(&user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepository *userRepository) Create(ctx context.Context, newUser *model.User) error {
	_, err := userRepository.db.NewInsert().Model(newUser).Exec(ctx)
	return err
}

func (userRepository *userRepository) UpdateById(ctx context.Context, id int64, updatedUser *model.User) error {
	_, err := userRepository.db.NewUpdate().Model(updatedUser).Where("id = ?", id).Exec(ctx)
	return err
}

func (userRepository *userRepository) DeleteById(ctx context.Context, id int64) error {
	_, err := userRepository.db.NewDelete().Model(&model.User{}).Where("id = ?", id).Exec(ctx)
	return err
}
