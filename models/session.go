package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Session struct {
	bun.BaseModel `bun:"session"`

	Token   string    `bun:"token" json:"token"`
	UserID  string    `bun:"user_id" json:"user_id"`
	Expires time.Time `bun:"expires" json:"expires"`
}
