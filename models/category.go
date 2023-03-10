package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"gopkg.in/guregu/null.v4"
)

type Category struct {
	bun.BaseModel `bun:"category"`

	ID          uuid.UUID   `bun:"id" json:"id"`
	CreatedAt   time.Time   `bun:"created_at" json:"created_at"`
	UserID      uuid.UUID   `bun:"user_id" json:"user_id"`
	Name        string      `bun:"name" json:"name"`
	Description null.String `bun:"description" json:"description"`
}
