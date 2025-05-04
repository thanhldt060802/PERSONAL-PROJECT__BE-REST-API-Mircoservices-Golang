package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Cart struct {
	bun.BaseModel `bun:"table:carts"`

	Id        int64     `bun:"id,pk,autoincrement"`
	UserId    int64     `bun:"user_id,notnull"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}
