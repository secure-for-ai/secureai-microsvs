package graphql

import (
	"github.com/graphql-go/graphql"
)

func getHealth(p graphql.ResolveParams) (interface{}, error) {
	return "hello world", nil
}
