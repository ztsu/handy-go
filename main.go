package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
)

type Request events.APIGatewayProxyRequest

type Response events.APIGatewayProxyResponse

type Deck struct {
	Name  string
	Cards []Card
}

type Card struct {
	Word        string
	Translation string
	IPA         string
}



func HandleRequest(ctx context.Context, event Request) (Response, error) {


	deck := Deck{
		Name: "",
		Cards: []Card{
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
