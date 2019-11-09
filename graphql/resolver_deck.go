package graphql

import (
	"context"
	"github.com/ztsu/handy-go/store/api"
	"net/http"
)

type deckResolver struct {
	*Resolver
}

func (deckResolver) User(ctx context.Context, obj *Deck) (*User, error) {
	client := api.NewClient(http.DefaultClient, "http://localhost:8080")

	user, err := client.GetUser(obj.User.ID)
	if err != nil {
		return nil, err
	}

	return &User{ID: user.ID, Email: user.Email}, nil
}

