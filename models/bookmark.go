package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"gopkg.in/guregu/null.v4"
)

type Bookmark struct {
	bun.BaseModel `bun:"bookmark"`

	ID          uuid.UUID   `bun:"id" json:"id"`
	CreatedAt   time.Time   `bun:"created_at" json:"created_at"`
	UserID      uuid.UUID   `bun:"user_id" json:"user_id"`
	CategoryID  uuid.UUID   `bun:"category_id" json:"category_id"`
	URL         string      `bun:"url" json:"url"`
	Image       null.String `bun:"image" json:"image"`
	Description null.String `bun:"description" json:"description"`
}
