package repository

import (
	"context"
	"fmt"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
	"thanhldt060802/utils"
)

type productRepository struct {
}

type ProductRepository interface {
	Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField) ([]model.Product, error)
	GetById(ctx context.Context, id int64) (*model.Product, error)
	GetByCategoryId(ctx context.Context, categoryId int64, offset int, limit int, sortFields []utils.SortField) ([]model.Product, error)
	Create(ctx context.Context, newProduct *model.Product) error
	UpdateById(ctx context.Context, id int64, updatedProduct *model.Product) error
	DeleteById(ctx context.Context, id int64) error

	GetAll(ctx context.Context) ([]model.Product, error)
}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

func (productRepository *productRepository) Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField) ([]model.Product, error) {
	var products []model.Product
	query := infrastructure.DB.NewSelect().Model(&products).
		Offset(offset).
		Limit(limit)
	for _, sortField := range sortFields {
		query = query.Order(fmt.Sprintf("%s %s", sortField.Field, sortField.Direction))
	}
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (productRepository *productRepository) GetById(ctx context.Context, id int64) (*model.Product, error) {
	var product model.Product
	err := infrastructure.DB.NewSelect().Model(&product).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (productRepository *productRepository) GetByCategoryId(ctx context.Context, categoryId int64, offset int, limit int, sortFields []utils.SortField) ([]model.Product, error) {
	var products []model.Product
	query := infrastructure.DB.NewSelect().Model(&products).Where("category_id = ?", categoryId).
		Offset(offset).
		Limit(limit)
	for _, sortField := range sortFields {
		query = query.Order(fmt.Sprintf("%s %s", sortField.Field, sortField.Direction))
	}
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (productRepository *productRepository) Create(ctx context.Context, newProduct *model.Product) error {
	_, err := infrastructure.DB.NewInsert().Model(newProduct).Exec(ctx)
	return err
}

func (productRepository *productRepository) UpdateById(ctx context.Context, id int64, updatedProduct *model.Product) error {
	_, err := infrastructure.DB.NewUpdate().Model(updatedProduct).Where("id = ?", id).Exec(ctx)
	return err
}

func (productRepository *productRepository) DeleteById(ctx context.Context, id int64) error {
	_, err := infrastructure.DB.NewDelete().Model(&model.Product{}).Where("id = ?", id).Exec(ctx)
	return err
}

func (productRepository *productRepository) GetAll(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	err := infrastructure.DB.NewSelect().Model(&products).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}
