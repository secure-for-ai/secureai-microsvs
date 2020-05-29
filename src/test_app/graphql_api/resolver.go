package graphql_api

import (
	"github.com/graphql-go/graphql"
	"log"
	"template2/test_app/constant"
	"template2/test_app/model"
)

func getUser(p graphql.ResolveParams) (interface{}, error) {
	var (
		ok       bool
		username string
	)
	if username, ok = p.Args["username"].(string); !ok {
		return false, constant.ErrParamEmpty
	}

	return model.GetUser(username)
	/*if password, _ = p.Args["password"].(string); password == "" {
		return model.GetUser(username)
	} else {
		return model.ManagerLogin(username, password)
	}*/
}

func createUser(p graphql.ResolveParams) (interface{}, error) {
	var (
		ok   bool
		err  error
		user = new(model.UserInfo)
	)

	if user.Username, ok = p.Args["username"].(string); !ok {
		return false, constant.ErrParamEmpty
	}
	if user.Nickname, ok = p.Args["nickname"].(string); !ok {
		return false, constant.ErrParamEmpty
	}
	if user.Email, ok = p.Args["email"].(string); !ok {
		return false, constant.ErrParamEmpty
	}

	if err = model.CreateUser(user); err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}
