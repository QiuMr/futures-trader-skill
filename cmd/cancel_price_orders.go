package cmd

import (
	"fmt"

	"futures-trader/config"
	"futures-trader/trader"

	"github.com/spf13/cobra"
)

var cancelPriceOrdersCmd = &cobra.Command{
	Use:   "cancel-price-orders",
	Short: "批量取消自动订单",
	Long:  `批量取消已创建的价格触发订单

参数说明:
  - settle: 结算货币，支持 btc 或 usdt
  - contract: 合约标识，可选，指定合约则只取消该合约的订单

示例:
  # 取消所有自动订单
  futures-trader cancel-price-orders

  # 只取消 BTC_USDT 合约的自动订单
  futures-trader cancel-price-orders --contract BTC_USDT`,
	RunE:  cancelPriceOrders,
}

var (
	cancelOrdersSettle   string
	cancelOrdersContract string
)

func init() {
	cancelPriceOrdersCmd.Flags().StringVar(&cancelOrdersSettle, "settle", "usdt", "结算货币 (btc 或 usdt)")
	cancelPriceOrdersCmd.Flags().StringVar(&cancelOrdersContract, "contract", "", "合约标识")
	rootCmd.AddCommand(cancelPriceOrdersCmd)
}

func cancelPriceOrders(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}
	if cfg == nil {
		return fmt.Errorf("未找到API密钥，请先使用 save-key 命令保存密钥")
	}

	if cancelOrdersSettle != "btc" && cancelOrdersSettle != "usdt" {
		return fmt.Errorf("settle必须是'btc'或'usdt'")
	}

	if cancelOrdersContract != "" && !trader.IsValidContract(cancelOrdersContract) {
		return fmt.Errorf("contract格式应为'基础货币_结算货币'，例如: BTC_USDT")
	}

	result, err := trader.CancelAllPriceOrders(cfg.APIKey, cfg.APISecret, cancelOrdersSettle, cancelOrdersContract)
	if err != nil {
		return fmt.Errorf("取消自动订单失败: %v", err)
	}

	printJSON(result)
	return nil
}
