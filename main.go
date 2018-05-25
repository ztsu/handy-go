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

var cors = map[string]string{
	"Access-Control-Allow-Origin": "*",
}

// HandleRequest handles the Request req
func HandleRequest(ctx context.Context, req Request) (Response, error) {

	decks := make([]deck, 2)

	decks[0] = deck{
		Name: "Test",
		Cards: []card{
			{"handy", "удобный", "ˈhændɪ"},
		},
	}

	decks[1] = deck{
		Name: "Yes, English can be weird. It can be understood",
		Cards: []card{
			{"through", "через", "θruː"},
			{"tough", "жесткий", "tʌf"},
			{"thorough", "полный", "ˈθʌrə"},
			{"thought", "мысль", "θɔːt"},
			{"though", "хотя", "ðəʊ"},

		},
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
