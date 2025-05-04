package service

import (
	"context"
	"encoding/json"
	"fmt"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/dto"
	"thanhldt060802/internal/model"
	"thanhldt060802/internal/repository"
	"thanhldt060802/utils"
	"time"

	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type productService struct {
	productRepository  repository.ProductRepository
	categoryRepository repository.CategoryRepository
}

type ProductService interface {
	GetProducts(ctx context.Context, reqDTO *dto.GetProductsWithQueryParamRequest) ([]model.Product, error)
	GetProductById(ctx context.Context, reqDTO *dto.GetProductByIdRequest) (*model.Product, error)
	GetProductsByCategoryId(ctx context.Context, reqDTO *dto.GetProductsByCategoryIdWithQueryParamRequest) ([]model.Product, error)
	CreateProduct(ctx context.Context, reqDTO *dto.CreateProductRequest) error
	UpdateProductById(ctx context.Context, reqDTO *dto.UpdateProductRequest) error
	DeleteProductById(ctx context.Context, reqDTO *dto.DeleteProductRequest) error

	SyncProductsToElasticsearch(ctx context.Context) error

	SearchProducts(ctx context.Context, reqDTO *dto.SearchProductsRequest) ([]model.Product, error)
}

func NewProductService(productRepository repository.ProductRepository, categoryRepository repository.CategoryRepository) ProductService {
	return &productService{
		productRepository:  productRepository,
		categoryRepository: categoryRepository,
	}
}

func (productService *productService) GetProducts(ctx context.Context, reqDTO *dto.GetProductsWithQueryParamRequest) ([]model.Product, error) {
	sortFields := utils.ParseSortBy(reqDTO.SortBy)

	products, err := productService.productRepository.Get(ctx, reqDTO.Offset, reqDTO.Limit, sortFields)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (productService *productService) GetProductById(ctx context.Context, reqDTO *dto.GetProductByIdRequest) (*model.Product, error) {
	foundProduct, err := productService.productRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
		return nil, err
	}

	return foundProduct, nil
}

func (productService *productService) GetProductsByCategoryId(ctx context.Context, reqDTO *dto.GetProductsByCategoryIdWithQueryParamRequest) ([]model.Product, error) {
	sortFields := utils.ParseSortBy(reqDTO.SortBy)

	products, err := productService.productRepository.GetByCategoryId(ctx, reqDTO.CategoryId, reqDTO.Offset, reqDTO.Limit, sortFields)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (productService *productService) CreateProduct(ctx context.Context, reqDTO *dto.CreateProductRequest) error {
	if _, err := productService.categoryRepository.GetById(ctx, reqDTO.Body.CategoryId); err != nil {
		return fmt.Errorf("id of category is not valid")
	}

	newProduct := model.Product{
		Name:               reqDTO.Body.Name,
		Description:        reqDTO.Body.Description,
		Price:              reqDTO.Body.Price,
		DiscountPercentage: reqDTO.Body.DiscountPercentage,
		Stock:              reqDTO.Body.Stock,
		ImageURL:           reqDTO.Body.ImageURL,
		CategoryId:         reqDTO.Body.CategoryId,
	}
	if err := productService.productRepository.Create(ctx, &newProduct); err != nil {
		return err
	}

	return nil
}

func (productService *productService) UpdateProductById(ctx context.Context, reqDTO *dto.UpdateProductRequest) error {
	foundProduct, err := productService.productRepository.GetById(ctx, reqDTO.Id)
	if err != nil {
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
		if _, err := productService.categoryRepository.GetById(ctx, *reqDTO.Body.CategoryId); err != nil {
			return fmt.Errorf("id of category is not valid")
		}
		foundProduct.CategoryId = *reqDTO.Body.CategoryId
	}
	foundProduct.UpdatedAt = time.Now().UTC()

	if err := productService.productRepository.UpdateById(ctx, reqDTO.Id, foundProduct); err != nil {
		return err
	}

	return nil
}

func (productService *productService) DeleteProductById(ctx context.Context, reqDTO *dto.DeleteProductRequest) error {
	if _, err := productService.productRepository.GetById(ctx, reqDTO.Id); err != nil {
		return fmt.Errorf("id of product is not valid")
	}

	if err := productService.productRepository.DeleteById(ctx, reqDTO.Id); err != nil {
		return err
	}

	return nil
}

func (productService *productService) SyncProductsToElasticsearch(ctx context.Context) error {
	products, err := productService.productRepository.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, product := range products {
		doc := map[string]interface{}{
			"id":                  product.Id,
			"name":                product.Name,
			"description":         product.Description,
			"price":               product.Price,
			"discount_percentage": product.DiscountPercentage,
			"stock":               product.Stock,
			"image_url":           product.ImageURL,
			"category_id":         product.CategoryId,
			"created_at":          product.CreatedAt,
			"updated_at":          product.UpdatedAt,
		}

		_, err := infrastructure.ElasticsearchClient.Index(
			"products", // Tên index
			esutil.NewJSONReader(doc),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (productService *productService) SearchProducts(ctx context.Context, reqDTO *dto.SearchProductsRequest) ([]model.Product, error) {
	res, err := infrastructure.ElasticsearchClient.Search(
		infrastructure.ElasticsearchClient.Search.WithContext(ctx),
		infrastructure.ElasticsearchClient.Search.WithIndex("products"),
		infrastructure.ElasticsearchClient.Search.WithQuery(reqDTO.Query),
		infrastructure.ElasticsearchClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var esResponse struct {
		Hits struct {
			Hits []struct {
				Source model.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&esResponse); err != nil {
		return nil, err
	}

	products := make([]model.Product, len(esResponse.Hits.Hits))
	for i, hit := range esResponse.Hits.Hits {
		products[i] = hit.Source
	}

	return products, nil
}
