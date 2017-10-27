package model

import (
	"reflect"
	"strings"
)

type metaModel struct {
	ModelName string
	Fields    []metaField
}

type metaField struct {
	Name string
	Tags []metaTag
}

type metaTag struct {
	Name  string
	Value string
}

func buildMetaModel(model interface{}) metaModel {
	return metaModel{getModelName(model),
		getModelFields(model)}
}

func getModelName(model interface{}) string {
	t := reflect.TypeOf(model)
	kind := t.Kind()

	if kind == reflect.Ptr {
		t = t.Elem()
		kind = t.Kind()
	}

	if kind != reflect.Struct {
		panic("Passed type is not a struct")
	}

	if t.Kind() == reflect.Ptr {
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
		tag := field.Tag.Get("model")
		kind := field.Type.Kind()

		if tag == "" || tag == "-" {
			continue
		}

		if kind == reflect.Struct || kind == reflect.Ptr || kind == reflect.Interface {
			panic("The type for field '" + field.Name + "' is invalid")
		}

		fields = append(fields, metaField{field.Name, parseTag(tag)})
	}

	return fields
}

func parseTag(tag string) []metaTag {
	tags := make([]metaTag, 0)
	elements := strings.Split(tag, ";")

	for _, e := range elements {
		nv := strings.Split(e, ":")

		if len(nv) == 1 {
			tags = append(tags, metaTag{"", nv[0]})
		} else if len(nv) == 2 {
			tags = append(tags, metaTag{nv[0], nv[1]})
		} else {
			panic("Too many or too few meta field tag separators")
		}
	}

	return tags
}
