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


// Card is a ...
type Card struct {
	Word        string
	Translation string
	IPA         string
}

var cors = map[string]string{
	"Access-Control-Allow-Origin": "*",
}

// HandleRequest handles the Request req
func HandleRequest(ctx context.Context, req Request) (Response, error) {

	cards := make(map[int][]Card, 2)

	cards[1] = []Card{
		{"handy", "удобный", "ˈhændɪ"},
	}

	cards[2] = []Card{
		{"through", "через", "θruː"},
		{"tough", "жесткий", "tʌf"},
		{"thorough", "полный", "ˈθʌrə"},
		{"thought", "мысль", "θɔːt"},
		{"though", "хотя", "ðəʊ"},
	}

	b, err := json.Marshal(cards)
	if err != nil {
		return Response{StatusCode: 500}, nil
	}

	return Response{StatusCode: 200, Body: string(b), Headers: cors}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
