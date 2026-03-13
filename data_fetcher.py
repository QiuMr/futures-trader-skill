import requests
import pandas as pd
import json
import time
import threading

from datetime import datetime

# 全局变量存储最新数据
market_data = None
data_lock = threading.Lock()


def get_kline_data(currency_pair: str, interval: str, limit: int):
    """
    获取指定交易对的K线数据

    Args:
        currency_pair (str): 交易对，格式为 '基础货币_计价货币'，例如 'BTC_USDT'
        interval (str): K线时间间隔，支持以下格式：
            - 分钟级：1m, 5m, 15m, 30m
            - 小时级：1h, 2h, 4h, 6h, 12h
            - 日级：1d, 3d
        limit (int): 获取的K线数量，最大值为1000

    Returns:
        str: JSON格式的K线数据，包含以下字段：
            - 时间戳：北京时间格式化的时间字符串
            - 开盘价：该时间段的起始价格
            - 最高价：该时间段的最高价格
            - 最低价：该时间段的最低价格
            - 收盘价：该时间段的结束价格
            - 基础货币成交量：以基础货币计量的交易量
            - 计价货币成交额：以计价货币计量的交易额
            - 是否闭合：标记该K线是否已完成
    """
    url = "https://api.gateio.ws/api/v4/spot/candlesticks"
    params = {"currency_pair": currency_pair, "interval": interval, "limit": limit}

    response = requests.get(url, params=params)
    data = response.json()

    # 构建 DataFrame
    df = pd.DataFrame(
        data,
        columns=[
            "时间戳",
            "计价货币成交额",
            "收盘价",
            "最高价",
            "最低价",
            "开盘价",
            "基础货币成交量",
            "是否闭合",
        ],
    )

    # 转换时间格式
    df["时间戳"] = (
        pd.to_datetime(
            df["时间戳"].astype("int64"), unit="s", utc=True
        )  # 先解析为 UTC 时间
        .dt.tz_convert("Asia/Shanghai")  # 转换为北京时间（UTC+8）
        .dt.strftime("%Y-%m-%d %H:%M:%S")  # 格式化输出
    )

    # 数值列转 float
    numeric_cols = [
        "计价货币成交额",
        "收盘价",
        "最高价",
        "最低价",
        "开盘价",
        "基础货币成交量",
    ]
    df[numeric_cols] = df[numeric_cols].astype(float)

    # 返回 JSON 格式
    return df.to_json(orient="records", force_ascii=False, indent=2)


def get_contract_volume(symbols):
    """
    获取指定币种的成交量数据
    Args:
        symbols: 币种名称列表，如 ['BTC_USDT', 'ETH_USDT']
    Returns:
        List[Dict]: 包含币种名称和24h成交额的字典列表
    """
    base_url = "https://api.gateio.ws/api/v4"
    settle = "usdt"

    try:
        # 获取行情数据
        resp = requests.get(f"{base_url}/futures/{settle}/tickers", timeout=10)
        resp.raise_for_status()
        tickers = resp.json()

        result = []
        # 将tickers转换为字典便于查找
        ticker_dict = {t["contract"]: t for t in tickers}

        for symbol in symbols:
            ticker = ticker_dict.get(symbol)
            if ticker:
                result.append(
                    {"symbol": symbol, "volume_24h": ticker["volume_24h_quote"]}
                )

        return result

    except requests.RequestException as e:
        print(f"请求失败: {e}")
        return []


def get_daily_market_data(symbols, interval, limit):
    """
    获取多个交易对的关键数据
    Args:
        symbols: 交易对列表，如 ['BTC_USDT', 'ETH_USDT']
        interval: 时间间隔，如 '1d', '1h', '1m'
        limit: 获取的数据条数
    Returns:
        str: 压缩格式的字符串，每个交易对的数据单独显示
    """
    table_lines = []

    def format_time(time_str, interval):
        """使用datetime库格式化时间"""
        dt = datetime.strptime(time_str, "%Y-%m-%d %H:%M:%S")
        if interval == "1d":
            return dt.strftime("%y-%m-%d")
        return dt.strftime("%y-%m-%d %H:%M")

    for symbol in symbols:
        # 获取数据
        data = get_kline_data(symbol, interval, limit)
        data = json.loads(data)

        if not data:
            continue

        # 获取时间范围
        start_time = format_time(data[0]["时间戳"], interval)
        end_time = format_time(data[-1]["时间戳"], interval)

        # 添加交易对标题和时间范围
        table_lines.append(f"{symbol} [{interval}] ({start_time}~{end_time}):")
        # 添加表头
        table_lines.append("时间\t最高价\t最低价\t成交额")

        for day_data in data:
            # 使用datetime格式化时间
            time_str = format_time(day_data["时间戳"], interval)

            # 将数据格式化为表格行
            line = f"{time_str}\t{float(day_data['最高价']):.2f}\t{float(day_data['最低价']):.2f}\t{int(day_data['计价货币成交额'])}"
            table_lines.append(line)

        # 每个交易对数据后添加空行，提高可读性
        table_lines.append("")

    # 用换行符连接所有行
    return "\n".join(table_lines)


