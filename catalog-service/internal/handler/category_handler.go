package handler

import (
	"context"
	"net/http"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/middleware"
	"thanhldt060802/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type CategoryHandler struct {
	categorieservice service.CategoryService
	authMiddleware   *middleware.AuthMiddleware
}

func NewCategoryHandler(api huma.API, categorieservice service.CategoryService, authMiddleware *middleware.AuthMiddleware) *CategoryHandler {
	categoryHandler := &CategoryHandler{
		categorieservice: categorieservice,
		authMiddleware:   authMiddleware,
	}

	// Get categories
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/categories",
		Summary:     "/categories",
		Description: "Get categories.",
		Tags:        []string{"Category"},
	}, categoryHandler.GetCategories)

	// Get category by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/categories/id/{id}",
		Summary:     "/categories/id/{id}",
		Description: "Get category by id.",
		Tags:        []string{"Category"},
	}, categoryHandler.GetCategoryById)

	// Get category by name
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/categories/name/{name}",
		Summary:     "/categories/name/{name}",
		Description: "Get category by name.",
		Tags:        []string{"Category"},
	}, categoryHandler.GetCategoryByName)

	// Create category
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/categories",
		Summary:     "/categories",
		Description: "Create category.",
		Tags:        []string{"Category"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, categoryHandler.CreateCategory)

	// Update category by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/categories/id/{id}",
		Summary:     "/categories/id/{id}",
		Description: "Update category by id.",
		Tags:        []string{"Category"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, categoryHandler.UpdateCategoryById)

	// Delete category by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/categories/id/{id}",
		Summary:     "/categories/id/{id}",
		Description: "Delete category by id.",
		Tags:        []string{"Category"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, categoryHandler.DeleteCategoryById)

	return categoryHandler
}

func (categoryHandler *CategoryHandler) GetCategories(ctx context.Context, reqDTO *dto.GetCategoriesRequest) (*dto.PaginationBodyResponseList[dto.CategoryView], error) {
	categories, err := categoryHandler.categorieservice.GetCategories(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get categories failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToListCategoryView(categories)
	res := &dto.PaginationBodyResponseList[dto.CategoryView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get categories successful"
	res.Body.Data = data
	res.Body.Total = len(data)
	return res, nil
}

func (categoryHandler *CategoryHandler) GetCategoryById(ctx context.Context, reqDTO *dto.GetCategoryByIdRequest) (*dto.BodyResponse[dto.CategoryView], error) {
	foundCategory, err := categoryHandler.categorieservice.GetCategoryById(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get category by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToCategoryView(foundCategory)
	res := &dto.BodyResponse[dto.CategoryView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get category by id successful"
	res.Body.Data = *data
	return res, nil
}

func (categoryHandler *CategoryHandler) GetCategoryByName(ctx context.Context, reqDTO *dto.GetCategoryByNameRequest) (*dto.BodyResponse[dto.CategoryView], error) {
	foundCategory, err := categoryHandler.categorieservice.GetCategoryByName(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get category by name failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToCategoryView(foundCategory)
	res := &dto.BodyResponse[dto.CategoryView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get category by name successful"
	res.Body.Data = *data
	return res, nil
}

func (categoryHandler *CategoryHandler) CreateCategory(ctx context.Context, reqDTO *dto.CreateCategoryRequest) (*dto.SuccessResponse, error) {
	if err := categoryHandler.categorieservice.CreateCategory(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Create category failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Create category successful"
	return res, nil
}

func (categoryHandler *CategoryHandler) UpdateCategoryById(ctx context.Context, reqDTO *dto.UpdateCategoryByIdRequest) (*dto.SuccessResponse, error) {
	if err := categoryHandler.categorieservice.UpdateCategoryById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update category failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Update category successful"
	return res, nil
}

func (categoryHandler *CategoryHandler) DeleteCategoryById(ctx context.Context, reqDTO *dto.DeleteCategoryByIdRequest) (*dto.SuccessResponse, error) {
	if err := categoryHandler.categorieservice.DeleteCategoryById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete category failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Delete category successful"
	return res, nil
}
