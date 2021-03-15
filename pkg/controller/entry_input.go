package controller

import (
	"time"

	"github.com/google/uuid"
	"github.com/rpagliuca/serverless-book-reading-tracker/pkg/entity"
)

type EntryInput struct {
	BookID        *string `json:"book_id"`
	StartTime     *string `json:"start_time"`
	EndTime       *string `json:"end_time"`
	StartLocation *int64  `json:"start_location"`
	EndLocation   *int64  `json:"end_location"`
	Version       *int64  `json:"version"`
}

func NewEntryFromEntryInput(username string, e EntryInput) entity.Entry {
	id, err := uuid.NewUUID()
	if err != nil {
		id = uuid.New()
	}
	idString := id.String()
	version := int64(1)
	if e.Version != nil {
		version = *e.Version
	}
	return entity.Entry{
		Username:      stringPointer(username),
		ID:            &idString,
		BookID:        e.BookID,
		StartTime:     toTime(e.StartTime),
		EndTime:       toTime(e.EndTime),
		StartLocation: e.StartLocation,
		EndLocation:   e.EndLocation,
		DateCreated:   timePointer(time.Now()),
		DateModified:  timePointer(time.Now()),
		Version:       &version,
	}
}

func toTime(str *string) *time.Time {
	if str == nil {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *str)
	if err != nil {
		return nil
	}
	return &t
}

func stringPointer(str string) *string {
	return &str
}

func timePointer(time time.Time) *time.Time {
	return &time
}
