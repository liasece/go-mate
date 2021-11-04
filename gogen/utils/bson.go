package utils

import (
	"reflect"
	"strings"
)

func GetFieldBSONName(f reflect.StructField) string {
	bsonFiled := f.Name
	if v := f.Tag.Get("bson"); v != "" {
		bsonFiled = v
		vs := strings.Split(v, ",")
		if len(vs) > 1 {
			bsonFiled = vs[0]
		}
	}
	return bsonFiled
}
