package repo

import (
	"github.com/liasece/go-mate/src/utils"
	"github.com/liasece/gocoder"
	"github.com/liasece/gocoder/cde"
	"github.com/liasece/gocoder/cdt"
	"go.mongodb.org/mongo-driver/bson"
)

type FieldSorterField struct {
	f   gocoder.Field
	opt string
	gf  gocoder.Field
}

func newFieldSorterField(f gocoder.Field, opt string, gt gocoder.Type) *FieldSorterField {
	fieldName := f.GetName() + opt
	return &FieldSorterField{
		f:   f,
		opt: opt,
		gf:  cde.Field(fieldName, gt, ""),
	}
}

func fieldSorterFieldsToGocoder(mfs []*FieldSorterField) []gocoder.Field {
	fs := make([]gocoder.Field, 0)
	for _, v := range mfs {
		fs = append(fs, v.gf)
	}
	return fs
}

func getFieldSorterFields(f gocoder.Field) []*FieldSorterField {
	bsonFiled := utils.GetFieldBSONName(f)
	if bsonFiled == "-" {
		return nil
	}
	fs := make([]*FieldSorterField, 0)
	fs = append(fs, newFieldSorterField(f, "", cdt.Int().TackPtr()))
	return fs
}

func getFieldSorterMethodToBSON(st gocoder.Struct, fs []*FieldSorterField) gocoder.Codable {
	c := gocoder.NewCode()
	receiver := cde.Receiver("f", st.GetType())
	f := cde.Method("ToBSON", receiver, nil, []gocoder.Type{cde.Type(bson.D{})})
	setC := gocoder.NewCode()
	var sortV gocoder.Value
	initSortV := func() {
		if sortV == nil {
			sortV = cde.Value("sort", bson.D{})
			setC.C(
				sortV.AutoSet(cde.Make(bson.D{}, 0)),
			)
		}
	}
	for _, f := range fs {
		rf := receiver.GetValue().Dot(f.gf.GetName())
		bsonFiled := utils.GetFieldBSONName(f.f)
		if bsonFiled == "-" {
			continue
		}
		var sorter gocoder.Value
		switch f.opt {
		case "":
			initSortV()
			// sorter = setV.Index(bsonFiled).Set(rf)
			sorter = cde.Value(`sort = append(sort, primitive.E{Key: "`+bsonFiled+`", Value: *f.`+f.gf.GetName()+`})`, nil)
		default:
			// log.Info("unknown opt", log.Any("opt", f.opt))
		}
		if sorter == nil {
			continue
		}
		setC.C(
			cde.PtrCheckerNotNil(rf).C(
				sorter,
			),
		)
	}
	f.C(
		setC,
		cde.Return(sortV),
	)
	c.C(f)
	return c
}
