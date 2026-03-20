package cmd

import (
	"fmt"

	"futures-trader/config"
	"futures-trader/trader"

	"github.com/spf13/cobra"
)

var createOrderCmd = &cobra.Command{
	Use:   "create-order",
	Short: "创建期货订单",
	Long:  `创建期货交易订单，支持开仓和平仓操作

参数说明:
  - contract: 合约标识，格式为"基础货币_结算货币"，例如: BTC_USDT
  - size: 交易张数，正数表示开多/平空，负数表示开空/平多
  - price: 委托价格，市价单可设为"0"或不填
  - tif: 订单有效时间，支持 gtc(挂单取消), ioc(立即成交或取消), fok(全部成交或取消), poc(部分成交取消)

开仓/平仓模式:
  - 开仓: 只指定 size 和 price，不设置 close 或 reduce-only
  - 平仓: 设置 close=true 或 reduce-only=true
  - 双仓模式: 设置 auto-size="close_long" 或 "close_short" 指定平仓方向

示例:
  # 开多仓 1张BTC_USDT，价格50000
  futures-trader create-order --contract BTC_USDT --size 1 --price 50000

  # 开空仓 1张BTC_USDT，市价
  futures-trader create-order --contract BTC_USDT --size -1 --price 0

  # 平多仓（单仓模式）
  futures-trader create-order --contract BTC_USDT --size 1 --close

  # 平多仓（双仓模式）
  futures-trader create-order --contract BTC_USDT --size 1 --auto-size close_long`,
	RunE:  createOrder,
}

var (
	orderSettle     string
	orderContract   string
	orderSize       int
	orderPrice      string
	orderTif        string
	orderText       string
	orderReduceOnly bool
	orderClose      bool
	orderAutoSize   string
	orderStpAct     string
	orderIceberg    int
)

func init() {
	createOrderCmd.Flags().StringVar(&orderSettle, "settle", "usdt", "结算货币 (btc 或 usdt)")
	createOrderCmd.Flags().StringVar(&orderContract, "contract", "", "合约标识 (如 BTC_USDT)")
	createOrderCmd.Flags().IntVar(&orderSize, "size", 0, "交易张数 (正数=买，负数=卖)")
	createOrderCmd.Flags().StringVar(&orderPrice, "price", "", "委托价格 (市价单可设为0)")
	createOrderCmd.Flags().StringVar(&orderTif, "tif", "gtc", "订单类型 (gtc/ioc/fok/poc)")
	createOrderCmd.Flags().StringVar(&orderText, "text", "", "自定义订单ID (必须以t-开头)")
	createOrderCmd.Flags().BoolVar(&orderReduceOnly, "reduce-only", false, "只减仓模式")
	createOrderCmd.Flags().BoolVar(&orderClose, "close", false, "平仓模式")
	createOrderCmd.Flags().StringVar(&orderAutoSize, "auto-size", "", "双仓模式平仓方向 (close_long/close_short)")
	createOrderCmd.Flags().StringVar(&orderStpAct, "stp-act", "", "自成交策略 (co/cn/cb/-)")
	createOrderCmd.Flags().IntVar(&orderIceberg, "iceberg", 0, "冰山委托显示数量")
	createOrderCmd.MarkFlagRequired("contract")
	createOrderCmd.MarkFlagRequired("size")
	rootCmd.AddCommand(createOrderCmd)
}

func createOrder(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}
	if cfg == nil {
		return fmt.Errorf("未找到API密钥，请先使用 save-key 命令保存密钥")
	}

	if orderSettle != "btc" && orderSettle != "usdt" {
		return fmt.Errorf("settle必须是'btc'或'usdt'")
	}

	if !trader.IsValidContract(orderContract) {
		return fmt.Errorf("contract格式应为'基础货币_结算货币'，例如: BTC_USDT")
	}

	validTif := []string{"gtc", "ioc", "fok", "poc"}
	if !trader.IsValidTif(orderTif, validTif) {
		return fmt.Errorf("tif必须是: %v", validTif)
	}

	if orderAutoSize != "" && !trader.IsValidAutoSize(orderAutoSize) {
		return fmt.Errorf("auto-size必须是'close_long'或'close_short'")
	}

	if orderStpAct != "" && !trader.IsValidStpAct(orderStpAct) {
		return fmt.Errorf("stp-act必须是'co', 'cn', 'cb'或'-'")
	}

	if orderText != "" && !trader.IsValidText(orderText) {
		return fmt.Errorf("自定义ID必须以't-'开头且长度不超过28字节")
	}

	if orderClose && orderReduceOnly {
		return fmt.Errorf("close 和 reduce-only 不能同时为 true")
	}

	if orderClose && orderAutoSize != "" {
		return fmt.Errorf("close 和 auto-size 不能同时使用")
	}

	if orderReduceOnly && orderAutoSize != "" {
		return fmt.Errorf("reduce-only 和 auto-size 不能同时使用")
	}

	if orderAutoSize != "" && orderSize > 0 {
		return fmt.Errorf("auto-size 仅在平仓时使用，size 应为负数或0")
	}

	result, err := trader.CreateFuturesOrder(
		cfg.APIKey, cfg.APISecret, orderSettle, orderContract,
		orderSize, orderPrice, orderTif, orderText,
		orderReduceOnly, orderClose, orderAutoSize, orderStpAct, orderIceberg,
	)
	if err != nil {
		return fmt.Errorf("创建订单失败: %v", err)
	}

	printJSON(result)
	return nil
}
