package gq

import (
	"context"
	"github.com/graphql-go/graphql"
	gh "github.com/graphql-go/handler"
	"net/http"
)

var (
	handler *gh.Handler

	query = graphql.NewObject(graphql.ObjectConfig{
		Name: "query",
		Fields: graphql.Fields{
			"health": &graphql.Field{
				Type:        graphql.String,
				Description: "health check",
				Resolve:     getHealth,
			},
		},
	})

	mutation = graphql.NewObject(graphql.ObjectConfig{
		Name: "query",
		Fields: graphql.Fields{
			"health": &graphql.Field{
				Type:        graphql.String,
				Description: "health check",
				Resolve:     getHealth,
			},
		},
	})
)

func init() {
	schemaConfig := graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
	}

	schema, _ := graphql.NewSchema(schemaConfig)

	isProd := false
	handler = gh.New(&gh.Config{
		Schema: &schema,
		// GraphiQL: !isProd,
		Pretty:     !isProd,
		Playground: !isProd,
	})
}

func Graphql(w http.ResponseWriter, r *http.Request) {
	/* jwt */
	/*token := r.Header.Get("Authorization")
	  user, ok := validateJWT(token)
	  if !ok && isProd {
	      resJSONError(w, http.StatusUnauthorized, constant.ErrorMsgUnAuth)
	      return
	  }*/

	//ctx := context.WithValue(context.Background(), constant.JWTContextKey, user)
	ctx := context.WithValue(context.Background(), "user", "key")
	handler.ContextHandler(ctx, w, r)
}
