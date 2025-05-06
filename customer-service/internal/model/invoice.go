package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Invoice struct {
	bun.BaseModel `bun:"table:invoices"`

	Id          int64     `bun:"id,pk,autoincrement" json:"id"`
	UserId      int64     `bun:"user_id,notnull" json:"user_id"`
	TotalAmount int64     `bun:"total_amount,notnull" json:"total_amount"`
	Status      string    `bun:"status,notnull" json:"status"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}

// Integrate with Elasticsearch

var InvoiceSchemaElasticsearch = `
{
  "mappings": {
    "properties": {
      "id": { "type": "long" },
	  "user_id": { "type": "long" },
	  "total_amount": { "type": "long" },
      "status": {
        "type": "text",
        "analyzer": "standard",
        "fields": {
          "keyword": { "type": "keyword" }
        }
      },
	  "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}`

var MapSortFieldInvoiceSchemaElasticsearch = map[string]string{
	"id":           "id",
	"user_id":      "user_id",
	"total_amount": "total_amount",
	"status":       "status.keyword",
	"created_at":   "created_at",
	"updated_at":   "updated_at",
}

type InvoiceReport struct {
	Sum       float64    `json:"sum"`
	Avg       float64    `json:"avg"`
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
}
