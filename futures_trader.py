import requests
import time
import hashlib
import hmac
import json

# API配置
HOST = "https://api.gateio.ws"
PREFIX = "/api/v4"


def _gen_sign(api_key, api_secret, method, url, query_string=None, payload_string=None):
    """
    生成Gate.io API请求签名

    Args:
        api_key (str): API密钥
        api_secret (str): API密钥
        method (str): HTTP方法
        url (str): 请求URL
        query_string (str, optional): 查询参数字符串
        payload_string (str, optional): 请求体字符串

    Returns:
        dict: 包含签名的请求头
    """
    t = time.time()
    m = hashlib.sha512()
    payload = payload_string.encode("utf-8") if payload_string else b""
    m.update(payload)
    hashed_payload = m.hexdigest()

    s = f"{method}\n{url}\n{query_string or ''}\n{hashed_payload}\n{t}"

    sign = hmac.new(
        api_secret.encode("utf-8"), s.encode("utf-8"), hashlib.sha512
    ).hexdigest()

    return {"KEY": api_key, "Timestamp": str(t), "SIGN": sign}


def create_futures_order(
    api_key,
    api_secret,
    settle,
    contract,
    size,
    price=None,
    tif="gtc",
    text=None,
    reduce_only=False,
    close=False,
    auto_size=None,
    stp_act=None,
    iceberg=0,
):
    """
    创建合约交易订单（开仓或平仓）

    开仓示例：
    create_futures_order(
        api_key="your_api_key",
        api_secret="your_api_secret",
        settle="usdt",
        contract="BTC_USDT",
        size=1,  # 正数表示买入开多
        price="113011.11",
        reduce_only=True,  # 只减仓模式
    )

    平仓示例：
    create_futures_order(
        api_key="your_api_key",
        api_secret="your_api_secret",
        settle="usdt",
        contract="BTC_USDT",
        size=-1,  # 负数表示卖出平多
        price="113011.11",
        reduce_only=True,  # 只减仓模式
    )

    市价单示例：
    create_futures_order(
        api_key="your_api_key",
        api_secret="your_api_secret",
        settle="usdt",
        contract="BTC_USDT",
        size=1,  # 正数表示买入开多
        reduce_only=True,  # 只减仓模式
        price="0",  # 注意：必须是字符串"0"
        tif="ioc"   # 市价单必须使用ioc模式
    )

    :param api_key: API密钥
    :param api_secret: API密钥
    :param settle: 结算货币 (btc 或 usdt)
    :param contract: 合约标识 (如 "FIL_USDT")
    :param size: 交易张数 (正数=买，负数=卖)
    :param price: 委托价格 (字符串格式，市价单可设为None)
    :param tif: 订单类型 (gtc/ioc/fok/poc, 默认gtc)
    :param text: 自定义订单ID (必须以"t-"开头)
    :param reduce_only: 只减仓模式 (默认False)
    :param close: 平仓模式 (需要同时设置size为0)
    :param auto_size: 双仓模式平仓方向 ("close_long"或"close_short")
    :param stp_act: 自成交策略 ("co"/"cn"/"cb"/"-")
    :param iceberg: 冰山委托显示数量 (0为完全不隐藏)
    :return: 订单详情字典
    """
    if not contract:
        raise ValueError("合约标识(contract)是必填参数")
    if size is None:
        raise ValueError("交易张数(size)是必填参数")

    order_data = {
        "contract": contract,
        "size": size,
        "iceberg": iceberg,
        "reduce_only": reduce_only,
        "close": close,
        "tif": tif,
    }

    if price is not None:
        order_data["price"] = str(price)
    if text is not None:
        if not text.startswith("t-") or len(text[2:]) > 28:
            raise ValueError("自定义ID必须以't-'开头且长度不超过28字节")
        order_data["text"] = text
    if auto_size is not None:
        if auto_size not in ["close_long", "close_short"]:
            raise ValueError("auto_size必须是'close_long'或'close_short'")
        order_data["auto_size"] = auto_size
    if stp_act is not None:
        if stp_act not in ["co", "cn", "cb", "-"]:
            raise ValueError("stp_act必须是'co', 'cn', 'cb'或'-'")
        order_data["stp_act"] = stp_act

    url_path = f"/futures/{settle}/orders"
    full_url = PREFIX + url_path
    payload_string = json.dumps(order_data)

    sign_headers = _gen_sign(
        api_key=api_key,
        api_secret=api_secret,
        method="POST",
        url=full_url,
        query_string="",
        payload_string=payload_string,
    )

    headers = {"Accept": "application/json", "Content-Type": "application/json"}
    headers.update(sign_headers)

    response = requests.post(HOST + full_url, headers=headers, data=payload_string)

    if response.status_code != 201:
        error_msg = f"下单失败! 状态码: {response.status_code}"
        try:
            error_details = response.json()
            error_msg += f", 错误信息: {error_details.get('message', '未知错误')}"
            if "label" in error_details:
                error_msg += f" [标签: {error_details['label']}]"
        except Exception:
            error_msg += f", 响应内容: {response.text}"
        raise Exception(error_msg)

    return response.json()


