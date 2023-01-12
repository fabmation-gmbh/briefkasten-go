package models

import (
	"time"

	"github.com/uptrace/bun"
	"gopkg.in/guregu/null.v4"
)

type Category struct {
	bun.BaseModel `bun:"category"`

	ID          string      `bun:"id" json:"id"`
	CreatedAt   time.Time   `bun:"created_at" json:"created_at"`
	UserID      string      `bun:"user_id" json:"user_id"`
	Name        string      `bun:"name" json:"name"`
	Description null.String `bun:"description" json:"description"`
}
