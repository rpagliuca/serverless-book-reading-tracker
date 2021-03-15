package persistence

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/rpagliuca/serverless-book-reading-tracker/pkg/entity"
)

func ListOneEntry(username, UUID string) (entity.Entry, error) {

	if username == "" || UUID == "" {
		return entity.Entry{}, errors.New(ERROR_REQUIRED)
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Read
	query := &dynamodb.QueryInput{
		TableName:              aws.String(entriesTable),
		KeyConditionExpression: aws.String("username = :username AND id = :id"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {
				S: aws.String(username),
			},
			":id": {
				S: aws.String(UUID),
			},
		},
	}

	result, err := svc.Query(query)

	if err != nil {
		return entity.Entry{}, err
	}

	if len(result.Items) > 1 {
		return entity.Entry{}, entity.MoreThanOneRecordFound(errors.New("Query should return only 1 result, but more results were found"))
	}

	if len(result.Items) == 0 {
		return entity.Entry{}, entity.RecordNotFound(errors.New("Entry not found"))
	}

	entry := NewEntryFromItem(result.Items[0])

	return entry, nil
}

func ListEntries(username string) ([]entity.Entry, error) {

	if username == "" {
		return []entity.Entry{}, errors.New(ERROR_REQUIRED)
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Read
	query := &dynamodb.QueryInput{
		TableName:              aws.String(entriesTable),
		KeyConditionExpression: aws.String("username = :username"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":username": {
				S: aws.String(username),
			},
		},
	}
	result, err := svc.Query(query)

	if err != nil {
		return []entity.Entry{}, err
	}

	list := []entity.Entry{}
	for _, v := range result.Items {
		list = append(list, NewEntryFromItem(v))
	}

	return list, nil
}

func NewEntryFromItem(m map[string]*dynamodb.AttributeValue) entity.Entry {
	return entity.Entry{
		ID:            m["id"].S,
		Username:      m["username"].S,
		BookID:        m["book_id"].S,
		StartTime:     toTime(m["start_time"].S),
		EndTime:       toTime(m["end_time"].S),
		StartLocation: toInt(m["start_location"].S),
		EndLocation:   toInt(m["end_location"].S),
		DateCreated:   toTime(m["date_created"].S),
		DateModified:  toTime(m["date_modified"].S),
		Version:       toInt(m["version"].S),
	}
}

func InsertEntry(entry entity.Entry) error {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	fmt.Printf("Inserting %+v\n", entry)

	// Write
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"username":       str(entry.Username),
			"id":             str(entry.ID),
			"book_id":        str(entry.BookID),
			"start_time":     str(entry.StartTime),
			"end_time":       str(entry.EndTime),
			"start_location": str(entry.StartLocation),
			"end_location":   str(entry.EndLocation),
			"date_created":   str(entry.DateCreated),
			"date_modified":  str(entry.DateModified),
			"version":        str(entry.Version),
		},
		TableName: aws.String(entriesTable),
	}
	_, err := svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func DeleteOneEntry(username, UUID string) error {

	if username == "" {
		return errors.New(ERROR_REQUIRED)
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Read
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(entriesTable),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(username),
			},
			"id": {
				S: aws.String(UUID),
			},
		},
	}

	_, err := svc.DeleteItem(input)

	if err != nil {
		return err
	}

	return nil
}
