package utils

import (
	"reflect"
	"strings"

	"github.com/liasece/gocoder"
)

func GetFieldBSONName(f gocoder.Field) string {
	if f == nil {
		return ""
	}
	bsonFiled := f.GetName()
	if v := reflect.StructTag(f.GetTag()).Get("bson"); v != "" {
		bsonFiled = v
		vs := strings.Split(v, ",")
		if len(vs) > 1 {
			bsonFiled = vs[0]
		}
	}
	return bsonFiled
}
