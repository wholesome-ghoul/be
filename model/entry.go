package model

import (
	"errors"
	"github/be/database"
)

type EntryError struct {
	Err error
}

func (e *EntryError) Error() string {
	return e.Err.Error()
}

type Entry struct {
	Content string `json:"content"`
	UserID  uint
}

func (entry *Entry) Save() (*Entry, error) {
	if entry.Content == "" {
		return &Entry{}, &EntryError{Err: errors.New("content is required")}
	}

	if _, err := database.Database.Exec("INSERT INTO entries (content, user_id) VALUES ($1, $2)", entry.Content, entry.UserID); err != nil {
		return &Entry{}, &EntryError{Err: errors.New("couldn't save entry")}
	}

	return entry, nil
}

func GetAllEntries(userID uint) ([]Entry, error) {
	var entries []Entry
	rows, err := database.Database.Query("SELECT content, user_id FROM entries WHERE user_id = $1", userID)
	if err != nil {
		return entries, &EntryError{Err: errors.New("couldn't get entries")}
	}

	for rows.Next() {
		var entry Entry
		if err := rows.Scan(&entry.Content, &entry.UserID); err != nil {
			return entries, &EntryError{Err: errors.New("couldn't get entries")}
		}

		entries = append(entries, entry)
	}

	return entries, nil
}
