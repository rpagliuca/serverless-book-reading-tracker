package entity

import "time"

type RecordNotFound error
type MoreThanOneRecordFound error

type Property int

const (
	Username Property = iota
	ID
	BookID
	StartTime
	EndTime
	StartLocation
	EndLocation
	DateCreated
	DateModified
	Version
)

type Entry struct {
	Username      *string    `json:"username"`
	ID            *string    `json:"id"`
	BookID        *string    `json:"book_id"`
	StartTime     *time.Time `json:"start_time"`
	EndTime       *time.Time `json:"end_time"`
	StartLocation *int64     `json:"start_location"`
	EndLocation   *int64     `json:"end_location"`
	DateCreated   *time.Time `json:"date_created"`
	DateModified  *time.Time `json:"date_modified"`
	Version       *int64     `json:"version"`
}
