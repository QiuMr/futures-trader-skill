package cmd

import (
	"fmt"

	"futures-trader/config"
	"futures-trader/trader"

	"github.com/spf13/cobra"
)

var createTriggerOrderCmd = &cobra.Command{
	Use:   "create-trigger-order",
	Short: "创建价格触发订单（止损/止盈单）",
	Long:  `创建价格触发订单，支持止损和止盈功能

参数说明:
  - contract: 合约标识，格式为"基础货币_结算货币"，例如: BTC_USDT
  - order-type: 订单类型，支持:
      close-long-position: 平多仓（市价）
      close-short-position: 平空仓（市价）
      plan-close-long-position: 计划平多仓（限价）
      plan-close-short-position: 计划平空仓（限价）
  - trigger-price: 触发价格
  - size: 平仓数量，0表示全部平仓
  - price: 成交价格，0表示市价，其他值表示限价
  - strategy-type: 触发策略，0=价格触发
  - price-type: 参考价格类型，0=最新成交价，1=买一价，2=卖一价
  - rule: 触发规则，1=大于等于，2=小于等于
  - expiration: 有效期（秒），默认86400（24小时）
  - tif: 订单类型，支持 gtc(挂单取消), ioc(立即成交或取消)

示例:
  # 止损平多仓，当价格低于45000时市价平仓
  futures-trader create-trigger-order --contract BTC_USDT --order-type close-long-position --trigger-price 45000 --size 1

  # 止盈平空仓，当价格低于55000时市价平仓
  futures-trader create-trigger-order --contract BTC_USDT --order-type close-short-position --trigger-price 55000 --size 1

  # 计划平多仓，当价格大于等于50000时限价50000平仓
  futures-trader create-trigger-order --contract BTC_USDT --order-type plan-close-long-position --trigger-price 50000 --price 50000 --size 1`,
	RunE:  createTriggerOrder,
}

var (
	triggerSettle       string
	triggerContract     string
	triggerOrderType    string
	triggerPrice        string
	triggerSize         int
	triggerPriceField   string
	triggerStrategyType int
	triggerPriceType    int
	triggerRule         int
	triggerExpiration   int
	triggerTif          string
	triggerText         string
	triggerClose        bool
	triggerReduceOnly   bool
	triggerAutoSize     string
)

func init() {
	createTriggerOrderCmd.Flags().StringVar(&triggerSettle, "settle", "usdt", "结算货币 (btc 或 usdt)")
	createTriggerOrderCmd.Flags().StringVar(&triggerContract, "contract", "", "合约标识 (如 BTC_USDT)")
	createTriggerOrderCmd.Flags().StringVar(&triggerOrderType, "order-type", "", "订单类型 (close-long-position/close-short-position/plan-close-long-position/plan-close-short-position)")
	createTriggerOrderCmd.Flags().StringVar(&triggerPrice, "trigger-price", "", "触发价格")
	createTriggerOrderCmd.Flags().IntVar(&triggerSize, "size", 0, "平仓数量 (0=全部平仓)")
	createTriggerOrderCmd.Flags().StringVar(&triggerPriceField, "price", "0", "成交价格 (0=市价)")
	createTriggerOrderCmd.Flags().IntVar(&triggerStrategyType, "strategy-type", 0, "触发策略 (0=价格触发)")
	createTriggerOrderCmd.Flags().IntVar(&triggerPriceType, "price-type", 0, "参考价格类型 (0=最新成交价)")
	createTriggerOrderCmd.Flags().IntVar(&triggerRule, "rule", 2, "触发规则 (1=大于等于 2=小于等于)")
	createTriggerOrderCmd.Flags().IntVar(&triggerExpiration, "expiration", 86400, "有效期 (秒)")
	createTriggerOrderCmd.Flags().StringVar(&triggerTif, "tif", "ioc", "订单类型 (gtc/ioc)")
	createTriggerOrderCmd.Flags().StringVar(&triggerText, "text", "api", "订单来源")
	createTriggerOrderCmd.Flags().BoolVar(&triggerClose, "close", false, "单仓模式全部平仓")
	createTriggerOrderCmd.Flags().BoolVar(&triggerReduceOnly, "reduce-only", false, "自动减仓")
	createTriggerOrderCmd.Flags().StringVar(&triggerAutoSize, "auto-size", "", "双仓模式平仓方向")
	createTriggerOrderCmd.MarkFlagRequired("contract")
	createTriggerOrderCmd.MarkFlagRequired("order-type")
	createTriggerOrderCmd.MarkFlagRequired("trigger-price")
	rootCmd.AddCommand(createTriggerOrderCmd)
}

func createTriggerOrder(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}
	if cfg == nil {
		return fmt.Errorf("未找到API密钥，请先使用 save-key 命令保存密钥")
	}

	if triggerClose && triggerAutoSize != "" {
		return fmt.Errorf("close 和 auto-size 不能同时使用")
	}

	if triggerReduceOnly && triggerAutoSize != "" {
		return fmt.Errorf("reduce-only 和 auto-size 不能同时使用")
	}

	result, err := trader.CreatePriceTriggerOrder(
		cfg.APIKey, cfg.APISecret, triggerSettle, triggerContract, triggerOrderType,
		triggerPrice, triggerSize, triggerPriceField, triggerStrategyType,
		triggerPriceType, triggerRule, triggerExpiration, triggerTif, triggerText,
		triggerClose, triggerReduceOnly, triggerAutoSize,
	)
	if err != nil {
		return fmt.Errorf("创建价格触发订单失败: %v", err)
	}

	printJSON(result)
	return nil
}
