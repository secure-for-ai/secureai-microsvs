package graphql

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/graphql-go/graphql"
	gh "github.com/graphql-go/handler"
	"net/http"
	"template2/lib/session"
	"template2/test_app/config"
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
			"getUser": &graphql.Field{
				Args:        userArgs,
				Type:        userType,
				Description: "get user info",
				Resolve:     getUser,
			},
			"listUser": &graphql.Field{
				Args:        userListArgs,
				Type:        userListType,
				Description: "list users",
				Resolve:     listUser,
			},
			"checkLogin": &graphql.Field{
				Args:        userArgs,
				Type:        userType,
				Description: "get user info",
				Resolve:     checkLogin,
			},
		},
	})

	mutation = graphql.NewObject(graphql.ObjectConfig{
		Name: "mutation",
		Fields: graphql.Fields{
			"health": &graphql.Field{
				Type:        graphql.String,
				Description: "health check",
				Resolve:     getHealth,
			},
			"createUser": &graphql.Field{
				Args:        userArgs,
				Type:        graphql.Boolean,
				Description: "create user",
				Resolve:     createUser,
			},
			"updateUser": &graphql.Field{
				Args:        userArgs,
				Type:        graphql.Boolean,
				Description: "update account",
				Resolve:     updateUser,
			},
			"deleteUser": &graphql.Field{
				Args:        idArgs,
				Type:        graphql.Boolean,
				Description: "delete account",
				Resolve:     deleteUser,
			},
			"login": &graphql.Field{
				Args:        usernameArgs,
				Type:        graphql.Boolean,
				Description: "login account",
				Resolve:     login,
			},
			"logout": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "logout account",
				Resolve:     logout,
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
		ResultCallbackFn: func(ctx context.Context, params *graphql.Params, result *graphql.Result, responseBody []byte) {
			//sess, _ := ctx.Value("sessions").(map[string]*sessions.Session)
			//if sess["SID"].IsNew {
			//	fmt.Println("save session *********")
			//	sess["SID"].Save(r, w)
			//}
		},
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
	sess, err := config.SessionStore.Get(r, "SID")

	if err != nil {
		return
	}

	sessionMap := map[string]*sessions.Session{
		"SID": sess,
	}

	collection := session.NewCollection(r, w, sessionMap)
	ctx := session.NewCollectionContext(context.Background(), collection)

	//ctx := context.WithValue(context.Background(), constant.JWTContextKey, user)
	//key := "sessions"
	//ctx := context.WithValue(context.Background(), &key, sesss)
	handler.ContextHandler(ctx, w, r)

	//if sess.IsNew {
	//	fmt.Println("save session *********")
	//	sess.Save(r, w)
	//}
}