def get_futures_account_balance(api_key, api_secret, settle):
    """
    获取合约账户信息

    参数:
    :param api_key: API密钥
    :param api_secret: API密钥
    :param settle: 结算货币 (btc 或 usdt)

    返回:
    :return: 字典包含账户信息和可用余额
    """
    url_path = f"/futures/{settle}/accounts"
    full_url = PREFIX + url_path

    sign_headers = _gen_sign(
        api_key=api_key,
        api_secret=api_secret,
        method="GET",
        url=full_url,
        query_string="",
    )

    headers = {"Accept": "application/json", "Content-Type": "application/json"}
    headers.update(sign_headers)

    response = requests.get(HOST + full_url, headers=headers)

    if response.status_code != 200:
        error_msg = f"获取账户信息失败! 状态码: {response.status_code}"
        try:
            error_details = response.json()
            error_msg += f", 错误信息: {error_details.get('message', '未知错误')}"
            if "label" in error_details:
                error_msg += f" [标签: {error_details['label']}]"
        except Exception:
            error_msg += f", 响应内容: {response.text}"
        raise Exception(error_msg)

    return response.json()


def get_real_positions(api_key, api_secret, settle):
    """
    获取用户真实持仓（当前持仓不为0的仓位）

    参数:
    :param api_key: API密钥
    :param api_secret: API密钥
    :param settle: 结算货币 (btc 或 usdt)

    返回:
    :return: 持仓信息列表
    """
    url_path = f"/futures/{settle}/positions"
    full_url = PREFIX + url_path
    query_params = "holding=true"

    sign_headers = _gen_sign(
        api_key=api_key,
        api_secret=api_secret,
        method="GET",
        url=full_url,
        query_string=query_params,
    )

    headers = {"Accept": "application/json", "Content-Type": "application/json"}
    headers.update(sign_headers)

    response = requests.get(f"{HOST}{full_url}?{query_params}", headers=headers)

    if response.status_code != 200:
        error_msg = f"获取持仓信息失败! 状态码: {response.status_code}"
        try:
            error_details = response.json()
            error_msg += f", 错误信息: {error_details.get('message', '未知错误')}"
            if "label" in error_details:
                error_msg += f" [标签: {error_details['label']}]"
        except Exception:
            error_msg += f", 响应内容: {response.text}"
        raise Exception(error_msg)

    return response.json()


