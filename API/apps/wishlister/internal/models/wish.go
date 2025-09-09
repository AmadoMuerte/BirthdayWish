package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Wish struct {
	bun.BaseModel `bun:"table:wishes,alias:wishes"`

	ID          int64     `bun:"id,pk,autoincrement" json:"id"`
	UserID      int64     `bun:"user_id,notnull" json:"user_id"`
	Link        string    `bun:"link" json:"link"`
	Price       float64   `bun:"price,type:decimal(10,2)" json:"price"`
	ImageURL    string    `bun:"image_url" json:"image"`
	Title       string    `bun:"title" json:"title"`
	Description string    `bun:"description" json:"description"`
	Priority    int32     `bun:"priority" json:"priority"`
	CreatedAt   time.Time `bun:"created_at,notnull,default:now()" json:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at,notnull,default:now()" json:"updated_at"`
}
