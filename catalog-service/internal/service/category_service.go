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

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

type CategoryService interface {
	GetCategories(ctx context.Context, reqDTO *dto.GetCategoriesWithQueryParamRequest) ([]model.Category, error)
	GetCategoryById(ctx context.Context, reqDTO *dto.GetCategoryByIdRequest) (*model.Category, error)
	CreateCategory(ctx context.Context, reqDTO *dto.CreateCategoryRequest) error
	UpdateCategoryById(ctx context.Context, reqDTO *dto.UpdateCategoryRequest) error
	DeleteCategoryById(ctx context.Context, reqDTO *dto.DeleteCategoryRequest) error
}

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepository: categoryRepository}
}

func (categoryService *categoryService) GetCategories(ctx context.Context, reqDTO *dto.GetCategoriesWithQueryParamRequest) ([]model.Category, error) {
	sortFields := utils.ParseSortBy(reqDTO.SortBy)

	categories, err := categoryService.categoryRepository.Get(ctx, reqDTO.Offset, reqDTO.Limit, sortFields)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (categoryService *categoryService) GetCategoryById(ctx context.Context, reqDTO *dto.GetCategoryByIdRequest) (*model.Category, error) {
	foundCategory, err := categoryService.categoryRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return nil, err
	}

	return foundCategory, nil
}

func (categoryService *categoryService) CreateCategory(ctx context.Context, reqDTO *dto.CreateCategoryRequest) error {
	newCategory := model.Category{
		Name:        reqDTO.Body.Name,
		Description: reqDTO.Body.Description,
	}
	if err := categoryService.categoryRepository.Create(ctx, &newCategory); err != nil {
		return err
	}

	return nil
}

func (categoryService *categoryService) UpdateCategoryById(ctx context.Context, reqDTO *dto.UpdateCategoryRequest) error {
	foundCategory, err := categoryService.categoryRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return fmt.Errorf("id of category is not valid")
	}

	if reqDTO.Body.Name != nil {
		if _, err := categoryService.categoryRepository.GetByName(ctx, *reqDTO.Body.Name); err == nil {
			return fmt.Errorf("name of category is already exists")
		}
		foundCategory.Name = *reqDTO.Body.Name
	}
	if reqDTO.Body.Description != nil {
		foundCategory.Description = *reqDTO.Body.Description
	}
	foundCategory.UpdatedAt = time.Now().UTC()

	if err := categoryService.categoryRepository.UpdateById(ctx, reqDTO.Id, foundCategory); err != nil {
		return err
	}

	return nil
}

func (categoryService *categoryService) DeleteCategoryById(ctx context.Context, reqDTO *dto.DeleteCategoryRequest) error {
	if _, err := categoryService.categoryRepository.GetById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("id of category is not valid")
	}

	if err := categoryService.categoryRepository.DeleteById(ctx, reqDTO.Id); err != nil {
		return err
	}

	return nil
}
