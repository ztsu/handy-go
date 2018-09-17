package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Request is a type alias for events.APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

// Response is a type alias for events.APIGatewayProxyResponse
type Response events.APIGatewayProxyResponse

type deck struct {
	ID          uint
	Name        string
	TypeOfCards string
}

var cors = map[string]string{
	"Access-Control-Allow-Origin": "*",
}

// HandleRequest handles the Request req
func HandleRequest(ctx context.Context, req Request) (Response, error) {
	decks := make([]deck, 2)

	decks[0] = deck{
		ID:          1,
		Name:        "Test",
		TypeOfCards: "Words",
	}

	decks[1] = deck{
		ID:          2,
		Name:        "Yes, English can be weird. It can be understood",
		TypeOfCards: "Words",
	}

	b, err := json.Marshal(decks)
	if err != nil {
		return Response{StatusCode: 500}, nil
	}

	return Response{StatusCode: 200, Body: string(b), Headers: cors}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
