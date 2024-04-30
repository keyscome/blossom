package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "blossom",
	Short: "CLI Tool for Blossom",
	Long:  `Blossom, a tool for the DevOps of Innovation, Like Kanban, Measure, Workflow, Artifact, Test, Configuration Management.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
