package main

import (
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/liasece/go-mate/gogen/writer/repo"
	"github.com/liasece/gocoder"
	"github.com/liasece/gocoder/cde"
	"github.com/liasece/log"
	"github.com/spf13/cobra"
)

var entityFile string
var outputFile string
var entityNames []string

func buildRunner() {
	tmpPaths := make([]string, 0)
	_ = tmpPaths
	tmpFiles := make([]string, 0)
	_ = tmpFiles
	path, entityFileName := filepath.Split(entityFile)
	log.Info("in", log.Any("entityFile", entityFile), log.Any("path", path), log.Any("entityNames", entityNames))
	entityFileBaseName := strings.TrimSuffix(entityFileName, ".go")

	fset := token.NewFileSet()
	// 这里取绝对路径，方便打印出来的语法树可以转跳到编辑器
	f, err := parser.ParseFile(fset, entityFile, nil, parser.AllErrors)
	if err != nil {
		log.Error("parser.ParseFile error", log.ErrorField(err))
		return
	}
	_ = f
	if outputFile == "" {
		outputFile = filepath.Join(entityFileBaseName + "_struct_gen.go")
	}

	c := gocoder.NewCode()
	for _, entity := range entityNames {
		t, err := cde.LoadTypeFromSource(entityFile, entity)
		if err != nil {
			log.Error("LoadTypeFromSource error", log.ErrorField(err), log.Any("entityFile", entityFile), log.Any("entity", entity))
		}
		enGameEntry := repo.NewRepositoryWriterByType(t.RefType(), entity)
		c.C(enGameEntry.GetFilterTypeCode(), enGameEntry.GetUpdaterTypeCode())
	}

	gocoder.WriteToFile(outputFile, c, gocoder.NewToCodeOpt().PkgName(f.Name.Name))
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
	buildRunnerCmd.Flags().StringVarP(&outputFile, "out", "o", "", "The output file path")

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(buildRunnerCmd)
	rootCmd.Execute()
}
