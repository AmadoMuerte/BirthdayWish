package models

import (
	"time"

	"github.com/uptrace/bun"
)

type ShareLink struct {
	bun.BaseModel `bun:"table:share_links,alias:share_links"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	Token     string    `bun:"token" json:"token"`
	UserID    int64     `bun:"user_id" json:"user_id"`
	ExpiresAt time.Time `bun:"expires_at" json:"expires_at"`
	CreatedAt time.Time `bun:"created_at" json:"created_at"`
}
