package cmd

import (
	"fmt"

	"futures-trader/config"

	"github.com/spf13/cobra"
)

var clearKeyCmd = &cobra.Command{
	Use:   "clear-key",
	Short: "清除保存的 API 密钥",
	Long:  `清除本地保存的 Gate.io API 密钥，清除后无法进行交易操作

示例:
  futures-trader clear-key`,
	RunE:  clearKey,
}

func init() {
	rootCmd.AddCommand(clearKeyCmd)
}

func clearKey(cmd *cobra.Command, args []string) error {
	err := config.ClearConfig()
	if err != nil {
		return fmt.Errorf("清除密钥失败: %v", err)
	}

	fmt.Println("API密钥已清除")
	return nil
}
