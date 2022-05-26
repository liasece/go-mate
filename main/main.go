package main

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/liasece/go-mate/gogen/writer"
	"github.com/liasece/go-mate/gogen/writer/repo"
	"github.com/liasece/gocoder"
	"github.com/liasece/gocoder/cde"
	"github.com/liasece/log"
	"github.com/spf13/cobra"
)

func buildRunner(cfg *BuildCfg) {
	tmpPaths := make([]string, 0)
	_ = tmpPaths
	tmpFiles := make([]string, 0)
	_ = tmpFiles
	path, _ := filepath.Split(cfg.EntityFile)
	log.Info("in", log.Any("entityFile", cfg.EntityFile), log.Any("entityPkg", calGoFilePkgName(cfg.EntityFile)), log.Any("path", path), log.Any("entityNames", cfg.EntityNames))

	optCode := gocoder.NewCode()
	repositoryInterfaceCode := gocoder.NewCode()
	for _, entity := range cfg.EntityNames {
		t, err := cde.LoadTypeFromSource(cfg.EntityFile, entity)
		if err != nil {
			log.Error("LoadTypeFromSource error", log.ErrorField(err), log.Any("entityFile", cfg.EntityFile), log.Any("entity", entity))
			continue
		}
		if t == nil {
			log.Error("buildRunner LoadTypeFromSource not found", log.Any("entity", entity), log.Any("cfg", cfg))
			continue
		}
		t.SetNamed(entity)
		if cfg.EntityPkg == "" {
			cfg.EntityPkg = calGoFilePkgName(cfg.EntityFile)
		}
		{
			filedNames := make([]string, 0)
			for i := 0; i < t.NumField(); i++ {
				filedNames = append(filedNames, t.Field(i).GetName())
			}
			log.Info("buildRunner filedNames", log.Any("entity", entity), log.Any("filedNames", filedNames))
		}
		enGameEntry := repo.NewRepositoryWriterByType(t, entity, cfg.EntityPkg, cfg.OutputFilterSuffix, cfg.OutputUpdaterSuffix, cfg.OutputTypeSuffix)
		optCode.C(enGameEntry.GetFilterTypeCode(), enGameEntry.GetUpdaterTypeCode())
		repositoryInterfaceCode.C(enGameEntry.GetEntityRepositoryInterfaceCode())
		if cfg.OutputProtoFile != "" {
			writer.StructToProto(cfg.OutputProtoFile, t, cfg.GetOutputProtoIndent())
			filterStr, _ := enGameEntry.GetFilterTypeStructCode()
			writer.StructToProto(cfg.OutputProtoFile, filterStr.GetType(), cfg.GetOutputProtoIndent())
		}
		if cfg.OutputRepositoryAdapterFile != "" {
			c, err := enGameEntry.GetEntityRepositoryCodeFromTmpl(cfg.RepositoryTmplPath)
			if err != nil {
				log.Error("buildRunner GetEntityRepositoryCodeFromTmpl error", log.ErrorField(err))
			} else {
				gocoder.WriteToFile(cfg.OutputRepositoryAdapterFile, c, gocoder.NewToCodeOpt().PkgName(""))
			}
		}
	}
	if cfg.OutputFile != "" {
		if cfg.OutputPkg == "" {
			cfg.OutputPkg = calGoFilePkgName(cfg.OutputFile)
		}
		gocoder.WriteToFile(cfg.OutputFile, optCode, gocoder.NewToCodeOpt().PkgName(cfg.OutputPkg))
	}
	if cfg.OutputRepositoryInterfaceFile != "" {
		gocoder.WriteToFile(cfg.OutputRepositoryInterfaceFile, repositoryInterfaceCode, gocoder.NewToCodeOpt().PkgName(calGoFilePkgName(cfg.OutputRepositoryInterfaceFile)))
	}
}

type BuildCfg struct {
	EntityFile         string   `arg:"name: file; short: f; usage: the file path of target entity; required;"`
	EntityNames        []string `arg:"name: name; short: n; usage: the name list of target entity; required"`
	EntityPkg          string   `arg:"name: entity-pkg; usage: the entity package path of target entity"`
	RepositoryTmplPath string   `arg:"name: repository-tmpl-path; usage: the repository gen from tmpl"`

	// output
	OutputFile                    string `arg:"name: out; short: o; usage: the output file path"`
	OutputPkg                     string `arg:"name: pkg; short: p; usage: the output pkg name"`
	OutputRepositoryInterfaceFile string `arg:"name: out-rep-inf-file; usage: output repository interface file"`
	OutputRepositoryAdapterFile   string `arg:"name: out-rep-adp-file; usage: output repository adapter file"`
	OutputFilterSuffix            string `arg:"name: out-filter-suffix; usage: output filter type name suffix"`
	OutputUpdaterSuffix           string `arg:"name: out-updater-suffix; usage: output updater type name suffix"`
	OutputTypeSuffix              string `arg:"name: out-type-suffix; usage: output type name suffix"`
	OutputProtoFile               string `arg:"name: out-proto-file; usage: output proto file"`
	OutputProtoIndent             string `arg:"name: out-proto-indent; usage: output proto file indent($4,$tab)"`
}

func (c *BuildCfg) GetOutputProtoIndent() string {
	if c.OutputProtoIndent == "" {
		return "\t"
	}
	if strings.HasPrefix(c.OutputProtoIndent, "$") {
		switch c.OutputProtoIndent[1:] {
		case "tab":
			return "\t"
		default:
			a, _ := strconv.Atoi(c.OutputProtoIndent[1:])
			return strings.Repeat(" ", a)
		}
	}
	return c.OutputProtoIndent
}

func main() {
	cfg := &BuildCfg{}

	var buildRunnerCmd = &cobra.Command{
		Use:   "buildRunner",
		Short: "build a go main.go to target entity folder",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			buildRunner(cfg)
		},
	}
	initFlag(buildRunnerCmd, cfg)

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(buildRunnerCmd)
	rootCmd.Execute()
}
