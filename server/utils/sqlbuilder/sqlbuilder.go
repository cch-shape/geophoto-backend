package sqlbuilder

import (
	"golang.org/x/exp/slices"
	"os"
	"reflect"
	"regexp"
	"strings"
)

func Select(model interface{}, tableName string, attachStmt ...string) string {
	t := reflect.TypeOf(model)
	for t.Kind() == reflect.Slice || t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var stmt strings.Builder
	stmt.WriteString("SELECT ")
	delim := ""
	for i := 0; i < t.NumField(); i++ {
		if colName, exist := t.Field(i).Tag.Lookup("db"); exist {
			stmt.WriteString(delim)
			stmt.WriteString(colName)
			delim = ", "
		} else if colName, exist := t.Field(i).Tag.Lookup("db_cal"); exist {
			stmt.WriteString(delim)
			stmt.WriteString(stringResolveEnv(colName))
			stmt.WriteByte(' ')
			stmt.WriteString(strings.ToLower(t.Field(i).Name))
			delim = ", "
		}
	}
	stmt.WriteString(" FROM `")
	stmt.WriteString(tableName)
	stmt.WriteString("` ")
	for i := 0; i < len(attachStmt); i++ {
		stmt.WriteString(attachStmt[i])
	}
	return stmt.String()
}

func Insert(model interface{}, tableName string, attachStmt ...string) string {
	t := reflect.TypeOf(model)
	for t.Kind() == reflect.Slice || t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var colBuilder strings.Builder
	var valBuilder strings.Builder
	delim := ""
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag
		if colName, exist := tags.Lookup("db"); exist {
			if !slices.Contains([]string{"auto", "key"}, tags.Get("db_prop")) {
				colBuilder.WriteString(delim)
				colBuilder.WriteString(colName)
				valBuilder.WriteString(delim)
				valBuilder.WriteByte(':')
				valBuilder.WriteString(t.Field(i).Tag.Get("json"))
				delim = ", "
			}
		}
	}
	col := colBuilder.String()
	val := valBuilder.String()

	var rtBuilder strings.Builder
	delim = ""
	for i := 0; i < t.NumField(); i++ {
		if colName, exist := t.Field(i).Tag.Lookup("db"); exist {
			rtBuilder.WriteString(delim)
			rtBuilder.WriteString(colName)
			delim = ", "
		} else if colName, exist := t.Field(i).Tag.Lookup("db_cal"); exist {
			rtBuilder.WriteString(delim)
			rtBuilder.WriteString(stringResolveEnv(colName))
			rtBuilder.WriteByte(' ')
			rtBuilder.WriteString(strings.ToLower(t.Field(i).Name))
			delim = ", "
		}
	}
	rt := rtBuilder.String()

	var stmt strings.Builder
	stmt.WriteString("INSERT INTO `")
	stmt.WriteString(tableName)
	stmt.WriteString("` (")
	stmt.WriteString(col)
	if len(attachStmt) > 0 {
		stmt.WriteString(delim)
		stmt.WriteString(attachStmt[0])
	}
	stmt.WriteString(") VALUES (")
	stmt.WriteString(val)
	if len(attachStmt) > 1 {
		stmt.WriteString(delim)
		stmt.WriteString(attachStmt[1])
	}
	stmt.WriteString(") RETURNING ")
	stmt.WriteString(rt)
	if len(attachStmt) > 2 {
		stmt.WriteString(delim)
		stmt.WriteString(attachStmt[2])
	}

	return stmt.String()
}

func Update(model interface{}, tableName string, attachStmt ...string) string {
	t := reflect.TypeOf(model)
	for t.Kind() == reflect.Slice || t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// UPDATE t SET A=a, B=b WHERE k=k
	key_db := "id"
	key_json := "id"
	var colBuilder strings.Builder
	delim := ""
	for i := 0; i < t.NumField(); i++ {
		tags := t.Field(i).Tag
		if colName, exist := tags.Lookup("db"); exist {
			if !slices.Contains([]string{"auto", "key"}, tags.Get("db_prop")) {
				colBuilder.WriteString(delim)
				colBuilder.WriteString(colName)
				colBuilder.WriteString("=:")
				colBuilder.WriteString(t.Field(i).Tag.Get("json"))
				delim = ", "
			} else if tags.Get("db_prop") == "key" {
				key_db = colName
				key_json = t.Field(i).Tag.Get("json")
			}
		}
	}
	col := colBuilder.String()

	var stmt strings.Builder
	stmt.WriteString("UPDATE `")
	stmt.WriteString(tableName)
	stmt.WriteString("` SET ")
	stmt.WriteString(col)
	if len(attachStmt) > 0 {
		stmt.WriteString(delim)
		stmt.WriteString(attachStmt[0])
	}
	stmt.WriteString(" WHERE ")
	stmt.WriteString(key_db)
	stmt.WriteString("=:")
	stmt.WriteString(key_json)
	if len(attachStmt) > 1 {
		stmt.WriteString(delim)
		stmt.WriteString(attachStmt[1])
	}

	return stmt.String()
}

func stringResolveEnv(str string) string {
	if re, err := regexp.Compile(`\$\{(.*?)\}`); err == nil {
		if matched := re.FindAllString(str, -1); len(matched) != 0 {
			for _, el := range matched {
				str = strings.Replace(str, el, os.Getenv(el[2:len(el)-1]), 1)
			}
		}
	}
	return str
}
