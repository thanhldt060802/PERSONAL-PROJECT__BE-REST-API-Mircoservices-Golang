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
	"time"

	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type invoiceElasticsearchRepository struct {
}

type InvoiceElasticsearchRepository interface {
	SyncAll(ctx context.Context, invoices []model.Invoice) error

	Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField, createdAtGTE string, createdAtLTE string) ([]model.Invoice, error)
	Sum(ctx context.Context, createdAtGTE string, createdAtLTE string) (*float64, error)
	SumAvg(ctx context.Context, createdAtGTE string, createdAtLTE string) (*model.InvoiceReport, error)
}

func NewInvoiceElasticsearchRepository() InvoiceElasticsearchRepository {
	return &invoiceElasticsearchRepository{}
}

func (invoiceElasticsearchRepository *invoiceElasticsearchRepository) Get(ctx context.Context, offset int, limit int, sortFields []utils.SortField, createdAtGTE string, createdAtLTE string) ([]model.Invoice, error) {
	mustConditions := []map[string]interface{}{}

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
				model.MapSortFieldInvoiceSchemaElasticsearch[sortField.Field]: sortField.Direction,
			})
		}
		query["sort"] = _sortFields
	}

	fmt.Println(query)

	// Convert query to JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("marshal query failed")
	}

	// Send request to Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Search(
		infrastructure.ElasticsearchClient.Search.WithContext(ctx),
		infrastructure.ElasticsearchClient.Search.WithIndex("invoices"),
		infrastructure.ElasticsearchClient.Search.WithBody(bytes.NewReader(queryJSON)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse response
	if res.IsError() {
		return nil, fmt.Errorf("get invoices from elasticsearch failed: %s", res.String())
	}
	var elasticsearchResponse struct {
		Hits struct {
			Hits []struct {
				Source model.Invoice `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&elasticsearchResponse); err != nil {
		return nil, fmt.Errorf("unmarshal elasticsearch response failed: %s", err.Error())
	}

	// Extract invoices
	invoices := make([]model.Invoice, len(elasticsearchResponse.Hits.Hits))
	for i, hit := range elasticsearchResponse.Hits.Hits {
		invoices[i] = hit.Source
	}

	return invoices, nil
}

func (invoiceElasticsearchRepository *invoiceElasticsearchRepository) SyncAll(ctx context.Context, invoices []model.Invoice) error {
	// Check if index already exists
	existsRes, err := infrastructure.ElasticsearchClient.Indices.Exists([]string{"invoices"})
	if err != nil {
		return fmt.Errorf("check index existence failed: %s", err.Error())
	}
	defer existsRes.Body.Close()

	// If index does not exists
	if existsRes.StatusCode == 404 {
		// Create index using custom invoice schema
		createRes, err := infrastructure.ElasticsearchClient.Indices.Create("invoice",
			infrastructure.ElasticsearchClient.Indices.Create.WithBody(bytes.NewReader([]byte(model.InvoiceSchemaElasticsearch))))
		if err != nil {
			return err
		}
		defer createRes.Body.Close()

		if createRes.IsError() {
			return fmt.Errorf("create invoices index on elasticsearch faield: %s", createRes.String())
		}

		// Create BulkIndexer on Elasticsearch
		indexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
			Client: infrastructure.ElasticsearchClient,
			Index:  "invoices",
		})
		if err != nil {
			return err
		}
		defer func() {
			if err := indexer.Close(ctx); err != nil {
				log.Printf("Close bulk indexer failed: %s", err.Error())
			}
		}()

		// Add all invoice to BulkIndexer on Elasticsearch
		for _, invoice := range invoices {
			data, err := json.Marshal(invoice)
			if err != nil {
				log.Printf("Marshal invoice with id = %d: %s", invoice.Id, err.Error())
				continue
			}

			// Add invoice to BulkIndexer
			err = indexer.Add(ctx, esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: strconv.FormatInt(invoice.Id, 10),
				Body:       bytes.NewReader(data),
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, resp esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("Bulk index failed: %s", err.Error())
					} else {
						log.Printf("Index invoice with id = %s failed: %s", item.DocumentID, resp.Error.Reason)
					}
				},
			})
			if err != nil {
				return err
			}
		}

		return nil
	}

	return fmt.Errorf("index invoices already exists after first sync all")
}

func (invoiceElasticsearchRepository *invoiceElasticsearchRepository) Sum(ctx context.Context, createdAtGTE string, createdAtLTE string) (*float64, error) {
	mustConditions := []map[string]interface{}{}

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
		"size": 0,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustConditions,
			},
		},
		"aggs": map[string]interface{}{
			"total_amount_sum": map[string]interface{}{
				"sum": map[string]interface{}{
					"field": "total_amount",
				},
			},
		},
	}

	fmt.Println(query)

	// Convert query to JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("marshal query failed")
	}

	// Send request to Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Search(
		infrastructure.ElasticsearchClient.Search.WithContext(ctx),
		infrastructure.ElasticsearchClient.Search.WithIndex("invoices"),
		infrastructure.ElasticsearchClient.Search.WithBody(bytes.NewReader(queryJSON)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse response
	if res.IsError() {
		return nil, fmt.Errorf("get invoices from elasticsearch failed: %s", res.String())
	}
	var elasticsearchResponse struct {
		Aggregations struct {
			TotalAmountSum struct {
				Value float64 `json:"value"`
			} `json:"total_amount_sum"`
		} `json:"aggregations"`
	}
	if err := json.NewDecoder(res.Body).Decode(&elasticsearchResponse); err != nil {
		return nil, fmt.Errorf("unmarshal elasticsearch response failed: %s", err.Error())
	}

	return &elasticsearchResponse.Aggregations.TotalAmountSum.Value, nil
}

func (invoiceElasticsearchRepository *invoiceElasticsearchRepository) SumAvg(ctx context.Context, createdAtGTE string, createdAtLTE string) (*model.InvoiceReport, error) {
	invoiceReport := &model.InvoiceReport{}

	mustConditions := []map[string]interface{}{}

	// If filtering by created_at in range or partial range
	createdAtRange := map[string]interface{}{}
	if createdAtGTE != "" {
		createdAtRange["gte"] = createdAtGTE
		startTime, err := time.Parse("2006-01-02T15:04:05", createdAtGTE)
		if err != nil {
			return nil, err
		}
		invoiceReport.StartTime = &startTime
	}
	if createdAtLTE != "" {
		createdAtRange["lte"] = createdAtLTE
		endTime, err := time.Parse("2006-01-02T15:04:05", createdAtLTE)
		if err != nil {
			return nil, err
		}
		invoiceReport.EndTime = &endTime
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
		"size": 0,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustConditions,
			},
		},
		"aggs": map[string]interface{}{
			"total_amount_sum": map[string]interface{}{
				"sum": map[string]interface{}{
					"field": "total_amount",
				},
			},
			"total_amount_avg": map[string]interface{}{
				"avg": map[string]interface{}{
					"field": "total_amount",
				},
			},
		},
	}

	fmt.Println(query)

	// Convert query to JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("marshal query failed")
	}

	// Send request to Elasticsearch
	res, err := infrastructure.ElasticsearchClient.Search(
		infrastructure.ElasticsearchClient.Search.WithContext(ctx),
		infrastructure.ElasticsearchClient.Search.WithIndex("invoices"),
		infrastructure.ElasticsearchClient.Search.WithBody(bytes.NewReader(queryJSON)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Parse response
	if res.IsError() {
		return nil, fmt.Errorf("get invoices from elasticsearch failed: %s", res.String())
	}
	var elasticsearchResponse struct {
		Aggregations struct {
			TotalAmountSum struct {
				Value float64 `json:"value"`
			} `json:"total_amount_sum"`
			TotalAmountAvg struct {
				Value float64 `json:"value"`
			} `json:"total_amount_avg"`
		} `json:"aggregations"`
	}
	if err := json.NewDecoder(res.Body).Decode(&elasticsearchResponse); err != nil {
		return nil, fmt.Errorf("unmarshal elasticsearch response failed: %s", err.Error())
	}

	invoiceReport.Sum = elasticsearchResponse.Aggregations.TotalAmountSum.Value
	invoiceReport.Avg = elasticsearchResponse.Aggregations.TotalAmountAvg.Value

	return invoiceReport, nil
}
