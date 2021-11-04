package main

import (
	"path/filepath"

	"github.com/liasece/go-mate/gogen/writer/repo"
	"github.com/liasece/gocoder"
	"github.com/liasece/gocoder/cde"
	"github.com/liasece/log"
	"github.com/spf13/cobra"
)

var entityFile string
var entityNames []string

// en := repo.NewRepositoryWriterByObj(entity.GameDetail{})
// filterStrT, _ := en.GetFilterTypeStructCode()
// updaterStrT, _ := en.GetFilterTypeStructCode()
// c := gocoder.NewCode()
// c.C(
// 	en.GetEntityRepositoryCode(filterStrT, updaterStrT),
// 	en.GetFilterTypeCode(),
// 	en.GetUpdaterTypeCode(),
// )
// err := gocoder.WriteToFile("gen/test_gen.go", c)
// if err != nil {
// 	log.Fatal("WriteToFileStr error", log.ErrorField(err))
// }

func buildRunner() {
	tmpPaths := make([]string, 0)
	tmpFiles := make([]string, 0)
	_ = tmpFiles
	path, _ := filepath.Split(entityFile)
	log.Info("in", log.Any("entityFile", entityFile), log.Any("path", path), log.Any("entityNames", entityNames))
	c := gocoder.NewCode()

	c.C(
		cde.Func("main", nil, nil).C(
			cde.Value("", repo.NewRepositoryWriterByObj).Call(cde.Value("entity.GameDetail{}", nil)),
			cde.Return(),
		),
	)

	mainPath := filepath.Join(path, "goMate")
	tmpPaths = append(tmpPaths, mainPath)
	gocoder.WriteToFile(filepath.Join(mainPath, "main.go"), c, gocoder.NewToCodeOpt().PkgName("main"))
	// for _, path := range tmpPaths {
	// 	os.RemoveAll(path)
	// }
	// for _, path := range tmpFiles {
	// 	os.Remove(path)
	// }
}

func main() {
	var buildRunnerCmd = &cobra.Command{
		Use:   "buildRunner",
		Short: "build a go main.go to target entity folder",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			buildRunner()
		},
	}
	buildRunnerCmd.Flags().StringVarP(&entityFile, "file", "f", "", "The file path of target entity")
	buildRunnerCmd.MarkFlagRequired("file")
	buildRunnerCmd.Flags().StringArrayVarP(&entityNames, "name", "n", nil, "The name list of target entity")
	buildRunnerCmd.MarkFlagRequired("name")

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(buildRunnerCmd)
	rootCmd.Execute()
}
