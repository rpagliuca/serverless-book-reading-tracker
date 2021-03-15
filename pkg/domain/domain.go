package domain

import (
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

func InsertEntry(e entity.Entry) error {
	err := persistence.InsertEntry(e)
	return err
}

func DeleteOneEntry(username, UUID string) error {
	err := persistence.DeleteOneEntry(username, UUID)
	return err
}
