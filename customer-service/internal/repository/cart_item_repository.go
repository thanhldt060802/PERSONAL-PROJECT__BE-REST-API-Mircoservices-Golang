package repository

import (
	"context"
	"fmt"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
	"thanhldt060802/utils"
)

type cartItemRepository struct {
}

type CartItemRepository interface {
	Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField) ([]model.CartItem, error)
	GetById(ctx context.Context, id int64) (*model.CartItem, error)
	GetByCartId(ctx context.Context, cartId int64, offset int, limit int, sortFields []utils.SortField) ([]model.CartItem, error)
	Create(ctx context.Context, newCartItem *model.CartItem) error
	UpdateById(ctx context.Context, id int64, updatedCartItem *model.CartItem) error
	DeleteById(ctx context.Context, id int64) error
}

func NewCartItemRepository() CartItemRepository {
	return &cartItemRepository{}
}

func (cartItemRepository *cartItemRepository) Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField) ([]model.CartItem, error) {
	var cartItems []model.CartItem
	query := infrastructure.DB.NewSelect().Model(&cartItems).
		Offset(offset).
		Limit(limit)
	for _, sortField := range sortFields {
		query = query.Order(fmt.Sprintf("%s %s", sortField.Field, sortField.Direction))
	}
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return cartItems, nil
}

func (cartItemRepository *cartItemRepository) GetById(ctx context.Context, id int64) (*model.CartItem, error) {
	var cartItem model.CartItem
	err := infrastructure.DB.NewSelect().Model(&cartItem).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func (cartItemRepository *cartItemRepository) GetByCartId(ctx context.Context, cartId int64, offset int, limit int, sortFields []utils.SortField) ([]model.CartItem, error) {
	var cartItems []model.CartItem
	query := infrastructure.DB.NewSelect().Model(&cartItems).Where("cart_id = ?", cartId).
		Offset(offset).
		Limit(limit)
	for _, sortField := range sortFields {
		query = query.Order(fmt.Sprintf("%s %s", sortField.Field, sortField.Direction))
	}
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return cartItems, nil
}

func (cartItemRepository *cartItemRepository) Create(ctx context.Context, newCartItem *model.CartItem) error {
	_, err := infrastructure.DB.NewInsert().Model(newCartItem).Exec(ctx)
	return err
}

func (cartItemRepository *cartItemRepository) UpdateById(ctx context.Context, id int64, updatedCartItem *model.CartItem) error {
	_, err := infrastructure.DB.NewUpdate().Model(updatedCartItem).Where("id = ?", id).Exec(ctx)
	return err
}

func (cartItemRepository *cartItemRepository) DeleteById(ctx context.Context, id int64) error {
	_, err := infrastructure.DB.NewDelete().Model(&model.CartItem{}).Where("id = ?", id).Exec(ctx)
	return err
}
