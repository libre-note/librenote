package model

import (
	"database/sql"
	"fmt"
	"time"
)

type Color string

const (
	ColorRed      Color = "red"
	ColorOrange   Color = "orange"
	ColorYellow   Color = "yellow"
	ColorGreen    Color = "green"
	ColorTeal     Color = "teal"
	ColorBlue     Color = "blue"
	ColorDarkblue Color = "dark blue"
	ColorPurple   Color = "purple"
	ColorPink     Color = "pink"
	ColorBrown    Color = "brown"
	ColorGray     Color = "gray"
)

func (e *Color) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Color(s)
	case string:
		*e = Color(s)
	default:
		return fmt.Errorf("unsupported scan type for Color: %T", src)
	}
	return nil
}

type NoteType string

const (
	NoteTypeNote NoteType = "note"
	NoteTypeList NoteType = "list"
)

func (e *NoteType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = NoteType(s)
	case string:
		*e = NoteType(s)
	default:
		return fmt.Errorf("unsupported scan type for NoteType: %T", src)
	}
	return nil
}

type Note struct {
	ID         int32          `json:"id"`
	UserID     int32          `json:"user_id"`
	Title      sql.NullString `json:"title"`
	Color      Color          `json:"color"`
	Type       NoteType       `json:"type"`
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
