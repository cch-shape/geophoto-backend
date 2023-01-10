package utils

import (
	"reflect"
	"strings"
)

func SelectStmt(model interface{}, tableName string, attachStmt ...*string) string {
	t := reflect.TypeOf(model)
	for t.Kind() == reflect.Slice || t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var stmt strings.Builder
	stmt.WriteString("SELECT ")
	delim := ""
	for i := 0; i < t.NumField(); i++ {
		stmt.WriteString(delim)
		stmt.WriteString(t.Field(i).Tag.Get("db"))
		delim = ", "
	}
	stmt.WriteString(" FROM `")
	stmt.WriteString(tableName)
	stmt.WriteByte('`')
	for i := 0; i < len(attachStmt); i++ {
		stmt.WriteString(*attachStmt[i])
	}
	return stmt.String()
}
