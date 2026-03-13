package trader

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	HOST   = "https://api.gateio.ws"
	PREFIX = "/api/v4"
)

type InitialOrder struct {
	Contract   string `json:"contract"`
	Size       int    `json:"size"`
	Price      string `json:"price"`
	Tif        string `json:"tif"`
	Text       string `json:"text"`
	Close      bool   `json:"close,omitempty"`
	ReduceOnly bool   `json:"reduce_only,omitempty"`
	AutoSize   string `json:"auto_size,omitempty"`
}

func GenSign(apiKey, apiSecret, method, path, query, payload string) (map[string]string, error) {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	m := sha512.New()
	m.Write([]byte(payload))
	hashedPayload := fmt.Sprintf("%x", m.Sum(nil))

	data := method + "\n" + path + "\n" + query + "\n" + hashedPayload + "\n" + timestamp
	h := hmac.New(sha512.New, []byte(apiSecret))
	h.Write([]byte(data))
	signature := fmt.Sprintf("%x", h.Sum(nil))

	headers := map[string]string{
		"KEY":       apiKey,
		"Timestamp": timestamp,
		"SIGN":      signature,
	}

	return headers, nil
}

func CreateFuturesOrder(apiKey, apiSecret, settle, contract string, size int, price string, tif, text string, reduceOnly, close bool, autoSize, stpAct string, iceberg int) (map[string]interface{}, error) {
	orderData := map[string]interface{}{
		"contract":    contract,
		"size":        size,
		"iceberg":     iceberg,
		"reduce_only": reduceOnly,
		"close":       close,
		"tif":         tif,
	}

	if price != "" {
		orderData["price"] = price
	}
	if text != "" {
		if !strings.HasPrefix(text, "t-") || len(text[2:]) > 28 {
			return nil, fmt.Errorf("自定义ID必须以't-'开头且长度不超过28字节")
		}
		orderData["text"] = text
	}
	if autoSize != "" {
		if autoSize != "close_long" && autoSize != "close_short" {
			return nil, fmt.Errorf("auto_size必须是'close_long'或'close_short'")
		}
		orderData["auto_size"] = autoSize
	}
	if stpAct != "" {
		if stpAct != "co" && stpAct != "cn" && stpAct != "cb" && stpAct != "-" {
			return nil, fmt.Errorf("stp_act必须是'co', 'cn', 'cb'或'-'")
		}
		orderData["stp_act"] = stpAct
	}

	payloadBytes, err := json.Marshal(orderData)
	if err != nil {
		return nil, err
	}
	payloadString := string(payloadBytes)

	signHeaders, err := GenSign(apiKey, apiSecret, "POST", PREFIX+"/futures/"+settle+"/orders", "", payloadString)
	if err != nil {
		return nil, err
	}

	fullURL := HOST + PREFIX + "/futures/" + settle + "/orders"
	req, err := http.NewRequest("POST", fullURL, strings.NewReader(payloadString))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("KEY", signHeaders["KEY"])
	req.Header.Set("Timestamp", signHeaders["Timestamp"])
	req.Header.Set("SIGN", signHeaders["SIGN"])

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("下单失败，状态码: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func GetFuturesAccountBalance(apiKey, apiSecret, settle string) (map[string]interface{}, error) {
	signHeaders, err := GenSign(apiKey, apiSecret, "GET", PREFIX+"/futures/"+settle+"/accounts", "", "")
	if err != nil {
		return nil, err
	}

	fullURL := HOST + PREFIX + "/futures/" + settle + "/accounts"
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("KEY", signHeaders["KEY"])
	req.Header.Set("Timestamp", signHeaders["Timestamp"])
	req.Header.Set("SIGN", signHeaders["SIGN"])

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("获取账户信息失败，状态码: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func GetRealPositions(apiKey, apiSecret, settle string) ([]map[string]interface{}, error) {
	queryParams := "holding=true"
	signHeaders, err := GenSign(apiKey, apiSecret, "GET", PREFIX+"/futures/"+settle+"/positions", queryParams, "")
	if err != nil {
		return nil, err
	}

	fullURL := HOST + PREFIX + "/futures/" + settle + "/positions?" + queryParams
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("KEY", signHeaders["KEY"])
	req.Header.Set("Timestamp", signHeaders["Timestamp"])
	req.Header.Set("SIGN", signHeaders["SIGN"])

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("获取持仓信息失败，状态码: %d", resp.StatusCode)
	}

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func CreatePriceTriggerOrder(apiKey, apiSecret, settle, contract, orderType, triggerPrice string, size int, price string, triggerStrategyType, triggerPriceType, triggerRule int, expiration int, tif, text string, close, reduceOnly bool, autoSize string) (map[string]interface{}, error) {
	if settle != "btc" && settle != "usdt" {
		return nil, fmt.Errorf("settle必须是'btc'或'usdt'")
	}

	if !strings.Contains(contract, "_") {
		return nil, fmt.Errorf("contract格式应为'基础货币_结算货币'")
	}

	validOrderTypes := []string{
		"close-long-position",
		"close-short-position",
		"plan-close-long-position",
		"plan-close-short-position",
	}

	valid := false
	for _, t := range validOrderTypes {
		if t == orderType {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fmt.Errorf("order_type必须是: %v", validOrderTypes)
	}

	_, err := strconv.ParseFloat(triggerPrice, 64)
	if err != nil {
		return nil, fmt.Errorf("trigger_price必须是数字字符串")
	}

	if triggerStrategyType != 0 && triggerStrategyType != 1 {
		return nil, fmt.Errorf("trigger_strategy_type必须是0或1")
	}

	if triggerPriceType < 0 || triggerPriceType > 2 {
		return nil, fmt.Errorf("trigger_price_type必须是0, 1或2")
	}

	if triggerRule != 1 && triggerRule != 2 {
		return nil, fmt.Errorf("trigger_rule必须是1或2")
	}

	if expiration < 60 || expiration > 2592000 {
		return nil, fmt.Errorf("expiration必须在60到2592000秒之间")
	}

	validTif := []string{"gtc", "ioc"}
	if !IsValidTif(tif, validTif) {
		return nil, fmt.Errorf("tif必须是: %v", validTif)
	}

	if text != "" && text != "api" {
		return nil, fmt.Errorf("text只能是空字符串或'api'")
	}

	orderData := map[string]interface{}{
		"contract":              contract,
		"order_type":            orderType,
		"trigger_price":         triggerPrice,
		"size":                  size,
		"price":                 price,
		"trigger_strategy_type": triggerStrategyType,
		"trigger_price_type":    triggerPriceType,
		"trigger_rule":          triggerRule,
		"expiration":            expiration,
		"tif":                   tif,
	}

	if text != "" {
		orderData["text"] = text
	}
	if close {
		orderData["close"] = close
	}
	if reduceOnly {
		orderData["reduce_only"] = reduceOnly
	}
	if autoSize != "" {
		orderData["auto_size"] = autoSize
	}

	payloadBytes, err := json.Marshal(orderData)
	if err != nil {
		return nil, err
	}
	payloadString := string(payloadBytes)

	signHeaders, err := GenSign(apiKey, apiSecret, "POST", PREFIX+"/futures/"+settle+"/price_orders", "", payloadString)
	if err != nil {
		return nil, err
	}

	fullURL := HOST + PREFIX + "/futures/" + settle + "/price_orders"
	req, err := http.NewRequest("POST", fullURL, strings.NewReader(payloadString))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("KEY", signHeaders["KEY"])
	req.Header.Set("Timestamp", signHeaders["Timestamp"])
	req.Header.Set("SIGN", signHeaders["SIGN"])

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("创建价格触发订单失败，状态码: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func GetPriceOrders(apiKey, apiSecret, settle, status, contract string, limit, offset int) ([]map[string]interface{}, error) {
	validStatus := []string{"open", "closed", "cancelled"}
	if !IsValidStatus(status, validStatus) {
		return nil, fmt.Errorf("status必须是: %v", validStatus)
	}

	queryParams := "status=" + status
	if contract != "" {
		queryParams += "&contract=" + url.QueryEscape(contract)
	}
	if limit > 0 {
		queryParams += "&limit=" + strconv.Itoa(limit)
	}
	if offset > 0 {
		queryParams += "&offset=" + strconv.Itoa(offset)
	}

	signHeaders, err := GenSign(apiKey, apiSecret, "GET", PREFIX+"/futures/"+settle+"/price_orders", queryParams, "")
	if err != nil {
		return nil, err
	}

	fullURL := HOST + PREFIX + "/futures/" + settle + "/price_orders?" + queryParams
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("KEY", signHeaders["KEY"])
	req.Header.Set("Timestamp", signHeaders["Timestamp"])
	req.Header.Set("SIGN", signHeaders["SIGN"])

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("查询自动订单失败，状态码: %d", resp.StatusCode)
	}

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func CancelAllPriceOrders(apiKey, apiSecret, settle, contract string) ([]map[string]interface{}, error) {
	queryParams := ""
	if contract != "" {
		queryParams = "contract=" + url.QueryEscape(contract)
	}

	signHeaders, err := GenSign(apiKey, apiSecret, "DELETE", PREFIX+"/futures/"+settle+"/price_orders", queryParams, "")
	if err != nil {
		return nil, err
	}

	fullURL := HOST + PREFIX + "/futures/" + settle + "/price_orders"
	if queryParams != "" {
		fullURL += "?" + queryParams
	}

	req, err := http.NewRequest("DELETE", fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("KEY", signHeaders["KEY"])
	req.Header.Set("Timestamp", signHeaders["Timestamp"])
	req.Header.Set("SIGN", signHeaders["SIGN"])

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("批量取消自动订单失败，状态码: %d", resp.StatusCode)
	}

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func IsValidContract(contract string) bool {
	if contract == "" {
		return false
	}
	return strings.Contains(contract, "_")
}

func IsValidTif(tif string, valid []string) bool {
	for _, v := range valid {
		if tif == v {
			return true
		}
	}
	return false
}

func IsValidStatus(status string, valid []string) bool {
	for _, v := range valid {
		if status == v {
			return true
		}
	}
	return false
}

func IsValidAutoSize(autoSize string) bool {
	return autoSize == "close_long" || autoSize == "close_short"
}

func IsValidStpAct(stpAct string) bool {
	return stpAct == "co" || stpAct == "cn" || stpAct == "cb" || stpAct == "-"
}

func IsValidText(text string) bool {
	if text == "" {
		return true
	}
	return strings.HasPrefix(text, "t-") && len(text) >= 3 && len(text) <= 31
}

func GetContractInfo(settle, contract string) (map[string]interface{}, error) {
	if settle != "btc" && settle != "usdt" {
		return nil, fmt.Errorf("settle必须是'btc'或'usdt'")
	}

	if !strings.Contains(contract, "_") {
		return nil, fmt.Errorf("contract格式应为'基础货币_结算货币'")
	}

	fullURL := HOST + PREFIX + "/futures/" + settle + "/contracts/" + contract
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("查询合约信息失败，状态码: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
