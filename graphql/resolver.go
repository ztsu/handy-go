//go:generate go run github.com/99designs/gqlgen

package graphql

type Resolver struct {
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Deck() DeckResolver {
	return &deckResolver{r}
}