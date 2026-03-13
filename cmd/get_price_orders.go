package cmd

import (
	"fmt"

	"futures-trader/config"
	"futures-trader/trader"

	"github.com/spf13/cobra"
)

var getPriceOrdersCmd = &cobra.Command{
	Use:   "get-price-orders",
	Short: "查询自动订单列表",
	Long:  `查询已创建的价格触发订单（止损/止盈单）

参数说明:
  - settle: 结算货币，支持 btc 或 usdt
  - status: 订单状态，支持:
      open: 未成交
      closed: 已成交
      cancelled: 已取消
  - contract: 合约标识，可选，指定合约则只查询该合约的订单
  - limit: 返回数量限制，可选，默认100
  - offset: 偏移量，可选

示例:
  # 查询所有未成交的自动订单
  futures-trader get-price-orders --status open

  # 查询 BTC_USDT 合约的已成交订单
  futures-trader get-price-orders --status closed --contract BTC_USDT

  # 分页查询，每页10条
  futures-trader get-price-orders --status open --limit 10 --offset 0`,
	RunE:  getPriceOrders,
}

var (
	priceOrdersSettle   string
	priceOrdersStatus   string
	priceOrdersContract string
	priceOrdersLimit    int
	priceOrdersOffset   int
)

func init() {
	getPriceOrdersCmd.Flags().StringVar(&priceOrdersSettle, "settle", "usdt", "结算货币 (btc 或 usdt)")
	getPriceOrdersCmd.Flags().StringVar(&priceOrdersStatus, "status", "open", "订单状态 (open/closed/cancelled)")
	getPriceOrdersCmd.Flags().StringVar(&priceOrdersContract, "contract", "", "合约标识")
	getPriceOrdersCmd.Flags().IntVar(&priceOrdersLimit, "limit", 100, "返回数量限制")
	getPriceOrdersCmd.Flags().IntVar(&priceOrdersOffset, "offset", 0, "偏移量")
	rootCmd.AddCommand(getPriceOrdersCmd)
}

func getPriceOrders(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}
	if cfg == nil {
		return fmt.Errorf("未找到API密钥，请先使用 save-key 命令保存密钥")
	}

	if priceOrdersSettle != "btc" && priceOrdersSettle != "usdt" {
		return fmt.Errorf("settle必须是'btc'或'usdt'")
	}

	validStatus := []string{"open", "closed", "cancelled"}
	if !trader.IsValidStatus(priceOrdersStatus, validStatus) {
		return fmt.Errorf("status必须是: %v", validStatus)
	}

	if priceOrdersContract != "" && !trader.IsValidContract(priceOrdersContract) {
		return fmt.Errorf("contract格式应为'基础货币_结算货币'，例如: BTC_USDT")
	}

	if priceOrdersLimit < 0 {
		return fmt.Errorf("limit不能为负数")
	}

	if priceOrdersOffset < 0 {
		return fmt.Errorf("offset不能为负数")
	}

	result, err := trader.GetPriceOrders(
		cfg.APIKey, cfg.APISecret, priceOrdersSettle, priceOrdersStatus, priceOrdersContract,
		priceOrdersLimit, priceOrdersOffset,
	)
	if err != nil {
		return fmt.Errorf("获取自动订单失败: %v", err)
	}

	printJSON(result)
	return nil
}
