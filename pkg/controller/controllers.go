package controller

import (
	"context"
	"encoding/json"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/rpagliuca/serverless-book-reading-tracker/pkg/auth"
	"github.com/rpagliuca/serverless-book-reading-tracker/pkg/domain"
	"github.com/rpagliuca/serverless-book-reading-tracker/pkg/entity"
)

func InsertHandler(ctx context.Context, req events.APIGatewayProxyRequest) (
	res events.APIGatewayProxyResponse,
	err error,
) {
	var entryInput EntryInput
	err = json.Unmarshal([]byte(req.Body), &entryInput)
	if err != nil {
		return createErrorResponse(500, err)
	}
	entry := NewEntryFromEntryInput(auth.UserFromContext(ctx), entryInput)
	err = domain.InsertEntry(entry)
	if err != nil {
		return createErrorResponse(500, err)
	}
	resp := map[string]interface{}{
		"success": true,
		"message": "OK",
		"id":      *entry.ID,
	}
	return createSuccessResponse(resp)
}

func ListHandler(ctx context.Context, req events.APIGatewayProxyRequest) (
	res events.APIGatewayProxyResponse,
	err error,
) {
	entries, err := domain.ListEntries(auth.UserFromContext(ctx))
	if err != nil {
		return createErrorResponse(500, err)
	}
	return createSuccessResponse(entries)
}

type getOneInput struct {
	UUID string `lambda:"path.uuid"` // a path parameter declared as :id
	//ShowSomething     bool     `lambda:"query.show_something"`   // a query parameter named "show_something"
	//AcceptedLanguages []string `lambda:"header.Accept-Language"` // a multi-value header parameter
}

func GetOneHandler(ctx context.Context, req events.APIGatewayProxyRequest) (
	res events.APIGatewayProxyResponse,
	err error,
) {
	var input getOneInput
	err = lmdrouter.UnmarshalRequest(req, false, &input)
	if err != nil {
		return lmdrouter.HandleError(err)
	}
	entry, err := domain.ListOneEntry(auth.UserFromContext(ctx), input.UUID)
	if err != nil {
		status := 500
		switch err.(type) {
		case entity.RecordNotFound:
			status = 404
		}
		return createErrorResponse(status, err)
	}
	return createSuccessResponse(entry)
}

func DeleteOneHandler(ctx context.Context, req events.APIGatewayProxyRequest) (
	res events.APIGatewayProxyResponse,
	err error,
) {
	var input getOneInput
	err = lmdrouter.UnmarshalRequest(req, false, &input)
	if err != nil {
		return lmdrouter.HandleError(err)
	}
	err = domain.DeleteOneEntry(auth.UserFromContext(ctx), input.UUID)
	if err != nil {
		status := 500
		switch err.(type) {
		case entity.RecordNotFound:
			status = 404
		}
		return createErrorResponse(status, err)
	}
	return createSuccessResponse(basicSuccessResponse)
}

var basicSuccessResponse = map[string]interface{}{
	"success": true,
	"message": "OK",
}

type loginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(ctx context.Context, req events.APIGatewayProxyRequest) (
	res events.APIGatewayProxyResponse,
	err error,
) {
	var input loginInput
	err = json.Unmarshal([]byte(req.Body), &input)
	if err != nil {
		return createErrorResponse(500, err)
	}
	token, expiration := auth.GetToken(input.Username, input.Password)
	resp := map[string]interface{}{
		"success":              true,
		"message":              "OK",
		"token":                token,
		"expiration_timestamp": expiration,
	}
	return createSuccessResponse(resp)
}