package graphql

import (
	"context"
	"encoding/json"
	"github.com/gorilla/sessions"
	"github.com/graphql-go/graphql"
	gh "github.com/graphql-go/handler"
	"net/http"
	"template2/lib/session"
	"template2/test_app/config"
	"template2/test_app/model"
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
		RootObjectFn: func(ctx context.Context, r *http.Request) map[string]interface{} {
			root := map[string]interface{}{}

			collection, _ := session.FromCollectionContext(ctx)

			root["session"] = collection
			s, _ := collection.Get("SID")

			userInfo, ok := s.Values["userInfo"]

			if !ok {
				return root
			}

			root["uid"] = userInfo.(model.UserInfo).UID
			root["userInfo"] = userInfo.(model.UserInfo)
			root["role"] = []string{"user", "manager", "admin"}
			return root

		},
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
	sess, err := config.SessionStore.Get(r, "SID")

	// something wrong with backend redis and database,
	// sent alert for maintenance
	if err != nil {
		var (
			returnCode int
			returnVal  map[string]interface{}
		)
		switch err {
		// get unexpected session, either malicious session.ID,
		// or session got lost in out database. In most cases,
		// it would be malicious cases.
		case session.ErrNil:
			fallthrough
		case session.ErrInvalidCookie:
			sess.Options.MaxAge = -1
			_ = sess.Save(r, w)

			returnVal = map[string]interface{}{
				"message":  "Not Authorized",
				"redirect": "http://localhost/login",
			}
			returnCode = 401

		case session.ErrStoreFail:
			fallthrough
		default:
			returnVal = map[string]interface{}{
				"message":  "Server Internal Error",
				"redirect": "http://localhost/Error/500",
			}
			returnCode = 500
		}

		returnStr, _ := json.Marshal(returnVal)
		http.Error(w, string(returnStr), returnCode)
		return
	}

	sessionMap := map[string]*sessions.Session{
		"SID": sess,
	}

	collection := session.NewCollection(r, w, sessionMap)
	ctx := session.NewCollectionContext(context.Background(), collection)

	handler.ContextHandler(ctx, w, r)
}