def create_price_trigger_order(
    api_key,
    api_secret,
    settle,
    contract,
    order_type,
    trigger_price,
    size=0,
    price="0",
    trigger_strategy_type=0,
    trigger_price_type=0,
    trigger_rule=2,
    expiration=86400,
    tif="ioc",
    text="api",
    close=False,
    reduce_only=False,
    auto_size=None,
):
    """
    创建价格触发订单（止损/止盈单）

    多单止损示例：
    create_price_trigger_order(
        api_key="your_api_key",
        api_secret="your_api_secret",
        settle="usdt",
        contract="BTC_USDT",
        order_type="close-long-position",
        trigger_price="113000",  # 低于当前价格设置止损
        price="0",  # 市价成交
        trigger_rule=2,  # 2表示小于等于触发价格时触发
        tif="ioc",
        close=True  # 单仓模式必须设置
    )

    空单止损示例：
    create_price_trigger_order(
        api_key="your_api_key",
        api_secret="your_api_secret",
        settle="usdt",
        contract="BTC_USDT",
        order_type="close-short-position",
        trigger_price="117000",  # 高于当前价格设置止损
        price="0",  # 市价成交
        trigger_rule=1,  # 1表示大于等于触发价格时触发
        tif="ioc",
        close=True  # 单仓模式必须设置
    )

    参数说明：
    :param api_key: API密钥
    :param api_secret: API密钥
    :param settle: 结算货币，必填
        - 可选值: "btc" 或 "usdt"
        - 示例: "usdt"

    :param contract: 合约标识，必填
        - 格式: "基础货币_结算货币"
        - 示例: "BTC_USDT"

    :param order_type: 订单类型，必填
        - "close-long-position": 仓位止盈止损，全部平多仓
        - "close-short-position": 仓位止盈止损，全部平空仓
        - "plan-close-long-position": 仓位计划止盈止损，全部或部分平多仓
        - "plan-close-short-position": 仓位计划止盈止损，全部或部分平空仓

    :param trigger_price: 触发价格，必填（字符串格式）
        - 止损示例：多单设低于当前价，空单设高于当前价
        - 止盈示例：多单设高于当前价，空单设低于当前价
        - 示例: "50000"

    :param size: 平仓数量，默认0（全部平仓）
        - 0: 全部平仓
        - 正数: 平空单（部分平空）
        - 负数: 平多单（部分平多）

    :param price: 成交价格，默认"0"（市价）
        - "0": 市价成交
        - 具体价格: 限价成交（字符串格式）

    :param trigger_strategy_type: 触发策略，默认0
        - 0: 价格触发
        - 1: 价差触发（暂不支持）

    :param trigger_price_type: 参考价格类型，默认0
        - 0: 最新成交价
        - 1: 标记价格
        - 2: 指数价格

    :param trigger_rule: 触发规则，默认2
        - 1: 大于等于触发价格时触发
        - 2: 小于等于触发价格时触发

    :param expiration: 有效期（秒），默认86400（24小时）

    :param tif: 订单类型，默认"ioc"
        - "gtc": GoodTillCancelled
        - "ioc": ImmediateOrCancelled（市价单必须使用）

    :param text: 订单来源，默认"api"
        - "api": API调用
        - "web": 网页
        - "app": 移动端

    :param close: 单仓模式全部平仓标志，默认False
        - 单仓模式全部平仓时必须设为True

    :param reduce_only: 自动减仓标志，默认False
        - True: 确保只平仓不开新仓

    :param auto_size: 双仓模式平仓方向，可选
        - "close_long": 平多头
        - "close_short": 平空头

    返回:
        dict: 包含订单ID的响应字典
        示例: {"id": 1432329}

    异常:
        ValueError: 参数验证失败
        Exception: API请求失败
    """
    if settle not in ["btc", "usdt"]:
        raise ValueError("settle必须是'btc'或'usdt'")

    if not contract or "_" not in contract:
        raise ValueError("contract格式应为'基础货币_结算货币'")

    valid_order_types = [
        "close-long-position",
        "close-short-position",
        "plan-close-long-position",
        "plan-close-short-position",
    ]
    if order_type not in valid_order_types:
        raise ValueError(f"order_type必须是: {', '.join(valid_order_types)}")

    if (
        not isinstance(trigger_price, str)
        or not trigger_price.replace(".", "").isdigit()
    ):
        raise ValueError("trigger_price必须是数字字符串")

    if trigger_strategy_type not in [0, 1]:
        raise ValueError("trigger_strategy_type必须是0或1")

    if trigger_price_type not in [0, 1, 2]:
        raise ValueError("trigger_price_type必须是0、1或2")

    if trigger_rule not in [1, 2]:
        raise ValueError("trigger_rule必须是1或2")

    if tif not in ["gtc", "ioc"]:
        raise ValueError("tif必须是'gtc'或'ioc'")

    if auto_size and auto_size not in ["close_long", "close_short"]:
        raise ValueError("auto_size必须是'close_long'或'close_short'")

    order_data = {
        "initial": {
            "contract": contract,
            "size": size,
            "price": price,
            "tif": tif,
            "text": text,
        },
        "trigger": {
            "strategy_type": trigger_strategy_type,
            "price_type": trigger_price_type,
            "price": trigger_price,
            "rule": trigger_rule,
            "expiration": expiration,
        },
        "order_type": order_type,
    }

    if close:
        order_data["initial"]["close"] = close
    if reduce_only:
        order_data["initial"]["reduce_only"] = reduce_only
    if auto_size:
        order_data["initial"]["auto_size"] = auto_size

    url_path = f"/futures/{settle}/price_orders"
    full_url = PREFIX + url_path
    payload_string = json.dumps(order_data)

    sign_headers = _gen_sign(
        api_key=api_key,
        api_secret=api_secret,
        method="POST",
        url=full_url,
        query_string="",
        payload_string=payload_string,
    )

    headers = {"Accept": "application/json", "Content-Type": "application/json"}
    headers.update(sign_headers)

    response = requests.post(HOST + full_url, headers=headers, data=payload_string)

    if response.status_code != 200:
        error_msg = f"创建价格触发订单失败! 状态码: {response.status_code}"
        try:
            error_details = response.json()
            error_msg += f", 错误信息: {error_details.get('message', '未知错误')}"
            if "label" in error_details:
                error_msg += f" [标签: {error_details['label']}]"
        except Exception:
            error_msg += f", 响应内容: {response.text}"
        raise Exception(error_msg)

    return response.json()


