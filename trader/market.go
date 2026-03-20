package trader

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type KlineData struct {
	Timestamp   int64   `json:"t"`
	VolumeQuote float64 `json:"v"`
	Close       string  `json:"c"`
	High        string  `json:"h"`
	Low         string  `json:"l"`
	Open        string  `json:"o"`
	VolumeBase  string  `json:"sum"`
	IsClosed    string  `json:"-"`
}

type TickerInfo struct {
	Contract       string `json:"contract"`
	Last           string `json:"last"`
	ChangePercent  string `json:"change_percentage"`
	Volume24h      string `json:"volume_24h"`
	VolumeQuote24h string `json:"volume_24h_quote"`
	High24h        string `json:"high_24h"`
	Low24h         string `json:"low_24h"`
	MarkPrice      string `json:"mark_price"`
	FundingRate    string `json:"funding_rate"`
	IndexPrice     string `json:"index_price"`
	TotalSize      string `json:"total_size"`
}

type FundingRateInfo struct {
	Contract        string `json:"name"`
	FundingRate     string `json:"funding_rate"`
	MarkPrice       string `json:"mark_price"`
	IndexPrice      string `json:"index_price"`
	NextFundingTime int64  `json:"funding_next_apply"`
}

func GetKlineData(contract, interval string, limit int) ([]KlineData, error) {
	validIntervals := []string{"1m", "5m", "15m", "30m", "1h", "2h", "4h", "6h", "12h", "1d", "3d", "7d"}
	if !IsValidInterval(interval, validIntervals) {
		return nil, fmt.Errorf("interval必须是: %v", validIntervals)
	}

	if limit < 1 || limit > 2000 {
		return nil, fmt.Errorf("limit必须在1到2000之间")
	}

	if !IsValidContract(contract) {
		return nil, fmt.Errorf("contract格式应为'基础货币_结算货币'")
	}

	queryParams := fmt.Sprintf("contract=%s&interval=%s&limit=%d",
		url.QueryEscape(contract), url.QueryEscape(interval), limit)

	fullURL := HOST + PREFIX + "/futures/usdt/candlesticks?" + queryParams
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("获取K线数据失败，状态码: %d", resp.StatusCode)
	}

	var result []KlineData
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func GetTicker(contract string) (*TickerInfo, error) {
	if contract == "" {
		return nil, fmt.Errorf("contract不能为空")
	}

	queryParams := "contract=" + url.QueryEscape(contract)
	fullURL := HOST + PREFIX + "/futures/usdt/tickers?" + queryParams
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("获取行情数据失败，状态码: %d", resp.StatusCode)
	}

	var result []TickerInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("未找到合约: %s", contract)
	}

	return &result[0], nil
}

func GetMultipleTickers(contracts []string) ([]TickerInfo, error) {
	if len(contracts) == 0 {
		return nil, fmt.Errorf("contracts不能为空")
	}

	var result []TickerInfo
	for _, contract := range contracts {
		ticker, err := GetTicker(contract)
		if err != nil {
			return nil, fmt.Errorf("获取 %s 行情失败: %v", contract, err)
		}
		result = append(result, *ticker)
	}

	return result, nil
}

func GetFundingRate(contract string) (*FundingRateInfo, error) {
	if !IsValidContract(contract) {
		return nil, fmt.Errorf("contract格式应为'基础货币_结算货币'")
	}

	fullURL := HOST + PREFIX + "/futures/usdt/contracts/" + url.QueryEscape(contract)
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("获取资金费率失败，状态码: %d", resp.StatusCode)
	}

	var result FundingRateInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func GetMultipleFundingRates(contracts []string) ([]FundingRateInfo, error) {
	if len(contracts) == 0 {
		return nil, fmt.Errorf("contracts不能为空")
	}

	var result []FundingRateInfo
	for _, contract := range contracts {
		fundingRate, err := GetFundingRate(contract)
		if err != nil {
			return nil, fmt.Errorf("获取 %s 资金费率失败: %v", contract, err)
		}
		result = append(result, *fundingRate)
	}

	return result, nil
}

func IsValidInterval(interval string, valid []string) bool {
	for _, v := range valid {
		if interval == v {
			return true
		}
	}
	return false
}
