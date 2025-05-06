package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"thanhldt060802/infrastructure"
	"thanhldt060802/internal/model"
	"thanhldt060802/utils"

	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type productElasticsearchRepository struct {
}

type ProductElasticsearchRepository interface {
	SyncAll(ctx context.Context, products []model.Product) error

	SyncCreating(ctx context.Context, newProduct *model.Product) error
	SyncUpdating(ctx context.Context, updatedProduct *model.Product) error
	SyncDeletingById(ctx context.Context, id int64) error

	Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField, name string, priceGTE string, priceLTE string, createdAtGTE string, createdAtLTE string) ([]model.Product, error)
}

func NewProductElasticsearchRepository() ProductElasticsearchRepository {
	return &productElasticsearchRepository{}
}

func (productElasticsearchRepository *productElasticsearchRepository) SyncAll(ctx context.Context, products []model.Product) error {
	// Check if index already exists
	existsRes, err := infrastructure.ElasticsearchClient.Indices.Exists([]string{"products"})
	if err != nil {
		return fmt.Errorf("check index existence failed: %s", err.Error())
	}
	defer existsRes.Body.Close()

	// If index does not exists
	if existsRes.StatusCode == 404 {
		// Create index using custom product schema
		createRes, err := infrastructure.ElasticsearchClient.Indices.Create("products",
			infrastructure.ElasticsearchClient.Indices.Create.WithBody(bytes.NewReader([]byte(model.ProductSchemaElasticsearch))))
		if err != nil {
			return err
		}
		defer createRes.Body.Close()

		if createRes.IsError() {
			return fmt.Errorf("create products index on elasticsearch faield: %s", createRes.String())
		}

		// Create BulkIndexer on Elasticsearch
		indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
			Client: infrastructure.ElasticsearchClient,
			Index:  "products",
		})
		if err != nil {
			return err
		}
		defer func() {
			if err := indexer.Close(ctx); err != nil {
				log.Printf("Close bulk indexer failed: %s", err.Error())
			}
		}()

		// Add all product to BulkIndexer on Elasticsearch
		for _, product := range products {
			data, err := json.Marshal(product)
			if err != nil {
				log.Printf("Marshal product with id = %d: %s", product.Id, err.Error())
				continue
			}

			// Add product to BulkIndexer
			err = indexer.Add(ctx, esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: strconv.FormatInt(product.Id, 10),
				Body:       bytes.NewReader(data),
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, resp esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("Bulk index failed: %s", err.Error())
					} else {
						log.Printf("Index product with id = %s failed: %s", item.DocumentID, resp.Error.Reason)
					}
				},
			})
			if err != nil {
				return err
			}
		}

		return nil
	}

	return fmt.Errorf("index products already exists after first sync all")
}

func (productElasticsearchRepository *productElasticsearchRepository) SyncCreating(ctx context.Context, newProduct *model.Product) error {
	// Add product to Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Index(
		"products",
		esutil.NewJSONReader(newProduct),
		infrastructure.ElasticsearchClient.Index.WithDocumentID(strconv.FormatInt(newProduct.Id, 10)),
		infrastructure.ElasticsearchClient.Index.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("add product to elasticsearch failed: %s", err.Error())
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("add product to elasticsearch failed: %s", res.String())
	}

	return nil
}

func (productElasticsearchRepository *productElasticsearchRepository) SyncUpdating(ctx context.Context, updatedProduct *model.Product) error {
	// Update product on Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Index(
		"products",
		esutil.NewJSONReader(updatedProduct),
		infrastructure.ElasticsearchClient.Index.WithDocumentID(strconv.FormatInt(updatedProduct.Id, 10)),
		infrastructure.ElasticsearchClient.Index.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("update product on elasticsearch failed: %s", err.Error())
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("update product on elasticsearch failed: %s", res.String())
	}

	return nil
}

func (productElasticsearchRepository *productElasticsearchRepository) SyncDeletingById(ctx context.Context, id int64) error {
	// Delete product from Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Delete(
		"products",
		strconv.FormatInt(id, 10),
		infrastructure.ElasticsearchClient.Delete.WithRefresh("true"),
	)
	if err != nil {
		return fmt.Errorf("delete product from elasticsearch failed: %s", err.Error())
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("delete product from elasticsearch failed: %s", res.String())
	}

	return nil
}

