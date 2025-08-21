package tools

import (
	"reflect"
	"strings"
)

func GetFieldToColumnMap(model any) map[string]string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	mapping := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		gormTag := field.Tag.Get("gorm")

		for _, part := range strings.Split(gormTag, ";") {
			if strings.HasPrefix(part, "column:") {
				columnName := part[len("column:"):]
				// Convierte la key a minúscula
				mapping[strings.ToLower(field.Name)] = columnName
			}
		}
	}

	return mapping
}
