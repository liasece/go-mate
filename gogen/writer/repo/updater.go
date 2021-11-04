package repo

import (
	"reflect"

	"github.com/liasece/go-mate/utils"

	"github.com/liasece/gocoder"
	"github.com/liasece/gocoder/cde"
	"github.com/liasece/log"
	"go.mongodb.org/mongo-driver/bson"
)

type FieldUpdaterField struct {
	f   reflect.StructField
	opt string
	gf  gocoder.Field
}

func newFieldUpdaterField(f reflect.StructField, opt string, gt gocoder.Type) *FieldUpdaterField {
	fieldName := f.Name + opt
	if opt == "JustDelete" {
		fieldName = opt
	}
	return &FieldUpdaterField{
		f:   f,
		opt: opt,
		gf:  cde.Field(fieldName, gt, ""),
	}
}

func fieldUpdaterFieldsToGocoder(mfs []*FieldUpdaterField) []gocoder.Field {
	fs := make([]gocoder.Field, 0)
	for _, v := range mfs {
		fs = append(fs, v.gf)
	}
	return fs
}

func getFieldUpdaterFields(f reflect.StructField) []*FieldUpdaterField {
	bsonFiled := utils.GetFieldBSONName(f)
	if bsonFiled == "-" {
		return nil
	}
	fs := make([]*FieldUpdaterField, 0)
	ft := f.Type
	if ft.Kind() == reflect.Ptr {
		ft = ft.Elem()
	}
	filterFt := cde.Type(ft).TackPtr()
	fs = append(fs, newFieldUpdaterField(f, "", filterFt))
	if ft.Kind() == reflect.Slice {
		fs = append(fs, newFieldUpdaterField(f, "Add", cde.Type(ft)))
		fs = append(fs, newFieldUpdaterField(f, "Del", cde.Type(ft)))
		fs = append(fs, newFieldUpdaterField(f, "Replace", cde.Type(ft).Elem().TackPtr()))
	}
	if ft.Kind() >= reflect.Int && ft.Kind() <= reflect.Float64 {
		fs = append(fs, newFieldUpdaterField(f, "Inc", filterFt))
	}
	return fs
}

func getFieldUpdaterMethodToBSON(st gocoder.Struct, fs []*FieldUpdaterField) gocoder.Codeable {
	c := gocoder.NewCode()
	receiver := cde.Receiver("f", st.GetType())
	f := cde.Method("ToBSON", receiver, nil, []gocoder.Type{cde.Type(bson.M{})})
	resV := cde.Value("res", bson.M{})
	setC := gocoder.NewCode()
	var addToSetV gocoder.Value
	var pullV gocoder.Value
	for _, f := range fs {
		rf := receiver.GetValue().Dot(f.gf.GetName())
		bsonFiled := utils.GetFieldBSONName(f.f)
		if bsonFiled == "-" {
			continue
		}
		var setter gocoder.Value
		switch f.opt {
		case "":
			setter = resV.Index(bsonFiled).Set(rf)
		case "Add":
			if addToSetV == nil {
				addToSetV = cde.Value("addToSet", bson.M{})
				setC.C(
					addToSetV.AutoSet(cde.Make(bson.M{})),
				)
			}
			setter = addToSetV.Index(bsonFiled).Set(cde.Value(`primitive.M{"$each": f.`+f.gf.GetName()+`}`, nil))
		case "Del":
			if pullV == nil {
				pullV = cde.Value("pull", bson.M{})
				setC.C(
					pullV.AutoSet(cde.Make(bson.M{})),
				)
			}
			setter = pullV.Index(bsonFiled).Set(cde.Value(`primitive.M{"$in": f.`+f.gf.GetName()+`}`, nil))
		case "Inc":
			setter = resV.Index(bsonFiled).Set(cde.Value(`primitive.M{"$inc": *f.`+f.gf.GetName()+` }`, nil))
		case "Replace":
			setter = resV.Index(bsonFiled + ".$").Set(rf)
		default:
			log.Info("unknown opt", log.Any("opt", f.opt))
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
	if addToSetV != nil {
		setC.C(
			resV.Index("$addToSet").Set(addToSetV),
		)
	}
	if pullV != nil {
		setC.C(
			resV.Index("$pull").Set(pullV),
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
