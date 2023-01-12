package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type TagOnBookmark struct {
	bun.BaseModel `bun:"tag_on_bookmark"`

	BookmarkID uuid.UUID `bun:"bookmark_id" json:"bookmark_id"`
	TagID      uuid.UUID `bun:"tag_id" json:"tag_id"`
}
