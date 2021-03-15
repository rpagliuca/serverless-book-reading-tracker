package main

import (
	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rpagliuca/serverless-book-reading-tracker/pkg/auth"
	"github.com/rpagliuca/serverless-book-reading-tracker/pkg/controller"
)

type Response = events.APIGatewayProxyResponse
type Request = events.APIGatewayProxyRequest

var router *lmdrouter.Router

func init() {
	router = lmdrouter.NewRouter("")
	router.Route("GET", "/entries/", controller.ListHandler, auth.AuthMiddleware)
	router.Route("POST", "/entries/", controller.InsertHandler, auth.AuthMiddleware)
	router.Route("GET", "/entries/:uuid", controller.GetOneHandler, auth.AuthMiddleware)
	router.Route("DELETE", "/entries/:uuid", controller.DeleteOneHandler, auth.AuthMiddleware)
	router.Route("PATCH", "/entries/:uuid", controller.PatchOneHandler, auth.AuthMiddleware)
	router.Route("POST", "/login", controller.LoginHandler)
}

// Allow mocking
var lambdaStart = lambda.Start

func main() {
	lambdaStart(router.Handler)
}
