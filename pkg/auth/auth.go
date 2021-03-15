package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/aquasecurity/lmdrouter"
	"github.com/aws/aws-lambda-go/events"
)

type key string

var keyUser = key("authUser")

func UserFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(keyUser).(string); ok {
		return v
	}
	return ""
}

func AuthMiddleware(next lmdrouter.Handler) lmdrouter.Handler {
	return func(ctx context.Context, req events.APIGatewayProxyRequest) (
		res events.APIGatewayProxyResponse,
		err error,
	) {
		token, err := extractToken(req.Headers["Authorization"])
		if err != nil {
			return createErrorResponse(401, err)
		}
		isValid, username := ParseToken(token)
		if !isValid || username == "" {
			return createErrorResponse(401, errors.New("Invalid token"))
		}
		ctx2 := context.WithValue(ctx, keyUser, username)
		res, err = next(ctx2, req)
		return res, err
	}
}

func extractToken(authHeader string) (string, error) {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", errors.New("Expecting header 'Authorization: Bearer <TOKEN>'")
	}
	if parts[0] != "Bearer" {
		return "", errors.New("Expecting header 'Authorization: Bearer <TOKEN>'")
	}
	return parts[1], nil
}

var DURATION_VALID_SECONDS = 60 * 60
var TOKEN = "YABADABADOO"

func GetToken(user, password string) (token string, timestamp int) {
	expiration := int(time.Now().Unix()) + DURATION_VALID_SECONDS
	return TOKEN, expiration
}

func ParseToken(token string) (isValid bool, username string) {
	if token == TOKEN {
		return true, "rafpag"
	}
	return false, ""
}

func createErrorResponse(statusCode int, err error) (Response, error) {
	body, err := json.Marshal(map[string]interface{}{
		"success": false,
		"message": err.Error(),
	})
	if err != nil {
		return Response{StatusCode: 500}, err
	}
	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      statusCode,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	return resp, nil
}

type Response = events.APIGatewayProxyResponse
type Request = events.APIGatewayProxyRequest
