package model

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:products"`

	Id                 int64           `bun:"id,pk,autoincrement"`
	Name               string          `bun:"name,notnull"`
	Description        string          `bun:"description,notnull"`
	Price              decimal.Decimal `bun:"price,notnull"`
	DiscountPercentage int32           `bun:"discount_percentage,notnull"`
	Stock              int32           `bun:"stock,notnull"`
	ImageURL           string          `bun:"image_url,notnull"`
	CategoryId         int64           `bun:"category_id,notnull"`
	CreatedAt          time.Time       `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt          time.Time       `bun:"updated_at,notnull,default:current_timestamp"`
}
