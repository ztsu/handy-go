package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
)

// Request is a type alias for events.APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

// Response is a type alias for events.APIGatewayProxyResponse
type Response events.APIGatewayProxyResponse

type deck struct {
	Name  string
	Cards []card
}

type card struct {
	Word        string
	Translation string
	IPA         string
}

// HandleRequest handles the Request req
func HandleRequest(ctx context.Context, req Request) (Response, error) {
	deck := deck{
		Name: "Test name 7",
		Cards: []card{
			{"handy", "удобный", "ˈhændɪ"},
		},
	}

	b, err := json.Marshal(deck)
	if err != nil {
		return Response{StatusCode: 500}, nil
	}

	return Response{StatusCode: 200, Body: string(b)}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
