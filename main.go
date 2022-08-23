package main

import (
	"github.com/liasece/go-mate/cmd"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "app"}

	{
		// buildRunner
		cfg := &cmd.BuildCfg{}
		var buildRunnerCmd = &cobra.Command{
			Use:   "buildRunner",
			Short: "build a go main.go to target entity folder",
			Long:  "",
			Run: func(c *cobra.Command, args []string) {
				cfg.AfterLoad()
				cmd.BuildRunner(cfg)
			},
		}
		cmd.InitFlag(buildRunnerCmd, cfg)
		rootCmd.AddCommand(buildRunnerCmd)
	}
	{
		// generate
		cfg := &cmd.GenerateCfg{}
		var generateCmd = &cobra.Command{
			Use:   "generate",
			Short: "build a go main.go to target entity folder",
			Long:  "",
			Run: func(c *cobra.Command, args []string) {
				cmd.Generate(cfg)
			},
		}
		cmd.InitFlag(generateCmd, cfg)
		rootCmd.AddCommand(generateCmd)
	}
	rootCmd.Execute()
}
