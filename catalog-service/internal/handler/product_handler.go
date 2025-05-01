package handler

import (
	"context"
	"net/http"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/middleware"
	"thanhldt060802/internal/service"

	"github.com/danielgtaylor/huma/v2"
)

type ProductHandler struct {
	productService service.ProductService
	authMiddleware *middleware.AuthMiddleware
}

func NewProductHandler(api huma.API, productService service.ProductService, authMiddleware *middleware.AuthMiddleware) *ProductHandler {
	productHandler := &ProductHandler{
		productService: productService,
		authMiddleware: authMiddleware,
	}

	// Get products
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/products",
		Summary:     "/products",
		Description: "Get products.",
		Tags:        []string{"Product"},
	}, productHandler.GetProducts)

	// Get product by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/products/id/{id}",
		Summary:     "/products/id/{id}",
		Description: "Get product by id.",
		Tags:        []string{"Product"},
	}, productHandler.GetProductById)

	// Create product
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/products",
		Summary:     "/products",
		Description: "Create product.",
		Tags:        []string{"Product"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, productHandler.CreateProduct)

	// Update product by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/products/id/{id}",
		Summary:     "/products/id/{id}",
		Description: "Update product by id.",
		Tags:        []string{"Product"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, productHandler.UpdateProductById)

	// Delete product by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/products/id{id}",
		Summary:     "/products/id{id}",
		Description: "Delete product by id.",
		Tags:        []string{"Product"},
		Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, productHandler.DeleteProductById)

	return productHandler
}

func (productHandler *ProductHandler) GetProducts(ctx context.Context, queryParam *dto.GetProductsRequestQueryParam) (*dto.SuccessResponse[[]dto.ProductDTO], error) {
	products, err := productHandler.productService.GetProducts(ctx, queryParam)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Message = "Get products failed"
		res.Error_ = "Internal Server Error"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToProductDTOs(products)
	res := &dto.SuccessResponse[[]dto.ProductDTO]{}
	res.Body.Status = http.StatusOK
	res.Body.Message = "Get products successful"
	res.Body.Data = data
	res.Body.Total = len(data)
	return res, nil
}

func (productHandler *ProductHandler) GetProductById(ctx context.Context, reqDTO *dto.GetProductByIdRequest) (*dto.SuccessResponse[*dto.ProductDTO], error) {
	id := reqDTO.Id

	foundProduct, err := productHandler.productService.GetProductById(ctx, id)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Message = "Get product by id failed"
		res.Error_ = "Bad Request"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToProductDTO(foundProduct)
	res := &dto.SuccessResponse[*dto.ProductDTO]{}
	res.Body.Status = http.StatusOK
	res.Body.Message = "Get product by id successful"
	res.Body.Data = data
	return res, nil
}

func (productHandler *ProductHandler) CreateProduct(ctx context.Context, reqDTO *dto.CreateProductRequest) (*dto.SuccessResponse[any], error) {
	if err := productHandler.productService.CreateProduct(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Message = "Create product failed"
		res.Error_ = "Bad Request"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse[any]{}
	res.Body.Status = http.StatusOK
	res.Body.Message = "Create product successful"
	return res, nil
}

func (productHandler *ProductHandler) UpdateProductById(ctx context.Context, reqDTO *dto.UpdateProductRequest) (*dto.SuccessResponse[any], error) {
	id := reqDTO.Id

	if err := productHandler.productService.UpdateProductById(ctx, id, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Message = "Update product failed"
		res.Error_ = "Bad Request"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse[any]{}
	res.Body.Status = http.StatusOK
	res.Body.Message = "Update product successful"
	return res, nil
}

func (productHandler *ProductHandler) DeleteProductById(ctx context.Context, reqDTO *dto.DeleteProductRequest) (*dto.SuccessResponse[any], error) {
	id := reqDTO.Id

	if err := productHandler.productService.DeleteProductById(ctx, id); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Message = "Delete product failed"
		res.Error_ = "Bad Request"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse[any]{}
	res.Body.Status = http.StatusOK
	res.Body.Message = "Delete product successful"
	return res, nil
}
