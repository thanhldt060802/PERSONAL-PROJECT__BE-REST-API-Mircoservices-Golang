package service

import (
	"context"
	"fmt"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/repository"
	"time"
)

type productService struct {
	productRepository repository.ProductRepository
}

type ProductService interface {
	IsProductExistedById(ctx context.Context, id int64) (bool, error)
	GetProducts(ctx context.Context, queryParam *dto.GetProductsRequestQueryParam) ([]model.Product, error)
	GetProductById(ctx context.Context, id int64) (*model.Product, error)
	CreateProduct(ctx context.Context, reqDTO *dto.CreateProductRequest) error
	UpdateProductById(ctx context.Context, id int64, reqDTO *dto.UpdateProductRequest) error
	DeleteProductById(ctx context.Context, id int64) error
}

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{productRepository: productRepository}
}

func (productService *productService) IsProductExistedById(ctx context.Context, id int64) (bool, error) {
	existed, err := productService.productRepository.ExistsById(ctx, id)
	if err != nil {
		return false, fmt.Errorf("get product by id failed: %w", err)
	}
	return existed, nil
}

func (productService *productService) GetProducts(ctx context.Context, queryParam *dto.GetProductsRequestQueryParam) ([]model.Product, error) {
	sortFields := dto.ParseSortBy(queryParam.SortBy)

	products, err := productService.productRepository.Get(ctx, queryParam.Offset, queryParam.Limit, sortFields)
	if err != nil {
		return nil, fmt.Errorf("get products failed: %w", err)
	}

	return products, nil
}

func (productService *productService) GetProductById(ctx context.Context, id int64) (*model.Product, error) {
	foundProduct, err := productService.productRepository.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get product by id failed: %w", err)
	}

	return foundProduct, nil
}

func (productService *productService) CreateProduct(ctx context.Context, reqDTO *dto.CreateProductRequest) error {
	// HTTP check valid for category id
	// ...
	newProduct := model.Product{
		Name:               reqDTO.Body.Name,
		Description:        reqDTO.Body.Description,
		Price:              reqDTO.Body.Price,
		DiscountPercentage: reqDTO.Body.DiscountPercentage,
		Stock:              reqDTO.Body.Stock,
		ImageURL:           reqDTO.Body.ImageURL,
		CategoryId:         reqDTO.Body.CategoryId,
	}

	return productService.productRepository.Create(ctx, &newProduct)
}

func (productService *productService) UpdateProductById(ctx context.Context, id int64, reqDTO *dto.UpdateProductRequest) error {
	foundProduct, err := productService.GetProductById(ctx, id)
	if err != nil {
		return err
	}
	if foundProduct == nil {
		return fmt.Errorf("id of product is not valid")
	}

	if reqDTO.Body.Name != nil {
		foundProduct.Name = *reqDTO.Body.Name
	}
	if reqDTO.Body.Description != nil {
		foundProduct.Description = *reqDTO.Body.Description
	}
	if reqDTO.Body.Price != nil {
		foundProduct.Price = *reqDTO.Body.Price
	}
	if reqDTO.Body.DiscountPercentage != nil {
		foundProduct.DiscountPercentage = *reqDTO.Body.DiscountPercentage
	}
	if reqDTO.Body.Stock != nil {
		foundProduct.Stock = *reqDTO.Body.Stock
	}
	if reqDTO.Body.ImageURL != nil {
		foundProduct.ImageURL = *reqDTO.Body.ImageURL
	}
	if reqDTO.Body.CategoryId != nil {
		// HTTP check valid for category id
		// ...
		foundProduct.CategoryId = *reqDTO.Body.CategoryId
	}
	foundProduct.UpdatedAt = time.Now().UTC()

	return productService.productRepository.UpdateById(ctx, id, foundProduct)
}

func (productService *productService) DeleteProductById(ctx context.Context, id int64) error {
	existed, err := productService.IsProductExistedById(ctx, id)
	if err != nil {
		return err
	}
	if !existed {
		return fmt.Errorf("id of product is not valid")
	}

	return productService.productRepository.DeleteById(ctx, id)
}
