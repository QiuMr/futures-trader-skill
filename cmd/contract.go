package cmd

import (
	"encoding/json"
	"fmt"

	"futures-trader/trader"

	"github.com/spf13/cobra"
)

var contractCmd = &cobra.Command{
	Use:   "contract",
	Short: "查询合约信息",
	Long: `查询指定合约的详细信息，包括价格、杠杆、保证金等

参数说明:
  - settle: 结算货币 (btc/usdt)
  - contract: 合约标识 (如 BTC_USDT)

示例:
  futures-trader contract --settle usdt --contract BTC_USDT`,
	RunE: getContractInfo,
}

var (
	contractName string
	settleCoin   string
)

func init() {
	contractCmd.Flags().StringVar(&settleCoin, "settle", "usdt", "结算货币 (btc/usdt)")
	contractCmd.Flags().StringVar(&contractName, "contract", "", "合约标识 (如 BTC_USDT)")
	contractCmd.MarkFlagRequired("contract")
	rootCmd.AddCommand(contractCmd)
}

func getContractInfo(cmd *cobra.Command, args []string) error {
	info, err := trader.GetContractInfo(settleCoin, contractName)
	if err != nil {
		return fmt.Errorf("查询合约信息失败: %v", err)
	}

	data, _ := json.MarshalIndent(info, "", "  ")
	fmt.Println(string(data))

	return nil
}