def get_multi_interval_market_data1(symbols=None, intervals=None, limits=None):
    """
    获取多个交易对在不同时间间隔下的市场数据

    Args:
        symbols (List[str], optional): 交易对列表，默认为 ['BTC_USDT', 'ETH_USDT']
        intervals (List[str], optional): 时间间隔列表，默认为 ['1d', '1h', '1m']
        limits (Dict[str, int], optional): 各时间间隔的数据条数限制，
                                         默认为 {'1d': 30, '1h': 10, '1m': 60}

    Returns:
        Dict[str, str]: 返回一个字典，键为时间间隔，值为对应的市场数据字符串
    """
    # 设置默认参数
    if symbols is None:
        symbols = ["BTC_USDT", "ETH_USDT"]
    if intervals is None:
        intervals = ["7d", "3d","1d","15m"]
    if limits is None:
        limits = {"7d": 52, "3d": 30, "1d": 30, "15m": 30}

    results = ""
    for interval in intervals:
        results += get_daily_market_data(symbols, interval, limits[interval]) + "\n"
    return results


def get_multiple_funding_rates(
    contracts: list = ["BTC_USDT", "ETH_USDT", "DOGE_USDT", "SOL_USDT"],
):
    """
    获取多个合约的资金费率和合约乘数

    Args:
        contracts (list): 合约标识列表，如 ['BTC_USDT', 'ETH_USDT']

    Returns:
        dict: 合约标识到资金费率和合约乘数的映射，如 
              {'BTC_USDT': {'资金费率': '0.0001', '合约乘数': '0.0001'}, 
               'ETH_USDT': {'资金费率': '0.0002', '合约乘数': '0.001'}}
    """
    base_url = "https://api.gateio.ws/api/v4"
    settle = "usdt"

    try:
        # 获取合约信息
        resp = requests.get(f"{base_url}/futures/{settle}/contracts", timeout=10)
        resp.raise_for_status()
        contracts_info = resp.json()

        # 将合约信息转换为字典便于查找
        contract_dict = {c["name"]: c for c in contracts_info}

        # 构建结果字典
        result = {}
        for contract in contracts:
            if contract in contract_dict:
                result[contract] = {
                    "资金费率": contract_dict[contract]["funding_rate"],
                    "合约乘数": contract_dict[contract]["quanto_multiplier"]
                }
            else:
                result[contract] = {
                    "资金费率": None,
                    "合约乘数": None
                }

        return result

    except requests.RequestException as e:
        print(f"请求失败: {e}")
        return {contract: {"资金费率": None, "合约乘数": None} for contract in contracts}



def get_multi_interval_market_data():
    """
    获取市场数据
    Returns:
        str: 市场数据字符串
    """
    global market_data
    while True:
        with data_lock:
            if market_data is not None:
                return market_data
        time.sleep(1)  # 等待1000ms后重试


def update_data():
    global market_data
    while True:
        try:
            symbols = ["BTC_USDT"]
            intervals = ["7d", "3d","1d","15m"]
            limits = {"7d": 52, "3d": 30, "1d": 30, "15m": 30}

            results = ""
            for interval in intervals:
                results += (
                    get_daily_market_data(symbols, interval, limits[interval]) + "\n"
                )

            with data_lock:
                market_data = results
        except Exception as e:
            print(f"更新数据出错: {e}")
        time.sleep(300)  # 每5分钟更新一次


# 启动数据更新线程
thread = threading.Thread(target=update_data, daemon=True)
thread.start()


if __name__ == "__main__":
    market_data = get_multi_interval_market_data()
    print(market_data)

    