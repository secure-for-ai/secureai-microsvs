package graphql

import (
	"github.com/graphql-go/graphql"
	"template2/test_app/constant"
	"template2/test_app/model"
)

var idArgs = graphql.FieldConfigArgument{
	"uid": &graphql.ArgumentConfig{
		Description: "uid",
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
		Type:        graphql.Int,
		Description: "creation time",
	},
	"updateTime": &graphql.ArgumentConfig{
		// change to Int64 scalar
		Type:        graphql.Int,
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
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if manager, ok := p.Source.(*model.UserInfo); ok {
					return manager.UID.Hex(), nil
				}
				return nil, constant.ErrParamEmpty
			},
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
			Type:        graphql.Int,
			Description: "creation time",
		},
		"updateTime": &graphql.Field{
			// change to Int64 scalar
			Type:        graphql.Int,
			Description: "update time",
		},
	},
})
