package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"futures-trader/config"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "futures-trader",
	Short: "Gate.io 期货交易 CLI 工具",
	Long:  `一个用于 Gate.io 期货交易的命令行工具，支持订单管理、持仓查询等功能

使用说明:
  1. 首先使用 save-key 命令保存您的 Gate.io API 密钥
  2. 然后可以使用其他命令进行交易操作
  3. 使用 clear-key 命令清除已保存的密钥

帮助:
  使用 [command] --help 查看具体命令的详细说明`,
}

var (
	cfg *config.Config
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.LoadConfig()
		if err != nil {
			return fmt.Errorf("加载配置失败: %v", err)
		}
		return nil
	}
}

func GetConfig() *config.Config {
	return cfg
}

func printJSON(data interface{}) {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON序列化错误: %v\n", err)
		return
	}
	fmt.Println(string(output))
}

func printError(msg string) {
	fmt.Fprintf(os.Stderr, "错误: %s\n", msg)
}
