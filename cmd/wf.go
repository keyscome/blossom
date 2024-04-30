package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var wfCmd = &cobra.Command{
	Use:   "wf",
	Short: "Workflow Engine",
	Long:  `Use Github Workflow, Jenkins, Tekton, Argo Workflow.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Workflow Engine: Github Workflow, Jenkins, Tekton, Argo Workflow.")
	},
}

func init() {
	rootCmd.AddCommand(wfCmd)
}
