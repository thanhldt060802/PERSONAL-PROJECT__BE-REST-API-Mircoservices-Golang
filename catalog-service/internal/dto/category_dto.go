package dto

import (
	"thanhldt060802/internal/model"
	"time"
)

type CategoryView struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToCategoryView(category *model.Category) *CategoryView {
	return &CategoryView{
		Id:          category.Id,
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

func ToListCategoryView(categorys []model.Category) []CategoryView {
	categoryViews := make([]CategoryView, len(categorys))
	for i, category := range categorys {
		categoryViews[i] = *ToCategoryView(&category)
	}
	return categoryViews
}
