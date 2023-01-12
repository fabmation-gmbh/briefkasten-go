package models

import (
	"time"

	"github.com/uptrace/bun"
	"gopkg.in/guregu/null.v4"
)

type Bookmark struct {
	bun.BaseModel `bun:"bookmark"`

	ID          string      `bun:"id" json:"id"`
	CreatedAt   time.Time   `bun:"created_at" json:"created_at"`
	UserID      string      `bun:"user_id" json:"user_id"`
	CategoryID  string      `bun:"category_id" json:"category_id"`
	URL         string      `bun:"url" json:"url"`
	Image       null.String `bun:"image" json:"image"`
	Description null.String `bun:"description" json:"description"`
}