def get_price_orders(
    api_key, api_secret, settle, status, contract=None, limit=None, offset=None
):
    """
    查询自动订单列表（止损/止盈单）

    查询所有未完成的订单：
    get_price_orders(
        api_key="your_api_key",
        api_secret="your_api_secret",
        settle="usdt",
        status="open"
    )

    查询特定合约的未完成订单：
    get_price_orders(
        api_key="your_api_key",
        api_secret="your_api_secret",
        settle="usdt",
        status="open",
        contract="BTC_USDT"
    )

    查询已完成的订单：
    get_price_orders(
        api_key="your_api_key",
        api_secret="your_api_secret",
        settle="usdt",
        status="finished"
    )

    :param api_key: API密钥
    :param api_secret: API密钥
    :param settle: 结算货币，必填
        - 可选值: "btc" 或 "usdt"
        - 示例: "usdt"

    :param status: 订单状态，必填
        - "open": 未完成的订单
        - "finished": 已完成的订单
        - 示例: "open"

    :param contract: 合约标识，可选
        - 格式: "基础货币_结算货币"
        - 不指定则返回所有合约的订单
        - 示例: "BTC_USDT"

    :param limit: 返回数量限制，可选
        - 限制返回的订单数量
        - 示例: 10（最多返回10条）

    :param offset: 偏移量，可选
        - 从第几个订单开始返回
        - 用于分页查询
        - 示例: 0（从第一个开始）

    :return: 订单列表，每个订单包含：
        - id: 订单ID
        - initial: 初始订单信息
        - trigger: 触发条件
        - status: 订单状态
        - finish_as: 完成方式
        - create_time: 创建时间
        - finish_time: 完成时间
        - order_type: 订单类型
    """
    query_params = f"status={status}"
    if contract:
        query_params += f"&contract={contract}"
    if limit is not None:
        query_params += f"&limit={limit}"
    if offset is not None:
        query_params += f"&offset={offset}"

    url_path = f"/futures/{settle}/price_orders"
    full_url = PREFIX + url_path

    sign_headers = _gen_sign(
        api_key=api_key,
        api_secret=api_secret,
        method="GET",
        url=full_url,
        query_string=query_params,
    )

    headers = {"Accept": "application/json", "Content-Type": "application/json"}
    headers.update(sign_headers)

    response = requests.get(f"{HOST}{full_url}?{query_params}", headers=headers)

    if response.status_code != 200:
        error_msg = f"查询自动订单失败! 状态码: {response.status_code}"
        try:
            error_details = response.json()
            error_msg += f", 错误信息: {error_details.get('message', '未知错误')}"
            if "label" in error_details:
                error_msg += f" [标签: {error_details['label']}]"
        except Exception:
            error_msg += f", 响应内容: {response.text}"
        raise Exception(error_msg)

    return response.json()


