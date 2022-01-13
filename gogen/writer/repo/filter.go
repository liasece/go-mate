package repo

import (
	"reflect"
	"strings"

	"github.com/liasece/go-mate/utils"
	"github.com/liasece/gocoder"
	"github.com/liasece/gocoder/cde"
	"go.mongodb.org/mongo-driver/bson"
)

type FieldFilterField struct {
	f   reflect.StructField
	opt string
	gf  gocoder.Field
}

func newFieldFilterField(f reflect.StructField, opt string, gt gocoder.Type) *FieldFilterField {
	fieldName := f.Name + opt
	if opt == "NoCount" {
		fieldName = opt
	}
	return &FieldFilterField{
		f:   f,
		opt: opt,
		gf:  cde.Field(fieldName, gt, ""),
	}
}

func fieldFilterFieldsToGocoder(mfs []*FieldFilterField) []gocoder.Field {
	fs := make([]gocoder.Field, 0)
	for _, v := range mfs {
		fs = append(fs, v.gf)
	}
	return fs
}

func getFieldFilterFields(f reflect.StructField) []*FieldFilterField {
	bsonFiled := utils.GetFieldBSONName(f)
	if bsonFiled == "-" {
		return nil
	}
	fs := make([]*FieldFilterField, 0)
	ft := f.Type
	if ft.Kind() == reflect.Ptr {
		ft = ft.Elem()
	}
	filterFt := cde.Type(ft).TackPtr()
	if (ft.Kind() != reflect.Struct && ft.Kind() != reflect.Map) || ft.String() == "time.Time" {
		fs = append(fs, newFieldFilterField(f, "Eq", filterFt))
		fs = append(fs, newFieldFilterField(f, "Ne", filterFt))
		if ft.Kind() != reflect.Slice {
			fs = append(fs, newFieldFilterField(f, "In", cde.Type(ft).Slice()))
			fs = append(fs, newFieldFilterField(f, "Nin", cde.Type(ft).Slice()))
		}
	}
	if (ft.Kind() >= reflect.Int && ft.Kind() <= reflect.Float64) || ft.String() == "time.Time" {
		fs = append(fs, newFieldFilterField(f, "Gt", filterFt))
		fs = append(fs, newFieldFilterField(f, "Gte", filterFt))
		fs = append(fs, newFieldFilterField(f, "Lt", filterFt))
		fs = append(fs, newFieldFilterField(f, "Lte", filterFt))
	} else if ft.Kind() == reflect.String {
		fs = append(fs, newFieldFilterField(f, "Reg", filterFt)) // 正则
	}
	return fs
}

func getFieldFilterMethodToBSON(st gocoder.Struct, fs []*FieldFilterField) gocoder.Codeable {
	c := gocoder.NewCode()
	receiver := cde.Receiver("f", st.GetType())
	f := cde.Method("ToBSON", receiver, nil, []gocoder.Type{cde.Type(bson.M{})})
	resV := cde.Value("res", bson.M{})
	setC := gocoder.NewCode()
	for _, f := range fs {
		rf := receiver.GetValue().Dot(f.gf.GetName())
		bsonFiled := utils.GetFieldBSONName(f.f)
		if bsonFiled == "-" {
			continue
		}
		var setter gocoder.Value
		switch f.opt {
		case "Eq":
			setter = resV.Index(bsonFiled).Set(rf)
		case "Ne":
			setter = resV.Index(bsonFiled).Set(cde.Value(`primitive.M{"$ne": *f.`+f.gf.GetName()+` }`, nil))
		case "In":
			setter = resV.Index(bsonFiled).Set(cde.Value(`primitive.M{"$in": f.`+f.gf.GetName()+` }`, nil))
		case "Nin":
			setter = resV.Index(bsonFiled).Set(cde.Value(`primitive.M{"$nin": f.`+f.gf.GetName()+` }`, nil))
		case "Gt", "Gte", "Lt", "Lte":
			setter = resV.Index(bsonFiled).Set(cde.Value(`primitive.M{"$`+strings.ToLower(f.opt)+`": *f.`+f.gf.GetName()+` }`, nil))
		case "Reg":
			setter = resV.Index(bsonFiled).Set(cde.Value(`primitive.Regex{Pattern: *f.`+f.gf.GetName()+`, Options: "i"}`, nil))
		default:
			// log.Info("unknown opt", log.Any("opt", f.opt))
		}
		if setter == nil {
			continue
		}
		setC.C(
			cde.PtrCheckerNotNil(rf).C(
				setter,
			),
		)
	}
	f.C(
		resV.AutoSet(cde.Make(bson.M{})),
		setC,
		cde.Return(resV),
	)
	c.C(f)
	return c
}
