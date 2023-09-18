package model

import (
	"errors"
	"github/be/database"

	"gorm.io/gorm"
)

type EntryError struct {
	Err error
}

func (e *EntryError) Error() string {
	return e.Err.Error()
}

type Entry struct {
	gorm.Model
	Content string `gorm:"type:text" json:"content"`
	UserID  uint
}

func (entry *Entry) Save() (*Entry, error) {
	if entry.Content == "" {
		return &Entry{}, &EntryError{Err: errors.New("content is required")}
	}

	if err := database.Database.Create(&entry).Error; err != nil {
		return &Entry{}, &EntryError{Err: err}
	}

	return entry, nil
}
