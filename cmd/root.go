package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd 是 Blossom 工具的根命令
var rootCmd = &cobra.Command{
	Use:   "blossom",
	Short: "Blossom DevOps Tool",
	Long:  "Blossom 是一个 DevOps 工具，能够在执行命令的同时显示实时监控信息。",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("请使用子命令，例如： 'blossom monitor' 来启动监控界面")
	},
}

// Execute 由 main 调用
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("执行错误:", err)
		os.Exit(1)
	}
}

func init() {
	// 定义全局标志：配置文件路径
	rootCmd.PersistentFlags().StringP("config", "c", "config.yml", "配置文件路径")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}
