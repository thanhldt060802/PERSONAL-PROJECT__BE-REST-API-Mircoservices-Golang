package model

import (
	"github.com/shopspring/decimal"
	"github.com/uptrace/bun"
)

type InvoiceDetail struct {
	bun.BaseModel `bun:"table:invoice_details"`

	Id                 int64           `bun:"id,pk,autoincrement"`
	InvoiceId          int64           `bun:"invoice_id,notnull"`
	ProductId          int64           `bun:"product_id,notnull"`
	Price              decimal.Decimal `bun:"price,notnull"`
	DiscountPercentage int32           `bun:"discount_percentage,notnull"`
	Quantity           int32           `bun:"quantity,notnull"`
	TotalPrice         decimal.Decimal `bun:"total_pice,notnull"`
}
