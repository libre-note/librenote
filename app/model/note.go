package model

import (
	"database/sql"
	"time"
)

//var Colors = map[string]string{
//	"red":       "red",
//	"orange":    "orange",
//	"yellow":    "yellow",
//	"green":     "green",
//	"teal":      "teal",
//	"blue":      "blue",
//	"dark blue": "dark blue",
//	"purple":    "purple",
//	"pink":      "pink",
//	"brown":     "brown",
//	"gray":      "gray",
//}
//
//var NoteTypes = map[string]string{
//	"note": "note",
//	"list": "list",
//}

type Note struct {
	ID         int32          `json:"id"`
	UserID     int32          `json:"user_id"`
	Title      sql.NullString `json:"title"`
	Color      string         `json:"color"`
	Type       string         `json:"type"`
	IsPinned   bool           `json:"is_pinned"`
	IsArchived bool           `json:"is_archived"`
	IsTrashed  bool           `json:"is_trashed"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

type NotesItem struct {
	ID        int32     `json:"id"`
	NoteID    int32     `json:"note_id"`
	Text      *string   `json:"text"`
	IsChecked bool      `json:"is_checked"`
	CreatedAt time.Time `json:"created_at"`
}

type NotesLabel struct {
	NoteID  int32 `json:"note_id"`
	LabelID int32 `json:"label_id"`
}
