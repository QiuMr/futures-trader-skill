package cmd

import (
	"fmt"
	"strings"
	"time"

	"futures-trader/trader"

	"github.com/spf13/cobra"
)

var marketCmd = &cobra.Command{
	Use:   "market",
	Short: "查询行情数据",
	Long: `查询 Gate.io 期货行情数据，包括K线、行情快照、资金费率等

子命令:
  kline      查询K线数据
  ticker     查询行情快照
  funding    查询资金费率

示例:
  # 查询BTC_USDT的K线数据
  futures-trader market kline --contract BTC_USDT --interval 1h --limit 24

  # 查询多个交易对的行情快照
  futures-trader market ticker --contracts BTC_USDT,ETH_USDT

  # 查询资金费率
  futures-trader market funding --contracts BTC_USDT,ETH_USDT`,
}

func init() {
	rootCmd.AddCommand(marketCmd)
}

var klineCmd = &cobra.Command{
	Use:   "kline",
	Short: "查询K线数据",
	Long: `查询指定交易对的K线数据

参数说明:
  - contract: 合约标识，格式为"基础货币_结算货币"，例如: BTC_USDT
  - interval: K线时间间隔，支持: 1m, 5m, 15m, 30m, 1h, 2h, 4h, 6h, 12h, 1d, 3d, 7d
  - limit: 获取的K线数量，最大值为2000

示例:
  # 查询1小时K线，24条
  futures-trader market kline --contract BTC_USDT --interval 1h --limit 24

  # 查询1分钟K线，60条
  futures-trader market kline --contract ETH_USDT --interval 1m --limit 60`,
	RunE: getKline,
}

var (
	klineContract string
	klineInterval string
	klineLimit    int
)

func init() {
	klineCmd.Flags().StringVar(&klineContract, "contract", "", "合约标识 (必填)")
	klineCmd.Flags().StringVar(&klineInterval, "interval", "1h", "K线时间间隔 (默认: 1h)")
	klineCmd.Flags().IntVar(&klineLimit, "limit", 100, "K线数量 (默认: 100)")
	klineCmd.MarkFlagRequired("contract")
	marketCmd.AddCommand(klineCmd)
}

func getKline(cmd *cobra.Command, args []string) error {
	if klineContract == "" {
		return fmt.Errorf("contract不能为空")
	}

	if klineInterval == "" {
		return fmt.Errorf("interval不能为空")
	}

	if klineLimit < 1 || klineLimit > 2000 {
		return fmt.Errorf("limit必须在1到2000之间")
	}

	result, err := trader.GetKlineData(klineContract, klineInterval, klineLimit)
	if err != nil {
		return fmt.Errorf("获取K线数据失败: %v", err)
	}

	printKlineData(result)
	return nil
}

func printKlineData(data []trader.KlineData) {
	if len(data) == 0 {
		fmt.Println("无数据")
		return
	}

	fmt.Printf("%-19s %12s %12s %12s %12s %15s %12s\n",
		"时间", "开盘价", "最高价", "最低价", "收盘价", "成交量(张)", "成交额(USDT)")
	fmt.Println(strings.Repeat("-", 110))

	for _, d := range data {
		timestamp := time.Unix(d.Timestamp, 0).Format("2006-01-02 15:04:05")
		fmt.Printf("%-19s %12s %12s %12s %12s %15.2f %12s\n",
			timestamp, d.Open, d.High, d.Low, d.Close, d.VolumeQuote, d.VolumeBase)
	}
}

var tickerCmd = &cobra.Command{
	Use:   "ticker",
	Short: "查询行情快照",
	Long: `查询指定交易对的行情快照，包含最新价、涨跌幅、成交量等

参数说明:
  - contract: 合约标识，格式为"基础货币_结算货币"，例如: BTC_USDT
  - contracts: 多个合约标识，用逗号分隔，例如: BTC_USDT,ETH_USDT

示例:
  # 查询单个交易对行情
  futures-trader market ticker --contract BTC_USDT

  # 查询多个交易对行情
  futures-trader market ticker --contracts BTC_USDT,ETH_USDT,DOGE_USDT`,
	RunE: getTicker,
}

