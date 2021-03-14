package main

import (
	"context"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response = events.APIGatewayProxyResponse
type Request = events.APIGatewayProxyRequest

var router *lmdrouter.Router

func init() {
	router = lmdrouter.NewRouter("/entries", authMiddleware)
	router.Route("GET", "/", listHandler)
	router.Route("POST", "/", insertHandler)
	router.Route("GET", "/:uuid", getOneHandler)
}

// Allow mocking
var lambdaStart = lambda.Start

func main() {
	lambdaStart(router.Handler)
}

func authMiddleware(next lmdrouter.Handler) lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (
		res events.APIGatewayProxyResponse,
		err error,
	) {
		res, err = next(ctx, req)
		return res, err
	}
}

func insertHandler(ctx context.Context, req events.APIGatewayProxyRequest) (
	res events.APIGatewayProxyResponse,
	err error,
) {
	return createSuccessResponse(insertEntry())
}

func listHandler(ctx context.Context, req events.APIGatewayProxyRequest) (
	res events.APIGatewayProxyResponse,
	err error,
) {
	return createSuccessResponse(listEntries())
}

type getOneInput struct {
	UUID string `lambda:"path.uuid"` // a path parameter declared as :id
	//ShowSomething     bool     `lambda:"query.show_something"`   // a query parameter named "show_something"
	//AcceptedLanguages []string `lambda:"header.Accept-Language"` // a multi-value header parameter
}

func getOneHandler(ctx context.Context, req events.APIGatewayProxyRequest) (
	res events.APIGatewayProxyResponse,
	err error,
) {
	var input getOneInput
	err = lmdrouter.UnmarshalRequest(req, false, &input)
	if err != nil {
		return lmdrouter.HandleError(err)
	}

	return createSuccessResponse(listOneEntry(input.UUID))
}
