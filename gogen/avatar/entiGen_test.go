package avatar

import (
	"testing"

	"solarland/backendv2/tools/gogen/avatar/entity"
	"solarland/backendv2/tools/gogen/writer/repo"

	"github.com/liasece/gocoder"
	"github.com/liasece/log"
)

func TestFoo(t *testing.T) {
	en := repo.NewRepositoryWriterByObj(entity.GameDetail{})
	c := gocoder.NewCode()
	c.C(
		en.GetFilterTypeCode(),
		en.GetUpdaterTypeCode(),
	)
	err := gocoder.WriteToFile("gen/test_gen.go", c)
	if err != nil {
		log.Fatal("WriteToFileStr error", log.ErrorField(err))
	}
}

func TestFoo1(t *testing.T) {
	en := repo.NewRepositoryWriterByObj(entity.GameDetail{})
	filterStrT, _ := en.GetFilterTypeStructCode()
	updaterStrT, _ := en.GetFilterTypeStructCode()
	c := gocoder.NewCode()
	c.C(
		en.GetEntityRepositoryCode(filterStrT, updaterStrT),
		en.GetFilterTypeCode(),
		en.GetUpdaterTypeCode(),
	)
	err := gocoder.WriteToFile("gen/test_gen.go", c)
	if err != nil {
		log.Fatal("WriteToFileStr error", log.ErrorField(err))
	}
}
