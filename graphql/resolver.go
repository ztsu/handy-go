package graphql

import (
	"context"
	"github.com/matoous/go-nanoid"
	"github.com/ztsu/handy-go/store"
	"github.com/ztsu/handy-go/store/api"
	"net/http"
)

type Resolver struct {
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct {
	*Resolver
}

func (r *queryResolver) Version(ctx context.Context) (string, error) {
	return "0", nil
}

type mutationResolver struct {
	*Resolver
}

func (r *mutationResolver) RegisterUser(ctx context.Context, input RegisterUser) (*User, error) {
	client := api.NewClient(http.DefaultClient, "http://localhost:8080")

	id, err := gonanoid.Nanoid()
	if err != nil {
		return nil, err
	}

	user := &store.User{
		ID:    id,
		Email: input.Email,
	}

	err = client.AddUser(user)
	if err != nil {
		return nil, err
	}

	return &User{ID: user.ID, Email: user.Email}, nil
}