var (
	tickerContract  string
	tickerContracts string
)

func init() {
	tickerCmd.Flags().StringVar(&tickerContract, "contract", "", "合约标识")
	tickerCmd.Flags().StringVar(&tickerContracts, "contracts", "", "多个合约标识，用逗号分隔")
	marketCmd.AddCommand(tickerCmd)
}

func getTicker(cmd *cobra.Command, args []string) error {
	var contracts []string

	if tickerContracts != "" {
		contracts = strings.Split(tickerContracts, ",")
	} else if tickerContract != "" {
		contracts = []string{tickerContract}
	} else {
		return fmt.Errorf("必须指定 contract 或 contracts 参数")
	}

	result, err := trader.GetMultipleTickers(contracts)
	if err != nil {
		return fmt.Errorf("获取行情数据失败: %v", err)
	}

	printTickerData(result)
	return nil
}

func printTickerData(data []trader.TickerInfo) {
	if len(data) == 0 {
		fmt.Println("无数据")
		return
	}

	fmt.Printf("%-15s %12s %10s %12s %15s %12s %12s\n",
		"合约", "最新价", "涨跌幅", "24h成交量", "24h成交额", "24h最高", "24h最低")
	fmt.Println(strings.Repeat("-", 95))

	for _, d := range data {
		fmt.Printf("%-15s %12s %10s %12s %15s %12s %12s\n",
			d.Contract, d.Last, d.ChangePercent, d.Volume24h, d.VolumeQuote24h, d.High24h, d.Low24h)
	}
}

var fundingCmd = &cobra.Command{
	Use:   "funding",
	Short: "查询资金费率",
	Long: `查询指定合约的资金费率信息

参数说明:
  - contract: 合约标识，格式为"基础货币_结算货币"，例如: BTC_USDT
  - contracts: 多个合约标识，用逗号分隔，例如: BTC_USDT,ETH_USDT

示例:
  # 查询单个合约资金费率
  futures-trader market funding --contract BTC_USDT

  # 查询多个合约资金费率
  futures-trader market funding --contracts BTC_USDT,ETH_USDT,DOGE_USDT`,
	RunE: getFunding,
}

var (
	fundingContract  string
	fundingContracts string
)

func init() {
	fundingCmd.Flags().StringVar(&fundingContract, "contract", "", "合约标识")
	fundingCmd.Flags().StringVar(&fundingContracts, "contracts", "", "多个合约标识，用逗号分隔")
	marketCmd.AddCommand(fundingCmd)
}

func getFunding(cmd *cobra.Command, args []string) error {
	var contracts []string

	if fundingContracts != "" {
		contracts = strings.Split(fundingContracts, ",")
	} else if fundingContract != "" {
		contracts = []string{fundingContract}
	} else {
		return fmt.Errorf("必须指定 contract 或 contracts 参数")
	}

	result, err := trader.GetMultipleFundingRates(contracts)
	if err != nil {
		return fmt.Errorf("获取资金费率失败: %v", err)
	}

	printFundingData(result)
	return nil
}

func printFundingData(data []trader.FundingRateInfo) {
	if len(data) == 0 {
		fmt.Println("无数据")
		return
	}

	fmt.Printf("%-15s %12s %12s %12s %20s\n",
		"合约", "资金费率", "标记价格", "指数价格", "下次结算时间")
	fmt.Println(strings.Repeat("-", 80))

	for _, d := range data {
		nextFundingTime := time.Unix(d.NextFundingTime, 0).Format("2006-01-02 15:04:05")
		fmt.Printf("%-15s %12s %12s %12s %20s\n",
			d.Contract, d.FundingRate, d.MarkPrice, d.IndexPrice, nextFundingTime)
	}
}