func (productElasticsearchRepository *productElasticsearchRepository) Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField, name string, priceGTE string, priceLTE string, createdAtGTE string, createdAtLTE string) ([]model.Product, error) {
	mustConditions := []map[string]interface{}{}

	// If filtering by name
	if name != "" {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match": map[string]interface{}{
				"name": name,
			},
		})
	}

	// If filtering by price in range or partial range
	priceRange := map[string]interface{}{}
	if priceGTE != "" {
		priceRange["gte"] = priceGTE
	}
	if priceLTE != "" {
		priceRange["lte"] = priceLTE
	}
	if len(priceRange) > 0 {
		mustConditions = append(mustConditions, map[string]interface{}{
			"range": map[string]interface{}{
				"price": priceRange,
			},
		})
	}

	// If filtering by created_at in range or partial range
	createdAtRange := map[string]interface{}{}
	if createdAtGTE != "" {
		createdAtRange["gte"] = createdAtGTE
	}
	if createdAtLTE != "" {
		createdAtRange["lte"] = createdAtLTE
	}
	if len(createdAtRange) > 0 {
		createdAtRange["format"] = "strict_date_optional_time" // For format YYYY-MM-ddTHH:mm:ss
		mustConditions = append(mustConditions, map[string]interface{}{
			"range": map[string]interface{}{
				"created_at": createdAtRange,
			},
		})
	}

	// If not filtering -> get all
	if len(mustConditions) == 0 {
		mustConditions = append(mustConditions, map[string]interface{}{
			"match_all": map[string]interface{}{},
		})
	}

	// Setup query
	query := map[string]interface{}{
		"from": offset,
		"size": limit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustConditions,
			},
		},
	}

	// Apply sorting to query
	if len(sortFields) > 0 {
		_sortFields := []map[string]interface{}{}
		for _, sortField := range sortFields {
			_sortFields = append(_sortFields, map[string]interface{}{
				model.MapSortFieldProductSchemaElasticsearch[sortField.Field]: sortField.Direction,
			})
		}
		query["sort"] = _sortFields
	}

	fmt.Println(query)

	// query := map[string]interface{}{
	// 	"from": 0,
	// 	"size": 20,
	// 	"query": map[string]interface{}{
	// 		"bool": map[string]interface{}{
	// 			"must": []interface{}{
	// 				map[string]interface{}{
	// 					"match": map[string]interface{}{
	// 						"name": "Áo",
	// 					},
	// 				},
	// 				map[string]interface{}{
	// 					"range": map[string]interface{}{
	// 						"price": map[string]interface{}{
	// 							"gte": "100000",
	// 							"lte": "1000000",
	// 						},
	// 					},
	// 				},
	// 				map[string]interface{}{
	// 					"range": map[string]interface{}{
	// 						"created_at": map[string]interface{}{
	// 							"format": "strict_date_optional_time",
	// 							"gte":    "2024-01-15T00:00:00",
	// 							"lte":    "2024-02-05T23:59:59",
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	"sort": []map[string]interface{}{
	// 		{"price": "asc"},
	// 		{"name.keyword": "desc"},
	// 	},
	// }

	// Convert query to JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("marshal query failed")
	}

	// Send request to Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Search(
		infrastructure.ElasticsearchClient.Search.WithContext(ctx),
		infrastructure.ElasticsearchClient.Search.WithIndex("products"),
		infrastructure.ElasticsearchClient.Search.WithBody(bytes.NewReader(queryJSON)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse response
	if res.IsError() {
		return nil, fmt.Errorf("get products from elasticsearch failed: %s", res.String())
	}
	var elasticsearchResponse struct {
		Hits struct {
			Hits []struct {
				Source model.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&elasticsearchResponse); err != nil {
		return nil, fmt.Errorf("unmarshal elasticsearch response failed: %s", err.Error())
	}

	// Extract products
	products := make([]model.Product, len(elasticsearchResponse.Hits.Hits))
	for i, hit := range elasticsearchResponse.Hits.Hits {
		products[i] = hit.Source
	}

	return products, nil
}
