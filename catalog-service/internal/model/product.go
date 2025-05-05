package model

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:products"`

	Id                 int64           `bun:"id,pk,autoincrement" json:"id"`
	Name               string          `bun:"name,notnull" json:"name"`
	Description        string          `bun:"description,notnull" json:"description"`
	Price              decimal.Decimal `bun:"price,notnull" json:"price"`
	DiscountPercentage int32           `bun:"discount_percentage,notnull" json:"discount_percentage"`
	Stock              int32           `bun:"stock,notnull" json:"stock"`
	ImageURL           string          `bun:"image_url,notnull" json:"image_url"`
	CategoryId         int64           `bun:"category_id,notnull" json:"category_id"`
	CreatedAt          time.Time       `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt          time.Time       `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}

var ProductMappingIndexForElasticsearch = `
{
  "settings": {
    "analysis": {
      "analyzer": {
        "ngram_analyzer": {
          "type": "custom",
          "tokenizer": "ngram_tokenizer",
          "filter": ["lowercase"]
        }
      },
      "tokenizer": {
        "ngram_tokenizer": {
          "type": "edge_ngram",
          "min_gram": 1,
          "max_gram": 10,
          "token_chars": ["digit"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": { "type": "long" },
      "name": { "type": "text", "analyzer": "standard" },
      "description": { "type": "text", "analyzer": "standard" },
      "price": {
        "type": "double",
        "fields": {
          "as_text": {
            "type": "text",
            "analyzer": "ngram_analyzer",
            "search_analyzer": "standard"
          }
        }
      },
      "discount_percentage": { "type": "integer" },
      "stock": { "type": "integer" },
      "image_url": { "type": "keyword" },
      "category_id": { "type": "long" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}`

var ProductValidSortField = map[string]string{
	"id":                  "id",
	"name":                "name.keyword",
	"description":         "description.keyword",
	"price":               "price.keyword",
	"discount_percentage": "discount_percentage",
	"stock":               "stock",
	"image_url":           "image_url.keyword",
	"category_id":         "category_id",
	"created_at":          "created_at",
	"updated_at":          "updated_at",
}
