package persistence

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const ERROR_REQUIRED = "Missing required parameters"

var entriesTable = os.Getenv("ENTRIES_TABLE")

func toInt(str *string) *int64 {
	if str == nil {
		return nil
	}
	i, err := strconv.ParseInt(*str, 10, 64)
	if err != nil {
		return nil
	}
	return &i
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

var trueValue = true
var nullAttribute = &dynamodb.AttributeValue{
	NULL: &trueValue,
}

func str(val interface{}) *dynamodb.AttributeValue {
	var str string
	switch v := val.(type) {
	case *string:
		if v == nil {
			return nullAttribute
		}
		str = *v
	case *int64:
		if v == nil {
			return nullAttribute
		}
		str = fmt.Sprintf("%d", *v)
	case *time.Time:
		if v == nil {
			return nullAttribute
		}
		str = v.Format(time.RFC3339)
	default:
		return nullAttribute
	}
	return &dynamodb.AttributeValue{
		S: aws.String(str),
	}
}
