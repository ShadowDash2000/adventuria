package pbhelper

import (
	"fmt"
	"strings"

	"github.com/pocketbase/dbx"
)

func DotExpand(parts ...string) string {
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for _, part := range parts[1:] {
		result += "." + part
	}
	return result
}

func Eq(a, b string) string {
	return a + " = " + b
}

func GreaterThan(a, b string) string {
	return a + " > " + b
}

func SliceToEqExp[V any](field string, values []V) []dbx.Expression {
	exp := make([]dbx.Expression, len(values))
	for i, value := range values {
		exp[i] = dbx.HashExp{field: value}
	}
	return exp
}

func SliceToAny[V any](values []V) []any {
	result := make([]any, len(values))
	for i, value := range values {
		result[i] = value
	}
	return result
}

func SliceToSqlString[V any](slice []V) string {
	quotedValues := make([]string, len(slice))
	for i, v := range slice {
		quotedValues[i] = fmt.Sprintf("'%v'", v)
	}
	return strings.Join(quotedValues, ", ")
}
