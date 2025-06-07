package apimodels

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID        int64     `bun:"id,pk,autoincrement" json:"id"`
	Email     string    `bun:"email,notnull,unique" json:"email"`
	Name      string    `bun:"name,notnull" json:"name"`
	Age       *int      `bun:"age" json:"age,omitempty"`
	Gender    *string   `bun:"gender" json:"gender,omitempty"`
	Password  string    `bun:"password,notnull" json:"password"`
	CreatedAt time.Time `bun:"created_at,notnull,default:now()" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:now()" json:"updated_at"`
}
