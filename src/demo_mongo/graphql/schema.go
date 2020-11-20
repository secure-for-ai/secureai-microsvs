package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/graphql-go/graphql"
	gh "github.com/graphql-go/handler"
	"log"
	"net/http"
	"template2/demo_mongo/config"
	"template2/demo_mongo/model"
	"template2/lib/session"
	"template2/lib/util"
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

	isProd = true

	gqInitializer initializer = prodInitializer{}
)

type initializer interface {
	loadSession(w http.ResponseWriter, r *http.Request) *session.Collection
}

type prodInitializer struct{}

func (prodInitializer) loadSession(w http.ResponseWriter, r *http.Request) *session.Collection {
	log.Println("production")
	sess, err := config.SessionStore.Get(r, "SID")

	/*	if sess.ID == "" {
			returnVal := map[string]interface{}{
				"message":  "Not Authorized",
				"redirect": "http://localhost/login",
			}
			returnCode := 401

			returnStr, _ := json.Marshal(returnVal)
			http.Error(w, string(returnStr), returnCode)
		}
	*/
	// something wrong with backend redis and database,
	// sent alert for maintenance
	if err != nil {
		fmt.Println(err)
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
			fmt.Println(err)
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
		return nil
	}

	if sess.ID == "" {
		sess.Values["data"] = map[string]interface{}{}
	}

	sessionMap := map[string]*sessions.Session{
		"SID": sess,
	}

	collection := session.NewCollection(r, w, sessionMap)

	return collection
}

type devInitializer struct{}

func (devInitializer) loadSession(w http.ResponseWriter, r *http.Request) *session.Collection {
	log.Println("dev")
	sess, _ := config.SessionStore.Get(r, "SID")

	//sess.ID, _ = config.SessionStore.Storage//IdGenerator().Generate(config.SessionStore.IdLength())
	sess.Values["uid"] = int64(0)
	sess.Values["data"] = map[string]interface{}{
		"userInfo": model.UserInfo{
			UID:        0,
			Username:   "testUsername",
			Nickname:   "testNickname",
			Email:      "test@test.com",
			CreateTime: 1400000,
			UpdateTime: 1400000,
		},
	}

	sessionMap := map[string]*sessions.Session{
		"SID": sess,
	}

	collection := session.NewCollection(r, w, sessionMap)

	return collection

}

func init() {
	schemaConfig := graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
	}

	schema, _ := graphql.NewSchema(schemaConfig)

	switch config.Conf.AppInfo.Env {
	case util.AppEnvProd:
		gqInitializer = prodInitializer{}
		isProd = true
	case util.AppEnvTest:
		gqInitializer = prodInitializer{}
		isProd = false
	case util.AppEnvDev:
		gqInitializer = devInitializer{}
		isProd = false
	default:
		gqInitializer = devInitializer{}
		isProd = false
	}

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

			userInfo, ok := s.Values["data"].(map[string]interface{})["userInfo"]

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
	collection := gqInitializer.loadSession(w, r)

	if collection == nil {
		return
	}

	ctx := session.NewCollectionContext(context.Background(), collection)

	handler.ContextHandler(ctx, w, r)
}
