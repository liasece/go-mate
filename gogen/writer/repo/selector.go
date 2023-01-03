package repo

import (
	"github.com/liasece/go-mate/utils"
	"github.com/liasece/gocoder"
	"github.com/liasece/gocoder/cde"
	"github.com/liasece/gocoder/cdt"
	"go.mongodb.org/mongo-driver/bson"
)

type FieldSelectorField struct {
	f   gocoder.Field
	opt string
	gf  gocoder.Field
}

func newFieldSelectorField(f gocoder.Field, opt string, gt gocoder.Type) *FieldSelectorField {
	fieldName := f.GetName() + opt
	return &FieldSelectorField{
		f:   f,
		opt: opt,
		gf:  cde.Field(fieldName, gt, ""),
	}
}

func fieldSelectorFieldsToGocoder(mfs []*FieldSelectorField) []gocoder.Field {
	fs := make([]gocoder.Field, 0)
	for _, v := range mfs {
		fs = append(fs, v.gf)
	}
	return fs
}

func getFieldSelectorFields(f gocoder.Field) []*FieldSelectorField {
	bsonFiled := utils.GetFieldBSONName(f)
	if bsonFiled == "-" {
		return nil
	}
	fs := make([]*FieldSelectorField, 0)
	fs = append(fs, newFieldSelectorField(f, "", cdt.Bool().TackPtr()))
	return fs
}

func getFieldSelectorMethodToBSON(st gocoder.Struct, fs []*FieldSelectorField) gocoder.Codable {
	c := gocoder.NewCode()
	receiver := cde.Receiver("f", st.GetType())
	f := cde.Method("ToBSON", receiver, nil, []gocoder.Arg{cde.Arg("", cde.Type(bson.M{}))})
	resV := cde.Value("res", bson.M{})
	selectorC := gocoder.NewCode()
	for _, f := range fs {
		rf := receiver.GetValue().Dot(f.gf.GetName())
		bsonFiled := utils.GetFieldBSONName(f.f)
		if bsonFiled == "-" {
			continue
		}
		var selector gocoder.Value
		switch f.opt {
		case "":
			selector = resV.Index(bsonFiled).Set(rf)
		default:
			// log.Info("unknown opt", log.Any("opt", f.opt))
		}
		if selector == nil {
			continue
		}
		selectorC.C(
			cde.PtrCheckerNotNil(rf).C(
				selector,
			),
		)
	}
	f.C(
		resV.AutoSet(cde.Make(bson.M{})),
		selectorC,
		cde.Return(resV),
	)
	c.C(f)
	return c
}
