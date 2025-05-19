package models

import (
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/pkg/models"
	"github.com/uptrace/bun"
)

type Wishlist struct {
	bun.BaseModel `bun:"table:wishlist,alias:w"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	UserID    int64     `bun:"user_id,notnull" json:"user_id"`
	Link      string    `bun:"link" json:"link"`
	Price     float64   `bun:"price,type:decimal(10,2)" json:"price"`
	Name      string    `bun:"name" json:"name"`
	CreatedAt time.Time `bun:"created_at,notnull,default:now()" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:now()" json:"updated_at"`

	User *models.User `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
}
