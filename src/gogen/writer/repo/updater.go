package repo

import (
	"reflect"

	"github.com/liasece/go-mate/src/utils"
	"github.com/liasece/gocoder"
	"github.com/liasece/gocoder/cde"
	"go.mongodb.org/mongo-driver/bson"
)

type FieldUpdaterField struct {
	f   gocoder.Field
	opt string
	gf  gocoder.Field
}

func newFieldUpdaterField(f gocoder.Field, opt string, gt gocoder.Type) *FieldUpdaterField {
	fieldName := ""
	if opt == "JustDelete" {
		fieldName = opt
	} else {
		fieldName = f.GetName() + opt
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

func getFieldUpdaterFields(f gocoder.Field, opts ...TypeOpt) []*FieldUpdaterField {
	cfg := &TypeCfg{}
	for _, opt := range opts {
		opt(cfg)
	}

	bsonFiled := utils.GetFieldBSONName(f)
	if bsonFiled == "-" {
		return nil
	}
	fs := make([]*FieldUpdaterField, 0)
	ft := f.GetType()
	if ft.Kind() == reflect.Ptr {
		ft = ft.Elem()
	}
	filterFt := cde.Type(ft)
	// if f.GetType().Kind() == reflect.Struct && ft.GetRowStr() != "" {
	// 	log.Error("reflect.Struct down to str", log.Any("name", ft.Name()), log.Any("pkg", ft.Package()), log.Any("ft", ft.String()), log.Any("rowStr", ft.GetRowStr()))
	// 	filterFt = cde.TypeD(ft.Package(), ft.Package()+"."+ft.Name())
	// }
	if ft.Kind() != reflect.Interface {
		filterFt = filterFt.TackPtr()
	}
	fs = append(fs, newFieldUpdaterField(f, "", filterFt))
	if ft.Kind() == reflect.Slice {
		if cfg.addSlicePBEmpty {
			fs = append(fs, newFieldUpdaterField(f, "PBEmpty", cde.Type(false)))
		}
		fs = append(fs, newFieldUpdaterField(f, "Add", cde.Type(ft)))
		fs = append(fs, newFieldUpdaterField(f, "Del", cde.Type(ft)))
		fs = append(fs, newFieldUpdaterField(f, "Replace", cde.Type(ft).Elem().TackPtr()))
	}
	if ft.Kind() >= reflect.Int && ft.Kind() <= reflect.Float64 {
		fs = append(fs, newFieldUpdaterField(f, "Inc", filterFt))
	}
	return fs
}

func getFieldUpdaterMethodToBSON(st gocoder.Struct, fs []*FieldUpdaterField) gocoder.Codable {
	c := gocoder.NewCode()
	receiver := cde.Receiver("f", st.GetType())
	f := cde.Method("ToBSON", receiver, nil, []gocoder.Type{cde.Type(bson.M{})})
	resV := cde.Value("res", bson.M{})
	setC := gocoder.NewCode()
	var addToSetV gocoder.Value
	var pullV gocoder.Value
	var setV gocoder.Value
	var incV gocoder.Value
	initSetV := func() {
		if setV == nil {
			setV = cde.Value("set", bson.M{})
			setC.C(
				setV.AutoSet(cde.Make(bson.M{})),
			)
		}
	}
	initIncV := func() {
		if incV == nil {
			incV = cde.Value("inc", bson.M{})
			setC.C(
				incV.AutoSet(cde.Make(bson.M{})),
			)
		}
	}
	for _, f := range fs {
		rf := receiver.GetValue().Dot(f.gf.GetName())
		bsonFiled := utils.GetFieldBSONName(f.f)
		if bsonFiled == "-" {
			continue
		}
		var setter gocoder.Value
		switch f.opt {
		case "":
			initSetV()
			setter = setV.Index(bsonFiled).Set(rf)
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
			initIncV()
			setter = incV.Index(bsonFiled).Set(rf)
		case "Replace":
			initSetV()
			setter = setV.Index(bsonFiled + ".$").Set(rf)
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
	if addToSetV != nil {
		setC.C(
			cde.If(cde.Len(addToSetV).GT(0)).C(
				resV.Index("$addToSet").Set(addToSetV),
			),
		)
	}
	if pullV != nil {
		setC.C(
			cde.If(cde.Len(pullV).GT(0)).C(
				resV.Index("$pull").Set(pullV),
			),
		)
	}
	if setV != nil {
		setC.C(
			cde.If(cde.Len(setV).GT(0)).C(
				resV.Index("$set").Set(setV),
			),
		)
	}
	if incV != nil {
		setC.C(
			cde.If(cde.Len(incV).GT(0)).C(
				resV.Index("$inc").Set(incV),
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
