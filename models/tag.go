package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
	"gopkg.in/guregu/null.v4"
)

type Tag struct {
	bun.BaseModel `bun:"tag"`

	ID          uuid.UUID   `bun:"id" json:"id"`
	CreatedAt   time.Time   `bun:"created_at" json:"created_at"`
	UserID      uuid.UUID   `bun:"user_id" json:"user_id"`
	Name        string      `bun:"name" json:"name"`
	Description null.String `bun:"description" json:"description"`
	Emoji       null.String `bun:"emoji" json:"emoji"`
}

func (t *Tag) Delete(ctx context.Context) error {
	_, err := db.NewDelete().
		Model(t).
		Where("user_id = ?", t.UserID).
		Where("id = ?", t.ID).
		Exec(ctx)

	return errors.Wrap(err, "unable to delete entry")
}

func (t *Tag) Create(ctx context.Context) error {
	_, err := db.NewInsert().
		Model(t).
		Returning("*").
		Exec(ctx)

	return errors.Wrap(err, "unable to insert tag")
}

func (t *Tag) Update(ctx context.Context) error {
	q := db.NewUpdate().
		Model(t).
		// NOTE: Some fields are immutable.
		Where("user_id = ?", t.UserID).
		Where("id = ?", t.ID).
		Returning("*")

	if t.Name != "" {
		q.Set("name = ?", t.Name)
	}
	if !t.Description.IsZero() {
		q.Set("description = ?", t.Description)
	}
	if !t.Emoji.IsZero() {
		q.Set("emoji = ?", t.Emoji)
	}

	_, err := q.Exec(ctx)

	return errors.Wrap(err, "unable to delete entry")
}

// GetTagsByUserID returns a list of tags from the user.
//
// This method is paginated, meaning that last ID and CreatedAt must be provided.
// For the first request the NilUUID and the zero value of time.Time can be used.
func GetTagsByUserID(ctx context.Context, userID uuid.UUID, startID uuid.UUID, startTime time.Time, limit uint) ([]Tag, error) {
	ret := make([]Tag, 0, limit)

	q := db.NewSelect().
		Model((*Tag)(nil)).
		Where("user_id = ?", userID).
		Order("created_at DESC")

	if startID != uuid.Nil && !startTime.IsZero() {
		q.Where("(id, created_at) =< (?, ?)", startID, startTime)
	}

	err := q.Limit(int(limit)).
		Scan(ctx, &ret)

	return ret, errors.Wrap(err, "unable to retrieve tags")
}

// DeleteTag deletes the tag entry.
func DeleteTag(ctx context.Context, userID, tagID uuid.UUID) error {
	t := Tag{
		ID:     tagID,
		UserID: userID,
	}

	return t.Delete(ctx)
}
