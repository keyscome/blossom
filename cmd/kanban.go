package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var kanbanCmd = &cobra.Command{
	Use:   "kanban",
	Short: "Kanban",
	Long:  `Visualizing the Cooperation Using a Kanban Board.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Visualizing cooperation using a Kanban board.")
	},
}

func init() {
	rootCmd.AddCommand(kanbanCmd)
}
