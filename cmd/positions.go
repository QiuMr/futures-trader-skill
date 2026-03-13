package cmd

import (
	"fmt"

	"futures-trader/config"
	"futures-trader/trader"

	"github.com/spf13/cobra"
)

var positionsCmd = &cobra.Command{
	Use:   "positions",
	Short: "查询持仓信息",
	Long:  `查询当前持有的期货仓位信息（仅显示持仓不为0的仓位）

参数说明:
  - settle: 结算货币，支持 btc 或 usdt

示例:
  # 查询 USDT 账户的持仓
  futures-trader positions --settle usdt

  # 查询 BTC 账户的持仓
  futures-trader positions --settle btc`,
	RunE:  getPositions,
}

var positionsSettle string

func init() {
	positionsCmd.Flags().StringVar(&positionsSettle, "settle", "usdt", "结算货币 (btc 或 usdt)")
	rootCmd.AddCommand(positionsCmd)
}

func getPositions(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}
	if cfg == nil {
		return fmt.Errorf("未找到API密钥，请先使用 save-key 命令保存密钥")
	}

	if positionsSettle != "btc" && positionsSettle != "usdt" {
		return fmt.Errorf("settle必须是'btc'或'usdt'")
	}

	result, err := trader.GetRealPositions(cfg.APIKey, cfg.APISecret, positionsSettle)
	if err != nil {
		return fmt.Errorf("获取持仓信息失败: %v", err)
	}

	printJSON(result)
	return nil
}
