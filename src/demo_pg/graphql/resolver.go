package graphql

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"log"
	"template2/demo_pg/constant"
	"template2/demo_pg/model"
	"template2/lib/session"
	"template2/lib/util"
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
}

func createUser(p graphql.ResolveParams) (interface{}, error) {
	var (
		ok   bool
		err  error
		user model.UserInfo
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

	if err = model.CreateUser(&user); err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}

func updateUser(p graphql.ResolveParams) (interface{}, error) {
	var (
		ok   bool
		err  error
		id   string
		user *model.UserInfo
	)
	if id, ok = p.Args["uid"].(string); !ok {
		return false, constant.ErrParamEmpty
	}
	// query whether user exist
	if user, err = model.GetUserById(id); err != nil {
		return false, err
	}

	updateFlag := false
	if username, ok := p.Args["username"].(string); ok && (username != user.Username) {
		user.Username = username
		updateFlag = true
	}
	if nickname, ok := p.Args["nickname"].(string); ok && (nickname != user.Nickname) {
		user.Nickname = nickname
		updateFlag = true
	}
	if email, ok := p.Args["email"].(string); ok && (email != user.Email) {
		user.Email = email
		updateFlag = true
	}

	// we only run update if any attribute being changed
	if updateFlag {
		user.UpdateTime = util.GetNowTimestamp()
		if err = model.UpdateUser(user); err != nil {
			return false, err
		}
	}
	return true, nil
}

func deleteUser(p graphql.ResolveParams) (interface{}, error) {
	id, _ := p.Args["uid"].(string)

	if err := model.DeleteUser(id); err != nil {
		return false, err
	}
	return true, nil
}

func listUser(p graphql.ResolveParams) (interface{}, error) {
	var (
		ok                   bool
		username             string
		page, perPage, count int64
		err                  error
		us                   *[]model.UserInfo
	)

	username, _ = p.Args["username"].(string)
	if page, ok = p.Args["page"].(int64); !ok || page < 1 {
		return false, constant.ErrParamEmpty
	}
	if perPage, ok = p.Args["perPage"].(int64); !ok || perPage < 10 {
		return false, constant.ErrParamEmpty
	}

	if count, us, err = model.ListUser(username, page, perPage); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"list":  us,
		"count": count,
	}, nil
}

func login(p graphql.ResolveParams) (interface{}, error) {
	var (
		user *model.UserInfo
		err  error
	)
	username, _ := p.Args["username"].(string)
	source := p.Source.(map[string]interface{})
	collection, ok := source["session"].(*session.Collection)

	if !ok {
		return false, constant.ErrSession
	}
	if user, err = model.GetUser(username); err != nil {
		return false, err
	}

	_ = collection.UpdateValue("SID", "uid", user.UID)
	_ = collection.UpdateValue("SID", "data", map[string]interface{}{"userInfo": user})
	//_ = collection.MaxAge("SID", 0)
	err = collection.Save("SID")

	if err != nil {
		return false, err
	}
	return true, nil
}

func logout(p graphql.ResolveParams) (interface{}, error) {
	source := p.Source.(map[string]interface{})
	collection, ok := source["session"].(*session.Collection)

	if !ok {
		return false, constant.ErrSession
	}

	_ = collection.MaxAge("SID", -1)
	err := collection.Save("SID")

	if err != nil {
		return false, err
	}

	return true, nil
}

func checkLogin(p graphql.ResolveParams) (interface{}, error) {
	fmt.Println(p.Args)
	source := p.Source.(map[string]interface{})

	if userInfo, ok := source["userInfo"]; ok {
		fmt.Println("login uid", source["uid"])
		return userInfo, nil
	}

	return nil, constant.ErrAccountNotLogin
}
