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

	// Get products by category id
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/products/category-id/{category_id}",
		Summary:     "/products/category-id/{category_id}",
		Description: "Get products by category id.",
		Tags:        []string{"Product"},
	}, productHandler.GetProductsByCategoryId)

	// Create product
	huma.Register(api, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/products",
		Summary:     "/products",
		Description: "Create product.",
		Tags:        []string{"Product"},
		// Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, productHandler.CreateProduct)

	// Update product by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodPut,
		Path:        "/products/id/{id}",
		Summary:     "/products/id/{id}",
		Description: "Update product by id.",
		Tags:        []string{"Product"},
		// Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, productHandler.UpdateProductById)

	// Delete product by id
	huma.Register(api, huma.Operation{
		Method:      http.MethodDelete,
		Path:        "/products/id/{id}",
		Summary:     "/products/id/{id}",
		Description: "Delete product by id.",
		Tags:        []string{"Product"},
		// Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, productHandler.DeleteProductById)

	// Sync all products to Elasticsearch
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/products/sync-to-elasticsearch",
		Summary:     "/products/sync-to-elasticsearch",
		Description: "Sync all products to Elasticsearch.",
		Tags:        []string{"Product"},
		// Middlewares: huma.Middlewares{authMiddleware.Authentication, authMiddleware.RequireAdmin},
	}, productHandler.SyncAllProductsToElasticsearch)

	// Get products with Elasticsearch
	huma.Register(api, huma.Operation{
		Method:      http.MethodGet,
		Path:        "/products/elasticsearch",
		Summary:     "/products/elasticsearch",
		Description: "Get products with Elasticsearch.",
		Tags:        []string{"Product"},
	}, productHandler.GetProducsWithElasticsearch)

	return productHandler
}

func (productHandler *ProductHandler) GetProducts(ctx context.Context, reqDTO *dto.GetProductsRequest) (*dto.PaginationBodyResponseList[dto.ProductView], error) {
	products, err := productHandler.productService.GetProducts(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get products failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToListProductView(products)
	res := &dto.PaginationBodyResponseList[dto.ProductView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get products successful"
	res.Body.Data = data
	res.Body.Total = len(data)
	return res, nil
}

func (productHandler *ProductHandler) GetProductById(ctx context.Context, reqDTO *dto.GetProductByIdRequest) (*dto.BodyResponse[dto.ProductView], error) {
	foundProduct, err := productHandler.productService.GetProductById(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Get product by id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToProductView(foundProduct)
	res := &dto.BodyResponse[dto.ProductView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get product by id successful"
	res.Body.Data = *data
	return res, nil
}

func (productHandler *ProductHandler) GetProductsByCategoryId(ctx context.Context, reqDTO *dto.GetProductsByCategoryIdRequest) (*dto.PaginationBodyResponseList[dto.ProductView], error) {
	products, err := productHandler.productService.GetProductsByCategoryId(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get products by category id failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToListProductView(products)
	res := &dto.PaginationBodyResponseList[dto.ProductView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get products category id successful"
	res.Body.Data = data
	res.Body.Total = len(data)
	return res, nil
}

func (productHandler *ProductHandler) CreateProduct(ctx context.Context, reqDTO *dto.CreateProductRequest) (*dto.SuccessResponse, error) {
	if err := productHandler.productService.CreateProduct(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Create product failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Create product successful"
	return res, nil
}

func (productHandler *ProductHandler) UpdateProductById(ctx context.Context, reqDTO *dto.UpdateProductByIdRequest) (*dto.SuccessResponse, error) {
	if err := productHandler.productService.UpdateProductById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Update product failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Update product successful"
	return res, nil
}

func (productHandler *ProductHandler) DeleteProductById(ctx context.Context, reqDTO *dto.DeleteProductByIdRequest) (*dto.SuccessResponse, error) {
	if err := productHandler.productService.DeleteProductById(ctx, reqDTO); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusBadRequest
		res.Code = "ERR_BAD_REQUEST"
		res.Message = "Delete product failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Delete product successful"
	return res, nil
}

func (productHandler *ProductHandler) SyncAllProductsToElasticsearch(ctx context.Context, reqDTO *struct{}) (*dto.SuccessResponse, error) {
	if err := productHandler.productService.SyncAllProductsToElasticsearch(ctx); err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Sync all products to Elasticsearch failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	res := &dto.SuccessResponse{}
	res.Body.Code = "OK"
	res.Body.Message = "Sync all products to Elasticsearch successful"
	return res, nil
}

func (productHandler *ProductHandler) GetProducsWithElasticsearch(ctx context.Context, reqDTO *dto.GetProductsWithElasticsearchRequest) (*dto.PaginationBodyResponseList[dto.ProductView], error) {
	products, err := productHandler.productService.GetProductsWithElasticsearch(ctx, reqDTO)
	if err != nil {
		res := &dto.ErrorResponse{}
		res.Status = http.StatusInternalServerError
		res.Code = "ERR_INTERNAL_SERVER"
		res.Message = "Get products with Elasticsearch failed"
		res.Details = []string{err.Error()}
		return nil, res
	}

	data := dto.ToListProductView(products)
	res := &dto.PaginationBodyResponseList[dto.ProductView]{}
	res.Body.Code = "OK"
	res.Body.Message = "Get products with Elasticsearch successful"
	res.Body.Data = data
	res.Body.Total = len(data)
	return res, nil
}