def cancel_all_price_orders(api_key, api_secret, settle, contract=None):
    """
    批量取消自动订单（止损/止盈单）

    参数说明：
    :param api_key: API密钥
    :param api_secret: API密钥
    :param settle: 结算货币，必填
        - 可选值: "btc" 或 "usdt"
        - 示例: "usdt"

    :param contract: 合约标识，可选
        - 格式: "基础货币_结算货币"
        - 不指定则取消所有合约的订单
        - 示例: "BTC_USDT"

    返回:
        list: 已取消订单的列表，每个订单包含：
            - id: 订单ID
            - initial: 初始订单信息
            - trigger: 触发条件
            - status: 订单状态
            - finish_as: 完成方式
            - create_time: 创建时间
            - finish_time: 完成时间
            - order_type: 订单类型
            - user: 用户ID
            - trade_id: 交易ID
            - reason: 取消原因

    异常:
        ValueError: 参数验证失败
        Exception: API请求失败

    使用示例:
        # 示例1: 取消所有USDT合约的自动订单
        cancelled_orders = cancel_all_price_orders(
            api_key="your_api_key",
            api_secret="your_api_secret",
            settle="usdt"
        )

        # 示例2: 仅取消BTC_USDT合约的自动订单
        cancelled_orders = cancel_all_price_orders(
            api_key="your_api_key",
            api_secret="your_api_secret",
            settle="usdt",
            contract="BTC_USDT"
        )
    """
    if settle not in ["btc", "usdt"]:
        raise ValueError("settle必须是'btc'或'usdt'")

    if contract and "_" not in contract:
        raise ValueError("contract格式应为'基础货币_结算货币'")

    url_path = f"/futures/{settle}/price_orders"
    full_url = PREFIX + url_path

    query_params = ""
    if contract:
        query_params = f"contract={contract}"

    sign_headers = _gen_sign(
        api_key=api_key,
        api_secret=api_secret,
        method="DELETE",
        url=full_url,
        query_string=query_params,
    )

    headers = {"Accept": "application/json", "Content-Type": "application/json"}
    headers.update(sign_headers)

    request_url = f"{HOST}{full_url}"
    if query_params:
        request_url += f"?{query_params}"

    response = requests.delete(request_url, headers=headers)

    if response.status_code != 200:
        error_msg = f"批量取消自动订单失败! 状态码: {response.status_code}"
        try:
            error_details = response.json()
            error_msg += f", 错误信息: {error_details.get('message', '未知错误')}"
            if "label" in error_details:
                error_msg += f" [标签: {error_details['label']}]"
        except Exception:
            error_msg += f", 响应内容: {response.text}"
        raise Exception(error_msg)

    return response.json()


