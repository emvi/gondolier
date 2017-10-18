package model

import (
	"reflect"
)

type metaModel struct {
	ModelName string
	Fields    []metaField
}

type metaField struct {
	Name string
	Tag  string
}

func buildMetaModel(model interface{}) metaModel {
	return metaModel{getModelName(model),
		getModelFields(model)}
}

func getModelName(model interface{}) string {
	if t := reflect.TypeOf(model); t.Kind() == reflect.Ptr {
		return t.Elem().Name()
	} else {
		return t.Name()
	}
}

func getModelFields(model interface{}) []metaField {
	val := reflect.ValueOf(model).Elem()
	fields := make([]metaField, 0)

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag
		fields = append(fields, metaField{field.Name, tag.Get("model")})
	}

	return fields
}
