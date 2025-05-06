package service

import (
	"context"
	"fmt"
	"strings"
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
	GetCategories(ctx context.Context, reqDTO *dto.GetCategoriesRequest) ([]model.Category, error)
	GetCategoryById(ctx context.Context, reqDTO *dto.GetCategoryByIdRequest) (*model.Category, error)
	GetCategoryByName(ctx context.Context, reqDTO *dto.GetCategoryByNameRequest) (*model.Category, error)
	CreateCategory(ctx context.Context, reqDTO *dto.CreateCategoryRequest) error
	UpdateCategoryById(ctx context.Context, reqDTO *dto.UpdateCategoryByIdRequest) error
	DeleteCategoryById(ctx context.Context, reqDTO *dto.DeleteCategoryByIdRequest) error
}

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepository: categoryRepository}
}

func (categoryService *categoryService) GetCategories(ctx context.Context, reqDTO *dto.GetCategoriesRequest) ([]model.Category, error) {
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

func (categoryService *categoryService) GetCategoryByName(ctx context.Context, reqDTO *dto.GetCategoryByNameRequest) (*model.Category, error) {
	foundCategory, err := categoryService.categoryRepository.GetByName(ctx, reqDTO.Name)
	if err != nil {
		return nil, err
	}

	return foundCategory, nil
}

func (categoryService *categoryService) CreateCategory(ctx context.Context, reqDTO *dto.CreateCategoryRequest) error {
	if _, err := categoryService.categoryRepository.GetByName(ctx, reqDTO.Body.Name); err != nil {
		return fmt.Errorf("name of category already exists")
	}

	newCategory := model.Category{
		Name: reqDTO.Body.Name,
	}
	if err := categoryService.categoryRepository.Create(ctx, &newCategory); err != nil {
		return err
	}

	return nil
}

func (categoryService *categoryService) UpdateCategoryById(ctx context.Context, reqDTO *dto.UpdateCategoryByIdRequest) error {
	foundCategory, err := categoryService.categoryRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return fmt.Errorf("id of category not found")
	}

	if reqDTO.Body.Name != nil && strings.EqualFold(foundCategory.Name, *reqDTO.Body.Name) {
		if _, err := categoryService.categoryRepository.GetByName(ctx, *reqDTO.Body.Name); err == nil {
			return fmt.Errorf("name of category already exists")
		}
		foundCategory.Name = *reqDTO.Body.Name
	}
	foundCategory.UpdatedAt = time.Now().UTC()

	if err := categoryService.categoryRepository.Update(ctx, foundCategory); err != nil {
		return err
	}

	return nil
}

func (categoryService *categoryService) DeleteCategoryById(ctx context.Context, reqDTO *dto.DeleteCategoryByIdRequest) error {
	if _, err := categoryService.categoryRepository.GetById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("id of category not found")
	}

	if err := categoryService.categoryRepository.DeleteById(ctx, reqDTO.Id); err != nil {
		return err
	}

	return nil
}
