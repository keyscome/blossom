package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmCmd = &cobra.Command{
	Use:   "cm",
	Short: "Configuration Management",
	Long:  `Use Ansible, Terraform, Helm.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Configuration Management: Ansible, Terraform, Helm.")
	},
}

func init() {
	rootCmd.AddCommand(cmCmd)
}
