package graphql

import "context"



type queryResolver struct {
	*Resolver
}

func (r *queryResolver) Version(ctx context.Context) (string, error) {
	return "0", nil
}

