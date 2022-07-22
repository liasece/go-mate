package repo

import (
	"reflect"
	"strings"

	"github.com/liasece/go-mate/src/utils"
	"github.com/liasece/gocoder"
	"github.com/liasece/gocoder/cde"
	"go.mongodb.org/mongo-driver/bson"
)

type FieldFilterField struct {
	f   gocoder.Field
	opt string
	gf  gocoder.Field
}

func newFieldFilterField(f gocoder.Field, opt string, gt gocoder.Type) *FieldFilterField {
	var fieldName string
	if opt == "NoCount" || opt == "Ands" || opt == "Ors" {
		fieldName = opt
	} else {
		fieldName = f.GetName() + opt
	}
	return &FieldFilterField{
		f:   f,
		opt: opt,
		gf:  cde.Field(fieldName, gt, ""),
	}
}

func fieldFilterFieldsToGocoder(mfs ...*FieldFilterField) []gocoder.Field {
	fs := make([]gocoder.Field, 0)
	for _, v := range mfs {
		fs = append(fs, v.gf)
	}
	return fs
}

func getFieldFilterFields(f gocoder.Field) []*FieldFilterField {
	bsonFiled := utils.GetFieldBSONName(f)
	if bsonFiled == "-" {
		return nil
	}
	fs := make([]*FieldFilterField, 0)
	ft := f.GetType()
	if ft.Kind() == reflect.Ptr {
		ft = ft.Elem()
	}
	filterFt := cde.Type(ft)
	if ft.Kind() != reflect.Interface {
		filterFt = filterFt.TackPtr()
	}
	if (ft.Kind() != reflect.Struct && ft.Kind() != reflect.Map) || ft.String() == "time.Time" {
		fs = append(fs, newFieldFilterField(f, "Eq", filterFt))
		fs = append(fs, newFieldFilterField(f, "Ne", filterFt))
		if ft.Kind() != reflect.Slice {
			fs = append(fs, newFieldFilterField(f, "In", cde.Type(ft).Slice()))
			fs = append(fs, newFieldFilterField(f, "Nin", cde.Type(ft).Slice()))
		} else {
			// slice
			fs = append(fs, newFieldFilterField(f, "Include", cde.Type(ft)))
			fs = append(fs, newFieldFilterField(f, "NotInclude", cde.Type(ft)))
			fs = append(fs, newFieldFilterField(f, "ElemEq", cde.Type(ft).Elem().TackPtr()))
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
		unPtrStr := ""
		if f.gf.GetType().IsPtr() {
			unPtrStr = "*"
		}
		var setter gocoder.Value
		switch f.opt {
		case "Eq":
			setter = resV.Index(bsonFiled).Set(rf)
		case "Ne":
			setter = resV.Index(bsonFiled).Set(cde.Value(`primitive.M{"$ne": `+unPtrStr+`f.`+f.gf.GetName()+` }`, nil))
		case "In":
			setter = resV.Index(bsonFiled).Set(cde.Value(`primitive.M{"$in": f.`+f.gf.GetName()+` }`, nil))
		case "Nin":
			setter = resV.Index(bsonFiled).Set(cde.Value(`primitive.M{"$nin": f.`+f.gf.GetName()+` }`, nil))
		case "Gt", "Gte", "Lt", "Lte":
			setter = resV.Index(bsonFiled).Set(cde.Value(`primitive.M{"$`+strings.ToLower(f.opt)+`": `+unPtrStr+`f.`+f.gf.GetName()+` }`, nil))
		case "Reg":
			setter = resV.Index(bsonFiled).Set(cde.Value(`primitive.Regex{Pattern: `+unPtrStr+`f.`+f.gf.GetName()+`, Options: "i"}`, nil))
		case "Include":
			setter = resV.Index(bsonFiled).Set(cde.Value(`primitive.M{"$elemMatch": primitive.M{"$in": f.`+f.gf.GetName()+`}}`, nil))
		case "NotInclude":
			setter = resV.Index(bsonFiled).Set(cde.Value(`primitive.M{"$elemMatch": primitive.M{"$nin": f.`+f.gf.GetName()+`}}`, nil))
		case "ElemEq":
			setter = resV.Index(bsonFiled).Set(cde.Value(`primitive.M{"$eq": f.`+f.gf.GetName()+` }`, nil))
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
	setC.C(
		cde.Value(`
			if len(f.Ands) > 0 {
				ands := []primitive.M{}
				for _, f := range f.Ands {
					b := f.ToBSON()
					if len(b) > 0 {
						ands = append(ands, b)
					}
				}
				if len(ands) > 0 {
					res["$and"] = ands
				}
			}`, nil),
	)
	setC.C(
		cde.Value(`
			if len(f.Ors) > 0 {
				ors := []primitive.M{}
				for _, f := range f.Ors {
					b := f.ToBSON()
					if len(b) > 0 {
						ors = append(ors, b)
					}
				}
				if len(ors) > 0 {
					res["$or"] = ors
				}
			}`, nil),
	)
	f.C(
		resV.AutoSet(cde.Make(bson.M{})),
		setC,
		cde.Return(resV),
	)
	c.C(f)
	return c
}
