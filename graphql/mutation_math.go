package graphql

import "context"

type mathMutationResolver struct {
	*Resolver
}

func (r *mathMutationResolver) Inc(ctx context.Context, obj *MathMutation, d int) (int, error) {
	return d + 1, nil
}

func (r *mathMutationResolver) Dec(ctx context.Context, obj *MathMutation, d int) (int, error) {
	return d - 1, nil
}