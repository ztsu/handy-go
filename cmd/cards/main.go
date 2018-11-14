package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/ztsu/handy-go/internal/handy"
)

// Request is a type alias for events.APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

// Response is a type alias for events.APIGatewayProxyResponse
type Response events.APIGatewayProxyResponse

var cors = map[string]string{
	"Access-Control-Allow-Origin": "*",
}

// HandleRequest handles the Request req
func HandleRequest(ctx context.Context, req Request) (Response, error) {

	cards := handy.SampleCards

	b, err := json.Marshal(cards)
	if err != nil {
		return Response{StatusCode: 500}, nil
	}

	return Response{StatusCode: 200, Body: string(b), Headers: cors}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
