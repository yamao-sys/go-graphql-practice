package lib

import (
	"app/graph"
	"app/graph/generated"
	"app/lib/auth"
	"app/services"
	"app/view"
	"context"
	"errors"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func GetGraphQLServer(authService services.AuthService) *handler.Server {
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graph.NewResolver(authService)}))

	srv.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)

		var errorCode int64
		var error error

		var re view.ViewError
		if errors.As(e, &re) {
			errorCode = re.Code
			error = re.Message
		}

		err.Extensions = map[string]interface{}{
			"code":  errorCode,
			"error": error,
		}

		return err
	})

	return srv
}

func GetGraphQLHttpHandler(srv *handler.Server) http.Handler {
	graphSrv := graph.Middleware(srv)
	return auth.Middleware(graphSrv)
}
