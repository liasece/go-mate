package main

import (
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "app"}

	{
		// buildRunner
		cfg := &BuildCfg{}
		var buildRunnerCmd = &cobra.Command{
			Use:   "buildRunner",
			Short: "build a go main.go to target entity folder",
			Long:  "",
			Run: func(cmd *cobra.Command, args []string) {
				cfg.AfterLoad()
				buildRunner(cfg)
			},
		}
		initFlag(buildRunnerCmd, cfg)
		rootCmd.AddCommand(buildRunnerCmd)
	}
	{
		// generate
		cfg := &GenerateCfg{}
		var generateCmd = &cobra.Command{
			Use:   "generate",
			Short: "build a go main.go to target entity folder",
			Long:  "",
			Run: func(cmd *cobra.Command, args []string) {
				generate(cfg)
			},
		}
		initFlag(generateCmd, cfg)
		rootCmd.AddCommand(generateCmd)
	}
	rootCmd.Execute()
}
