package validate

import (
	"github.com/go-playground/validator/v10"
	"reflect"
)

var v = validator.New()

func init() {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("json")
	})
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func Struct(model interface{}) []*map[string]interface{} {
	var errors []*map[string]interface{}

	err := v.Struct(model)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			el := map[string]interface{}{
				"failed_field": err.Field(),
				"tag":          err.Tag(),
			}
			if len(err.Param()) != 0 {
				el["value"] = err.Param()
			}
			errors = append(errors, &el)
		}
		return errors
	}
	return nil
}
