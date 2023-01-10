package utils

import (
	"reflect"
	"strings"
)

func SelectStmt(model interface{}, tableName string) string {
	stmt := "SELECT "
	t := reflect.TypeOf(model)
	for t.Kind() == reflect.Slice || t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	delim := ", "
	for i := 0; i < t.NumField(); i++ {
		stmt += t.Field(i).Tag.Get("db") + delim
	}
	return strings.TrimSuffix(stmt, delim) + " FROM `" + tableName + "`"
}
