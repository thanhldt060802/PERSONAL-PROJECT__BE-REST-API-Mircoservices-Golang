package repository

import (
	"context"
	"fmt"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
	"thanhldt060802/utils"
)

type cartRepository struct {
}

type CartRepository interface {
	Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField) ([]model.Cart, error)
	GetById(ctx context.Context, id int64) (*model.Cart, error)
	GetByUserId(ctx context.Context, userId int64) (*model.Cart, error)
	Create(ctx context.Context, newCart *model.Cart) error
	UpdateById(ctx context.Context, id int64, updatedCart *model.Cart) error
	DeleteById(ctx context.Context, id int64) error
}

func NewCartRepository() CartRepository {
	return &cartRepository{}
}

func (cartRepository *cartRepository) Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField) ([]model.Cart, error) {
	var carts []model.Cart
	query := infrastructure.DB.NewSelect().Model(&carts).
		Offset(offset).
		Limit(limit)
	for _, sortField := range sortFields {
		query = query.Order(fmt.Sprintf("%s %s", sortField.Field, sortField.Direction))
	}
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return carts, nil
}

func (cartRepository *cartRepository) GetById(ctx context.Context, id int64) (*model.Cart, error) {
	var cart model.Cart
	err := infrastructure.DB.NewSelect().Model(&cart).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (cartRepository *cartRepository) GetByUserId(ctx context.Context, userId int64) (*model.Cart, error) {
	var cart model.Cart
	err := infrastructure.DB.NewSelect().Model(&cart).Where("user_id = ?", userId).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (cartRepository *cartRepository) Create(ctx context.Context, newCart *model.Cart) error {
	_, err := infrastructure.DB.NewInsert().Model(newCart).Exec(ctx)
	return err
}

func (cartRepository *cartRepository) UpdateById(ctx context.Context, id int64, updatedCart *model.Cart) error {
	_, err := infrastructure.DB.NewUpdate().Model(updatedCart).Where("id = ?", id).Exec(ctx)
	return err
}

func (cartRepository *cartRepository) DeleteById(ctx context.Context, id int64) error {
	_, err := infrastructure.DB.NewDelete().Model(&model.Cart{}).Where("id = ?", id).Exec(ctx)
	return err
}
