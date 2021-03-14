package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

var entriesTable = os.Getenv("ENTRYS_TABLE")

func listOneEntry(UUID string) interface{} {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Read
	query := &dynamodb.QueryInput{
		TableName:              aws.String(entriesTable),
		KeyConditionExpression: aws.String("key1 = :v AND key2 = :u"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v": {
				S: aws.String("rafpag"),
			},
			":u": {
				S: aws.String(UUID),
			},
		},
	}

	result, err := svc.Query(query)

	if err != nil {
		return err.Error()
	}

	list := []Entity{}
	for _, v := range result.Items {
		list = append(list, NewEntityFromItem(v))
	}

	return list
}

func listEntries() interface{} {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Read
	query := &dynamodb.QueryInput{
		TableName:              aws.String(entriesTable),
		KeyConditionExpression: aws.String("key1 = :v"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":v": {
				S: aws.String("rafpag"),
			},
		},
	}
	result, err := svc.Query(query)

	if err != nil {
		return err.Error()
	}

	list := []Entity{}
	for _, v := range result.Items {
		list = append(list, NewEntityFromItem(v))
	}

	return list
}

type Entity struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	Other    string `json:"other"`
}

func NewEntityFromItem(m map[string]*dynamodb.AttributeValue) Entity {
	return Entity{
		Username: *m["key1"].S,
		UUID:     *m["key2"].S,
		Other:    *m["other"].S,
	}
}

func insertEntry() interface{} {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Write
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"key1": {
				S: aws.String("rafpag"),
			},
			"key2": {
				S: aws.String(uuid.New().String()),
			},
			"other": {
				S: aws.String(uuid.New().String()),
			},
		},
		TableName: aws.String(entriesTable),
	}
	_, err := svc.PutItem(input)
	if err != nil {
		return err.Error()
	}

	return true
}

func listTables() []string {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSession())

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	// create the input configuration instance
	input := &dynamodb.ListTablesInput{}

	tables := []string{}

	for {
		// Get the list of tables
		result, err := svc.ListTables(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeInternalServerError:
					fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return []string{}
		}

		for _, n := range result.TableNames {
			tables = append(tables, *n)
		}

		// assign the last read tablename as the start for our next call to the ListTables function
		// the maximum number of table names returned in a call is 100 (default), which requires us to make
		// multiple calls to the ListTables function to retrieve all table names
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}
	return tables
}
