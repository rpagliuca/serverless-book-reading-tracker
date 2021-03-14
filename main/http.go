package main

import (
	"bytes"
	"encoding/json"
)

func createSuccessResponse(data interface{}) (Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return createErrorResponse(err)
	}
	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
		},
	}
	return resp, nil
}

func createErrorResponse(err error) (Response, error) {
	body, err := json.Marshal(map[string]interface{}{
		"success": false,
		"message": err.Error(),
	})
	if err != nil {
		return Response{StatusCode: 400}, err
	}
	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      400,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	return resp, nil
}
