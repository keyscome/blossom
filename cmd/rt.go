package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rtCmd = &cobra.Command{
	Use:   "rt",
	Short: "Artifact",
	Long:  `Use Oringal Platform for Every Type of Artifact, Use Nexus for Private.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Artifact, Use Oringal Platform for Every Type of Artifact, Use Nexus for Private.")
	},
}

func init() {
	rootCmd.AddCommand(rtCmd)
}
