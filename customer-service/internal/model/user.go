package model

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	Id             int64     `bun:"id,pk,autoincrement"`
	FullName       string    `bun:"full_name,notnull"`
	Email          string    `bun:"email,notnull"`
	Username       string    `bun:"username,notnull"`
	HashedPassword string    `bun:"hashed_password,notnull"`
	Address        string    `bun:"address,notnull"`
	RoleName       string    `bun:"role_name,notnull"`
	CreatedAt      time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt      time.Time `bun:"updated_at,notnull,default:current_timestamp"`
}
