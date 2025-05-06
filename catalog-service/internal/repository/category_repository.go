package repository

import (
	"context"
	"fmt"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
	"thanhldt060802/utils"
)

type categoryRepository struct {
}

type CategoryRepository interface {
	Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField) ([]model.Category, error)
	GetById(ctx context.Context, id int64) (*model.Category, error)
	GetByName(ctx context.Context, name string) (*model.Category, error)
	Create(ctx context.Context, newCategory *model.Category) error
	Update(ctx context.Context, updatedCategory *model.Category) error
	DeleteById(ctx context.Context, id int64) error
}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepository{}
}

func (categoryRepository *categoryRepository) Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField) ([]model.Category, error) {
	var categories []model.Category

	query := infrastructure.DB.NewSelect().Model(&categories).
		Offset(offset).
		Limit(limit)
	for _, sortField := range sortFields {
		query = query.Order(fmt.Sprintf("%s %s", sortField.Field, sortField.Direction))
	}
	err := query.Scan(ctx)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (categoryRepository *categoryRepository) GetById(ctx context.Context, id int64) (*model.Category, error) {
	var category model.Category

	err := infrastructure.DB.NewSelect().Model(&category).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (categoryRepository *categoryRepository) GetByName(ctx context.Context, name string) (*model.Category, error) {
	var category model.Category

	err := infrastructure.DB.NewSelect().Model(&category).Where("LOWER(name) = LOWER(?)", name).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (categoryRepository *categoryRepository) Create(ctx context.Context, newCategory *model.Category) error {
	_, err := infrastructure.DB.NewInsert().Model(newCategory).Exec(ctx)

	return err
}

func (categoryRepository *categoryRepository) Update(ctx context.Context, updatedCategory *model.Category) error {
	_, err := infrastructure.DB.NewUpdate().Model(updatedCategory).Where("id = ?", updatedCategory.Id).Exec(ctx)

	return err
}

func (categoryRepository *categoryRepository) DeleteById(ctx context.Context, id int64) error {
	_, err := infrastructure.DB.NewDelete().Model(&model.Category{}).Where("id = ?", id).Exec(ctx)

	return err
}
