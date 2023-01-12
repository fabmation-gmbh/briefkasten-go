package models

import "github.com/uptrace/bun"

type TagOnBookmark struct {
	bun.BaseModel `bun:"tag_on_bookmark"`

	BookmarkID string `bun:"bookmark_id" json:"bookmark_id"`
	TagID      string `bun:"tag_id" json:"tag_id"`
}
