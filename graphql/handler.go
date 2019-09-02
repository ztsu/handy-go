package graphql

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"net/http"
)

func Handler() http.HandlerFunc {
	return handler.GraphQL(
		NewExecutableSchema(
			Config{
				Resolvers: &Resolver{},
			},
		),
		handler.IntrospectionEnabled(true),
		handler.ResolverMiddleware(LogMutationName),
	)
}

func LogMutationName(ctx context.Context, next graphql.Resolver) (res interface{}, err error) {
	//if rCtx := graphql.GetResolverContext(ctx); rCtx != nil {
	//	if rCtx.IsMethod {
	//		fmt.Printf("Mutation name: %#v\n", rCtx.Field.Name)
	//	}
	//}

	return next(ctx)
}
