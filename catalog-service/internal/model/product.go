package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:products"`

	Id                 int64     `bun:"id,pk,autoincrement" json:"id"`
	Name               string    `bun:"name,notnull" json:"name"`
	Description        string    `bun:"description,notnull" json:"description"`
	Sex                string    `bun:"sex,notnull" json:"sex"`
	Price              int64     `bun:"price,notnull" json:"price"`
	DiscountPercentage int32     `bun:"discount_percentage,notnull" json:"discount_percentage"`
	Stock              int32     `bun:"stock,notnull" json:"stock"`
	ImageURL           string    `bun:"image_url,notnull" json:"image_url"`
	CategoryId         int64     `bun:"category_id,notnull" json:"category_id"`
	CreatedAt          time.Time `bun:"created_at,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt          time.Time `bun:"updated_at,notnull,default:current_timestamp" json:"updated_at"`
}

// Integrate with Elasticsearch

var ProductSchemaElasticsearch = `
{
  "mappings": {
    "properties": {
      "id": { "type": "long" },
      "name": {
        "type": "text",
        "analyzer": "standard",
        "fields": {
          "keyword": { "type": "keyword" }
        }
      },
      "description": {
        "type": "text",
        "analyzer": "standard",
        "fields": {
          "keyword": { "type": "keyword" }
        }
      },
      "sex": {
        "type": "text",
        "analyzer": "standard",
        "fields": {
          "keyword": { "type": "keyword" }
        }
      },
      "price": { "type": "long" },
      "discount_percentage": { "type": "integer" },
      "stock": { "type": "integer" },
      "image_url": { "type": "keyword" },
      "category_id": { "type": "long" },
      "created_at": { "type": "date" },
      "updated_at": { "type": "date" }
    }
  }
}`

var MapSortFieldProductSchemaElasticsearch = map[string]string{
	"id":                  "id",
	"name":                "name.keyword",
	"description":         "description.keyword",
	"sex":                 "sex.keyword",
	"price":               "price",
	"discount_percentage": "discount_percentage",
	"stock":               "stock",
	"image_url":           "image_url.keyword",
	"category_id":         "category_id",
	"created_at":          "created_at",
	"updated_at":          "updated_at",
}
