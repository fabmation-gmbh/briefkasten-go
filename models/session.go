package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type Session struct {
	bun.BaseModel `bun:"session"`

	Token   string    `bun:"token" json:"token"`
	UserID  uuid.UUID `bun:"user_id" json:"user_id"`
	Expires time.Time `bun:"expires" json:"expires"`
	RawData string    `bun:"raw_data" json:"raw_data"`
}

// Create inserts the object into the table.
func (s *Session) Create(ctx context.Context) error {
	_, err := db.NewInsert().
		Model(s).
		Returning("*").
		Exec(ctx)

	return errors.Wrap(err, "unable to insert session into DB")
}
