package cmd

import (
	"fmt"

	"futures-trader/config"
	"futures-trader/trader"

	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "查询期货账户信息",
	Long:  `查询 Gate.io 期货账户信息，包括可用余额等

参数说明:
  - settle: 结算货币，支持 btc 或 usdt

示例:
  # 查询 USDT 账户信息
  futures-trader account --settle usdt

  # 查询 BTC 账户信息
  futures-trader account --settle btc`,
	RunE:  getAccount,
}

var accountSettle string

func init() {
	accountCmd.Flags().StringVar(&accountSettle, "settle", "usdt", "结算货币 (btc 或 usdt)")
	rootCmd.AddCommand(accountCmd)
}

func getAccount(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}
	if cfg == nil {
		return fmt.Errorf("未找到API密钥，请先使用 save-key 命令保存密钥")
	}

	if accountSettle != "btc" && accountSettle != "usdt" {
		return fmt.Errorf("settle必须是'btc'或'usdt'")
	}

	result, err := trader.GetFuturesAccountBalance(cfg.APIKey, cfg.APISecret, accountSettle)
	if err != nil {
		return fmt.Errorf("获取账户信息失败: %v", err)
	}

	printJSON(result)
	return nil
}
