package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var measureCmd = &cobra.Command{
	Use:   "measure",
	Short: "Measure",
	Long:  `Measure Time Duration, Performance.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Measure Time Duration, Performance.")
	},
}

func init() {
	rootCmd.AddCommand(measureCmd)
}
