package control

import (
	"notification-service/domain"
	"strconv"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

// TIME -------------------------

func serializeTime(value interface{}) interface{} {
	if v, ok := value.(time.Time); ok {
		return v.Format(domain.TimeFormat)
	}
	if v, ok := value.(*time.Time); ok {
		if v == nil {
			return ""
		}
		return v.Format(domain.TimeFormat)
	}
	if v, ok := value.(string); ok {
		return v
	}

	return ""
}

func unserializeTime(value interface{}) interface{} {
	switch value := value.(type) {
	case string:
		// try TIME_FORMAT
		t, err := time.Parse(domain.TimeFormat, value)
		if err == nil {
			return t
		}

		// try DATE_FORMAT
		t, err = time.Parse(domain.DateFormat, value)
		if err == nil {
			return t
		}

		return time.Time{}
	case []byte:
		return unserializeTime(string(value))
	case *string:
		if value == nil {
			return nil
		}
		return unserializeTime(*value)
	case time.Time:
		return value
	default:
		return time.Time{}
	}
}

var Time = graphql.NewScalar(
	graphql.ScalarConfig{
		Name:        "Time",
		Description: "A time fomatted to marketplace standards",
		Serialize:   serializeTime,
		ParseValue:  unserializeTime,
		ParseLiteral: func(valueAST ast.Value) interface{} {
			switch valueAST := valueAST.(type) {
			case *ast.StringValue:
				return unserializeTime(valueAST.Value)
			}

			return time.Time{}
		},
	},
)

// DATE -------------------------

func serializeDate(value interface{}) interface{} {
	if v, ok := value.(time.Time); ok {
		return v.Format(domain.DateFormat)
	}
	if v, ok := value.(*time.Time); ok {
		if v == nil {
			return ""
		}
		return v.Format(domain.DateFormat)
	}

	return ""
}

var Date = graphql.NewScalar(
	graphql.ScalarConfig{
		Name:        "Date",
		Description: "A date fomatted to marketplace standards",
		Serialize:   serializeDate,
		ParseValue:  unserializeTime,
		ParseLiteral: func(valueAST ast.Value) interface{} {
			switch valueAST := valueAST.(type) {
			case *ast.StringValue:
				return unserializeTime(valueAST.Value)
			}

			return time.Time{}
		},
	},
)

// BOOLEAN ------------------

var Boolean = graphql.NewScalar(
	graphql.ScalarConfig{
		Name:        "Bool",
		Description: "A boolean value",
		Serialize: func(value interface{}) interface{} {
			if v, ok := value.(bool); ok {
				return v
			}

			// when we want to convert time.Time -> bool, true if != time.Time{}
			if v, ok := value.(time.Time); ok {
				return !v.Equal(time.Time{})
			}
			return false
		},
		ParseValue: func(value interface{}) interface{} {
			if v, ok := value.(bool); ok {
				return v
			}
			if v, ok := value.(string); ok {
				val, err := strconv.ParseBool(v)
				if err == nil {
					return val
				}
			}
			return false
		},
		ParseLiteral: func(valueAST ast.Value) interface{} {
			switch valueAST := valueAST.(type) {
			case *ast.BooleanValue:
				return valueAST.Value
			}

			return false
		},
	},
)
