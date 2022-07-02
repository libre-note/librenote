package model

import (
	"database/sql"
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
	IsPinned   int8           `json:"is_pinned"`
	IsArchived int8           `json:"is_archived"`
	IsTrashed  int8           `json:"is_trashed"`
	CreatedAt  string         `json:"created_at"`
	UpdatedAt  string         `json:"updated_at"`
}

type NotesItem struct {
	ID        int32   `json:"id"`
	NoteID    int32   `json:"note_id"`
	Text      *string `json:"text"`
	IsChecked int8    `json:"is_checked"`
	CreatedAt string  `json:"created_at"`
}

type NotesLabel struct {
	NoteID  int32 `json:"note_id"`
	LabelID int32 `json:"label_id"`
}
