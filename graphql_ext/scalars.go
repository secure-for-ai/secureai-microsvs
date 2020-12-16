package graphql_ext

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"math"
	"strconv"
)

// As per the GraphQL Spec, Integers are only treated as valid when a valid
// 32-bit signed integer, providing the broadest support across platforms.
//
// n.b. JavaScript's integers are safe between -(2^53 - 1) and 2^53 - 1 because
// they are internally represented as IEEE 754 doubles.
func coerceInt64(value interface{}) interface{} {
	switch value := value.(type) {
	case bool:
		if value == true {
			return 1
		}
		return 0
	case *bool:
		if value == nil {
			return nil
		}
		return coerceInt64(*value)
	case int:
		return int64(value)
	case *int:
		if value == nil {
			return nil
		}
		return coerceInt64(*value)
	case int8:
		return int64(value)
	case *int8:
		if value == nil {
			return nil
		}
		return int64(*value)
	case int16:
		return int64(value)
	case *int16:
		if value == nil {
			return nil
		}
		return int64(*value)
	case int32:
		return int64(value)
	case *int32:
		if value == nil {
			return nil
		}
		return int64(*value)
	case int64:
		return value
	case *int64:
		if value == nil {
			return nil
		}
		return coerceInt64(*value)
	case uint:
		return int64(value)
	case *uint:
		if value == nil {
			return nil
		}
		return coerceInt64(*value)
	case uint8:
		return int64(value)
	case *uint8:
		if value == nil {
			return nil
		}
		return int64(*value)
	case uint16:
		return int64(value)
	case *uint16:
		if value == nil {
			return nil
		}
		return int64(*value)
	case uint32:
		return int64(value)
	case *uint32:
		if value == nil {
			return nil
		}
		return coerceInt64(*value)
	case uint64:
		if value > uint64(math.MaxInt64) {
			return nil
		}
		return int64(value)
	case *uint64:
		if value == nil {
			return nil
		}
		return coerceInt64(*value)
	case float32:
		if value < float32(math.MinInt64) || value > float32(math.MaxInt64) {
			return nil
		}
		return int64(value)
	case *float32:
		if value == nil {
			return nil
		}
		return coerceInt64(*value)
	case float64:
		if value < float64(math.MinInt64) || value > float64(math.MaxInt64) {
			return nil
		}
		return int64(value)
	case *float64:
		if value == nil {
			return nil
		}
		return coerceInt64(*value)
	case string:
		val, err := strconv.ParseFloat(value, 0)
		if err != nil {
			return nil
		}
		return coerceInt64(val)
	case *string:
		if value == nil {
			return nil
		}
		return coerceInt64(*value)
	}

	// If the value cannot be transformed into an int, return nil instead of '0'
	// to denote 'no integer found'
	return nil
}

// Int is the GraphQL Integer type definition.
var Int64 = graphql.NewScalar(graphql.ScalarConfig{
	Name: "Int64",
	Description: "The `Int` scalar type represents non-fractional signed whole numeric " +
		"values. Int can represent values between -(2^31) and 2^31 - 1. ",
	Serialize:  coerceInt64,
	ParseValue: coerceInt64,
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.IntValue:
			if intValue, err := strconv.ParseInt(valueAST.Value, 10, 64); err == nil {
				return intValue
			}
		}
		return nil
	},
})
