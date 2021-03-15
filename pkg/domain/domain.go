package domain

import (
	"errors"
	"time"

	"github.com/rpagliuca/serverless-book-reading-tracker/pkg/entity"
	"github.com/rpagliuca/serverless-book-reading-tracker/pkg/persistence"
)

func ListOneEntry(username, UUID string) (entity.Entry, error) {
	entry, err := persistence.ListOneEntry(username, UUID)
	return entry, err
}

func ListEntries(username string) ([]entity.Entry, error) {
	entries, err := persistence.ListEntries(username)
	return entries, err
}

func InsertEntry(username string, e entity.Entry) error {
	e.Username = &username
	now := time.Now()
	e.DateCreated = &now
	e.DateModified = &now
	version := int64(1)
	e.Version = &version
	err := persistence.InsertEntry(e)
	return err
}

func DeleteOneEntry(username, UUID string) error {
	err := persistence.DeleteOneEntry(username, UUID)
	return err
}

func PatchEntry(username string, e entity.Entry) error {
	e.Username = &username
	now := time.Now()
	e.DateModified = &now
	p := []entity.Property{}
	if e.Version == nil {
		return errors.New("Property 'version' must be incremented for PATCH operations")
	}
	if e.StartLocation != nil {
		p = append(p, entity.StartLocation)
	}
	if e.StartTime != nil {
		p = append(p, entity.StartTime)
	}
	if e.EndLocation != nil {
		p = append(p, entity.EndLocation)
	}
	if e.EndTime != nil {
		p = append(p, entity.EndTime)
	}
	if e.BookID != nil {
		p = append(p, entity.BookID)
	}
	p = append(p, entity.Version)
	p = append(p, entity.DateModified)
	// At least one 2 properties must be patched to be valid (version plus others)
	if len(p) < 3 {
		return errors.New("At least one other property besides 'version' and 'date_modified' must be patched")
	}
	err := persistence.PatchEntry(e, p)
	return err
}
