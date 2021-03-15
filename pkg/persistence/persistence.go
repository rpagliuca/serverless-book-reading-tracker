package persistence

import (
	"errors"
	"fmt"
	"strings"

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

func strValue(m map[string]*dynamodb.AttributeValue, key string) *string {
	if v, ok := m[key]; ok {
		return v.S
	}
	var empty string
	return &empty
}

func NewEntryFromItem(m map[string]*dynamodb.AttributeValue) entity.Entry {
	return entity.Entry{
		ID:            strValue(m, "id"),
		Username:      strValue(m, "username"),
		BookID:        strValue(m, "book_id"),
		StartTime:     toTime(strValue(m, "start_time")),
		EndTime:       toTime(strValue(m, "end_time")),
		StartLocation: toInt(strValue(m, "start_location")),
		EndLocation:   toInt(strValue(m, "end_location")),
		DateCreated:   toTime(strValue(m, "date_created")),
		DateModified:  toTime(strValue(m, "date_modified")),
		Version:       toInt(strValue(m, "version")),
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

type Key string

const (
	KeyVersion       Key = "version"
	KeyStartLocation Key = "start_location"
	KeyStartTime     Key = "start_time"
	KeyEndLocation   Key = "end_location"
	KeyEndTime       Key = "end_time"
	KeyBookID        Key = "book_id"
	KeyDateModified  Key = "date_modified"
)

type KeyWithValue struct {
	Key   Key
	Value *dynamodb.AttributeValue
}

var propGetter = map[entity.Property]func(entity.Entry) KeyWithValue{
	entity.Version: func(e entity.Entry) KeyWithValue {
		return KeyWithValue{KeyVersion, str(e.Version)}
	},
	entity.StartLocation: func(e entity.Entry) KeyWithValue {
		return KeyWithValue{KeyStartLocation, str(e.StartLocation)}
	},
	entity.StartTime: func(e entity.Entry) KeyWithValue {
		return KeyWithValue{KeyStartTime, str(e.StartTime)}
	},
	entity.EndLocation: func(e entity.Entry) KeyWithValue {
		return KeyWithValue{KeyEndLocation, str(e.EndLocation)}
	},
	entity.EndTime: func(e entity.Entry) KeyWithValue {
		return KeyWithValue{KeyEndTime, str(e.EndTime)}
	},
	entity.BookID: func(e entity.Entry) KeyWithValue {
		return KeyWithValue{KeyBookID, str(e.BookID)}
	},
	entity.DateModified: func(e entity.Entry) KeyWithValue {
		return KeyWithValue{KeyDateModified, str(e.DateModified)}
	},
}

func PatchEntry(entry entity.Entry, patchedProperties []entity.Property) error {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	if len(patchedProperties) < 1 {
		return errors.New("At least one patched property is required")
	}

	updateExpressionParts := []string{}
	attributeValues := map[string]*dynamodb.AttributeValue{}
	for i, p := range patchedProperties {
		keyWithValue := propGetter[p](entry)
		uniqueKey := fmt.Sprintf(":value%d", i)
		attributeValues[uniqueKey] = keyWithValue.Value
		updateExpressionParts = append(updateExpressionParts, string(keyWithValue.Key)+"="+uniqueKey)
	}

	partsJoined := strings.Join(updateExpressionParts, ", ")

	updateExpression := "SET " + partsJoined

	// Update
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(entriesTable),
		Key: map[string]*dynamodb.AttributeValue{
			"username": str(entry.Username),
			"id":       str(entry.ID),
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: attributeValues}
	_, err := svc.UpdateItem(input)
	if err != nil {
		return err
	}

	return nil
}
