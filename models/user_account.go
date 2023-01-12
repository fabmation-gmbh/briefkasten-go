package models

import (
	"time"

	"github.com/uptrace/bun"
	"gopkg.in/guregu/null.v4"
)

type UserAccount struct {
	bun.BaseModel `bun:"user_account"`

	ID            string      `bun:"id" json:"id"`
	CreatedAt     time.Time   `bun:"created_at" json:"created_at"`
	Name          string      `bun:"name" json:"name"`
	Email         string      `bun:"email" json:"email"`
	EmailVerified bool        `bun:"email_verified" json:"email_verified"`
	Image         null.String `bun:"image" json:"image"`
}
