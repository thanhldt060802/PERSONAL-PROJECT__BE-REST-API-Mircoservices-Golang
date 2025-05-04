package model

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type Invoice struct {
	bun.BaseModel `bun:"table:invoices"`

	Id          int64           `bun:"id,pk,autoincrement"`
	UserId      int64           `bun:"user_id,notnull"`
	TotalAmount decimal.Decimal `bun:"total_amount,notnull"`
	Status      string          `bun:"status,notnull"`
	CreatedAt   time.Time       `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt   time.Time       `bun:"updated_at,notnull,default:current_timestamp"`
}
