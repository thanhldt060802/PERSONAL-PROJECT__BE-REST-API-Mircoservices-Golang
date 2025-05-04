package model

import (
	"github.com/uptrace/bun"
)

type CartItem struct {
	bun.BaseModel `bun:"table:cart_items"`

	Id        int64 `bun:"id,pk,autoincrement"`
	CartId    int64 `bun:"cart_id,notnull"`
	ProductId int64 `bun:"product_id,notnull"`
	Quantity  int32 `bun:"quantity,notnull"`
}
