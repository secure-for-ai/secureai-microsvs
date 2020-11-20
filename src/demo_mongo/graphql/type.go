package graphql

import (
	"github.com/graphql-go/graphql"
	"template2/lib/graphql_ext"
)

var idArgs = graphql.FieldConfigArgument{
	"uid": &graphql.ArgumentConfig{
		Description: "uid",
		Type:        graphql.String,
	},
}

var usernameArgs = graphql.FieldConfigArgument{
	"username": &graphql.ArgumentConfig{
		Description: "username",
		Type:        graphql.String,
	},
}

var userArgs = graphql.FieldConfigArgument{
	"uid": &graphql.ArgumentConfig{
		Type:        graphql.String,
		Description: "id",
	},
	"username": &graphql.ArgumentConfig{
		Type:        graphql.String,
		Description: "username",
	},
	"nickname": &graphql.ArgumentConfig{
		Type:        graphql.String,
		Description: "nickname",
	},
	"email": &graphql.ArgumentConfig{
		Type:        graphql.String,
		Description: "email",
	},
	"createTime": &graphql.ArgumentConfig{
		// change to Int64 scalar
		Type:        graphql_ext.Int64,
		Description: "creation time",
	},
	"updateTime": &graphql.ArgumentConfig{
		// change to Int64 scalar
		Type:        graphql_ext.Int64,
		Description: "update time",
	},
}

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "user",
	Description: "user",
	Fields: graphql.Fields{
		"uid": &graphql.Field{
			Type:        graphql.ID,
			Description: "user id",
			//Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			//	if user, ok := p.Source.(*model.UserInfo); ok {
			//		return user.UID, nil
			//	}
			//	if user, ok := p.Source.(model.UserInfo); ok {
			//		return user.UID, nil
			//	}
			//	return nil, constant.ErrParamEmpty
			//},
		},
		"username": &graphql.Field{
			Type:        graphql.String,
			Description: "username",
		},
		"nickname": &graphql.Field{
			Type:        graphql.String,
			Description: "nickname",
		},
		"email": &graphql.Field{
			Type:        graphql.String,
			Description: "email",
		},
		"createTime": &graphql.Field{
			// change to Int64 scalar
			Type:        graphql_ext.Int64,
			Description: "creation time",
		},
		"updateTime": &graphql.Field{
			// change to Int64 scalar
			Type:        graphql_ext.Int64,
			Description: "update time",
		},
	},
})

var userListType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "userList",
	Description: "user list",
	Fields: graphql.Fields{
		"list": &graphql.Field{
			Type:        graphql.NewList(userType),
			Description: "user list",
		},
		"count": &graphql.Field{
			Type:        graphql_ext.Int64,
			Description: "number of users",
		},
	},
})

var userListArgs = graphql.FieldConfigArgument{
	"username": &graphql.ArgumentConfig{
		Type:        graphql.String,
		Description: "username",
	},
	"perPage": &graphql.ArgumentConfig{
		Type:        graphql.NewNonNull(graphql_ext.Int64),
		Description: "perPage",
	},
	"page": &graphql.ArgumentConfig{
		Type:        graphql.NewNonNull(graphql_ext.Int64),
		Description: "perPage",
	},
}
