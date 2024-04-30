package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test",
	Long:  `Function Tests, Benchmark.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Test: Function Tests, Benchmarks.")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