def get_my_trades(
    api_key,
    api_secret,
    settle,
    contract=None,
    order=None,
    limit=None,
    offset=None,
    last_id=None,
):
    """
    查询个人成交记录

    参数说明：
    :param api_key: API密钥
    :param api_secret: API密钥
    :param settle: 结算货币，必填
        - 可选值: "btc" 或 "usdt"
        - 示例: "usdt"

    :param contract: 合约标识，可选
        - 格式: "基础货币_结算货币"
        - 示例: "BTC_USDT"

    :param order: 委托ID，可选
        - 如果指定则返回该委托相关的成交记录
        - 示例: 12345678

    :param limit: 返回数量限制，可选
        - 限制返回的记录数量
        - 示例: 10（最多返回10条）

    :param offset: 偏移量，可选
        - 从第几条记录开始返回
        - 用于分页查询
        - 示例: 0（从第一条开始）

    :param last_id: 上个列表的最后一条记录的ID，可选
        - 用于分页查询
        - 注意：该参数不再继续支持，建议使用get_my_trades_timerange查询更久的数据

    返回:
        list: 成交记录列表，每条记录包含：
            - id: 成交记录ID
            - create_time: 成交时间
            - contract: 合约标识
            - order_id: 关联订单ID
            - size: 成交数量
            - close_size: 平仓数量
            - price: 成交价格
            - role: 成交角色（taker/maker）
            - text: 订单自定义信息
            - fee: 成交手续费
            - point_fee: 成交点卡手续费
    """
    if settle not in ["btc", "usdt"]:
        raise ValueError("settle必须是'btc'或'usdt'")

    query_params = []
    if contract:
        query_params.append(f"contract={contract}")
    if order is not None:
        query_params.append(f"order={order}")
    if limit is not None:
        query_params.append(f"limit={limit}")
    if offset is not None:
        query_params.append(f"offset={offset}")
    if last_id is not None:
        query_params.append(f"last_id={last_id}")

    query_string = "&".join(query_params)

    url_path = f"/futures/{settle}/my_trades"
    full_url = PREFIX + url_path

    sign_headers = _gen_sign(
        api_key=api_key,
        api_secret=api_secret,
        method="GET",
        url=full_url,
        query_string=query_string,
    )

    headers = {"Accept": "application/json", "Content-Type": "application/json"}
    headers.update(sign_headers)

    request_url = f"{HOST}{full_url}"
    if query_string:
        request_url += f"?{query_string}"

    response = requests.get(request_url, headers=headers)

    if response.status_code != 200:
        error_msg = f"查询成交记录失败! 状态码: {response.status_code}"
        try:
            error_details = response.json()
            error_msg += f", 错误信息: {error_details.get('message', '未知错误')}"
            if "label" in error_details:
                error_msg += f" [标签: {error_details['label']}]"
        except Exception:
            error_msg += f", 响应内容: {response.text}"
        raise Exception(error_msg)

    return response.json()


if __name__ == "__main__":
    # 配置API密钥
    API_KEY = "b4c72c3b0be13e98c038673a834e95dc"
    API_SECRET = "2be3a3378730d386108c8dbb0eec4ae31818dc54d6913fa65bbacfcc7c781a60"
    SETTLE = "usdt"

    # 获取余额
    balance = get_futures_account_balance(API_KEY, API_SECRET, SETTLE)
    print(f"可用余额: {balance}")

    # 获取持仓
    positions = get_real_positions(API_KEY, API_SECRET, SETTLE)
    print(f"持仓: {positions}")

    # 查询自动订单
    orders = get_price_orders(API_KEY, API_SECRET, SETTLE, "open")
    print(f"自动订单: {orders}")

    # # 查询最近的所有成交记录
    # trades = get_my_trades(API_KEY, API_SECRET, SETTLE,limit=20)
    # print(f"最近的成交记录: {trades}")
