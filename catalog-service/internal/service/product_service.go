package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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

	SyncAllProductsToElasticsearch(ctx context.Context) error
	SyncAddProductToElasticsearch(ctx context.Context, newProduct *model.Product) error
	SyncUpdateProductOnElasticsearch(ctx context.Context, updatedProduct *model.Product) error
	SyncDeleteProductFromElasticsearch(ctx context.Context, id int64) error

	SearchProducts(ctx context.Context, reqDTO *dto.SearchProductsWithQueryParamRequest) ([]model.Product, error)
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

	if err := productService.SyncAddProductToElasticsearch(ctx, &newProduct); err != nil {
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

	if err := productService.SyncUpdateProductOnElasticsearch(ctx, foundProduct); err != nil {
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

	if err := productService.SyncDeleteProductFromElasticsearch(ctx, reqDTO.Id); err != nil {
		return err
	}

	return nil
}

// func (productService *productService) SyncAllProductsToElasticsearch(ctx context.Context) error {
// 	products, err := productService.productRepository.GetAll(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	for _, product := range products {
// 		// doc := map[string]interface{}{
// 		// 	"id":                  product.Id,
// 		// 	"name":                product.Name,
// 		// 	"description":         product.Description,
// 		// 	"price":               product.Price,
// 		// 	"discount_percentage": product.DiscountPercentage,
// 		// 	"stock":               product.Stock,
// 		// 	"image_url":           product.ImageURL,
// 		// 	"category_id":         product.CategoryId,
// 		// 	"created_at":          product.CreatedAt,
// 		// 	"updated_at":          product.UpdatedAt,
// 		// }

// 		_, err := infrastructure.ElasticsearchClient.Index(
// 			"products", // Tên index
// 			esutil.NewJSONReader(product),
// 			infrastructure.ElasticsearchClient.Index.WithDocumentID(strconv.FormatInt(product.Id, 10)),
// 		)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (productService *productService) SyncAllProductsToElasticsearch(ctx context.Context) error {
	// Tạo index với product mapping
	res, err := infrastructure.ElasticsearchClient.Indices.Create("products",
		infrastructure.ElasticsearchClient.Indices.Create.WithBody(bytes.NewReader([]byte(model.ProductMappingIndexForElasticsearch))))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("create index on elasticsearch faield: %s", res.String())
	}

	// Sync dữ liệu từ database theo model.Product
	products, err := productService.productRepository.GetAll(ctx)
	if err != nil {
		return err
	}

	indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client: infrastructure.ElasticsearchClient,
		Index:  "products",
	})
	if err != nil {
		return err
	}
	defer func() {
		if cerr := indexer.Close(ctx); cerr != nil {
			log.Printf("error closing bulk indexer: %v", cerr)
		}
	}()

	// Lặp qua tất cả sản phẩm và chuẩn bị dữ liệu để index vào Elasticsearch
	for _, product := range products {
		// Chuyển đổi product thành JSON
		data, err := json.Marshal(product)
		if err != nil {
			log.Printf("failed to marshal product ID %d: %v", product.Id, err)
			continue
		}

		// Tạo một BulkIndexerItem để thêm vào Elasticsearch
		err = indexer.Add(ctx, esutil.BulkIndexerItem{
			Action:     "index",                           // hành động là "index" (tạo hoặc cập nhật tài liệu)
			DocumentID: strconv.FormatInt(product.Id, 10), // ID của tài liệu Elasticsearch
			Body:       bytes.NewReader(data),             // Chuyển đổi dữ liệu JSON thành io.Reader
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, resp esutil.BulkIndexerResponseItem, err error) {
				// Nếu có lỗi khi thêm dữ liệu vào Elasticsearch
				if err != nil {
					log.Printf("bulk index error: %v", err)
				} else {
					log.Printf("failed to index product ID %s: %s", item.DocumentID, resp.Error.Reason)
				}
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (productService *productService) SyncAddProductToElasticsearch(ctx context.Context, newProduct *model.Product) error {
	// Index the product in Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Index(
		"products", // Index name
		esutil.NewJSONReader(newProduct),
		infrastructure.ElasticsearchClient.Index.WithDocumentID(strconv.FormatInt(newProduct.Id, 10)), // Use product ID as the document ID
		infrastructure.ElasticsearchClient.Index.WithRefresh("true"),                                  // Refresh index immediately
	)
	if err != nil {
		return fmt.Errorf("failed to add product to elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %s", res.String())
	}

	return nil
}

func (productService *productService) SyncUpdateProductOnElasticsearch(ctx context.Context, updatedProduct *model.Product) error {
	// Index the product in Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Index(
		"products", // Index name
		esutil.NewJSONReader(updatedProduct),
		infrastructure.ElasticsearchClient.Index.WithDocumentID(strconv.FormatInt(updatedProduct.Id, 10)), // Use product ID as the document ID
		infrastructure.ElasticsearchClient.Index.WithRefresh("true"),                                      // Refresh index immediately
	)
	if err != nil {
		return fmt.Errorf("failed to update product on elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %s", res.String())
	}

	return nil
}

func (productService *productService) SyncDeleteProductFromElasticsearch(ctx context.Context, id int64) error {
	// Index the product in Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Delete(
		"products",                // Index name
		strconv.FormatInt(id, 10), // Document ID
		infrastructure.ElasticsearchClient.Delete.WithRefresh("true"), // Refresh index immediately
	)
	if err != nil {
		return fmt.Errorf("failed to delete product from elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error: %s", res.String())
	}

	return nil
}

func (productService *productService) SearchProducts(ctx context.Context, reqDTO *dto.SearchProductsWithQueryParamRequest) ([]model.Product, error) {
	// Build Elasticsearch query
	// esQuery := map[string]interface{}{
	// 	"from": reqDTO.Offset,
	// 	"size": reqDTO.Limit,
	// 	"query": map[string]interface{}{
	// 		"match": map[string]interface{}{
	// 			"name": reqDTO.Query,
	// 		},
	// 	},
	// }

	// Build Elasticsearch query
	esQuery := map[string]interface{}{
		"from": reqDTO.Offset,
		"size": reqDTO.Limit,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  reqDTO.Query,
				"fields": []string{"name", "description", "price.as_text"},
			},
		},
		// "sort": utils.ParseSortByUsingElasticsearchForProduct(reqDTO.SortBy),
	}

	// Convert query to JSON
	esQueryJSON, err := json.Marshal(esQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query")
	}

	// Send request to Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Search(
		infrastructure.ElasticsearchClient.Search.WithContext(ctx),
		infrastructure.ElasticsearchClient.Search.WithIndex("products"),
		infrastructure.ElasticsearchClient.Search.WithBody(bytes.NewReader(esQueryJSON)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Send request to Elasticsearch
	// res, err := infrastructure.ElasticsearchClient.Search(
	// 	infrastructure.ElasticsearchClient.Search.WithContext(ctx),
	// 	infrastructure.ElasticsearchClient.Search.WithIndex("products"),
	// 	infrastructure.ElasticsearchClient.Search.WithQuery(reqDTO.Query),
	// 	infrastructure.ElasticsearchClient.Search.WithPretty(),
	// )
	// if err != nil {
	// 	return nil, err
	// }
	// defer res.Body.Close()

	// Parse response
	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error, %v", res.String())
	}
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

	// Extract products
	products := make([]model.Product, len(esResponse.Hits.Hits))
	for i, hit := range esResponse.Hits.Hits {
		products[i] = hit.Source
	}

	return products, nil
}
