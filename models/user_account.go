package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"gopkg.in/guregu/null.v4"
)

type UserAccount struct {
	bun.BaseModel `bun:"user_account"`

	ID            uuid.UUID   `bun:"id" json:"id"`
	CreatedAt     time.Time   `bun:"created_at" json:"created_at"`
	Name          string      `bun:"name" json:"name"`
	Email         string      `bun:"email" json:"email"`
	EmailVerified bool        `bun:"email_verified" json:"email_verified"`
	Image         null.String `bun:"image" json:"image"`
}

// Create inserts the object into the table.
func (u *UserAccount) Create(ctx context.Context) error {
	_, err := db.NewInsert().
		Model(u).
		Returning("*").
		Exec(ctx)

	return errors.Wrap(err, "unable to insert user into DB")
}

// GetOrCreateUser returns the user object from the database.
// The user will be searched by its username.
// If the user does not exist, it will be created.
func GetOrCreateUser(ctx context.Context, u UserAccount) (UserAccount, error) {
	var ret UserAccount

	err := db.NewSelect().
		Model(&ret).
		Where("email = ?", u.Email).
		Limit(1).
		Scan(ctx)
	if err != nil && IsNoRows(err) {
		return u, u.Create(ctx)
	}

	return ret, errors.Wrap(err, "unable to find user in DB")
}
