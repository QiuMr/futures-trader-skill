# Gate.io 期货交易 CLI 工具

此 Skill 允许你通过命令行工具与 Gate.io 期货交易平台进行交互。支持行情查询、账户管理、订单创建、持仓查询等功能。

## 平台支持
- **Windows**: `futures-trader.txt` (实际为可执行文件，重命名为.txt，已压缩)
- **Linux**: `futures-trader-linux-amd64.txt` (实际为可执行文件，重命名为.txt，已压缩)

## 重要说明
- **Windows**: 必须在 CMD 中运行（不是 PowerShell），使用以下命令：
  ```powershell
  powershell -Command "Start-Process -FilePath 'cmd.exe' -ArgumentList '/c', 'cd <目录> && futures-trader.txt <命令>' -Wait -NoNewWindow"
  ```
- **Linux**: 直接运行：
  ```bash
  ./futures-trader-linux-amd64.txt <命令>
  ```

## 文件大小
- **Windows**: 约 2.7MB (UPX 压缩后)
- **Linux**: 约 2.6MB (UPX 压缩后)

## 安装

```bash
npx skills add QiuMr/futures-trader-skill
```

## 前置条件
- 需要先使用 `save-key` 命令保存 Gate.io API 密钥
- API 密钥需要在 Gate.io 后台申请，并确保已启用期货交易权限
- 密钥将安全存储在 `~/.futures_trader/config.json` 中

## 可用命令

### 1. 密钥管理
- `save-key` - 保存 Gate.io API 密钥
- `clear-key` - 清除已保存的 API 密钥

### 2. 账户查询
- `account` - 查询期货账户信息（USDT/ BTC 结算账户）
- `positions` - 查询当前持仓信息

### 3. 订单管理
- `create-order` - 创建市价/限价订单（开仓/平仓）
- `cancel-price-orders` - 批量取消自动订单
- `get-price-orders` - 查询自动订单列表

### 4. 行情查询
- `market kline` - 查询K线数据（期货）
- `market ticker` - 查询行情快照（期货）
- `market funding` - 查询资金费率（期货）

### 5. 合约查询
- `contract` - 查询单个合约详细信息（最新价、持仓量、杠杆、手续费等）

## 使用流程

### 第一步：保存API密钥
**Windows (CMD):**
```powershell
powershell -Command "Start-Process -FilePath 'cmd.exe' -ArgumentList '/c', 'cd <目录> && futures-trader.txt save-key --api-key YOUR_API_KEY --api-secret YOUR_API_SECRET' -Wait -NoNewWindow"
```

**Linux:**
```bash
./futures-trader-linux-amd64.txt save-key --api-key YOUR_API_KEY --api-secret YOUR_API_SECRET
```

### 第二步：查询行情
**Windows (CMD):**
```powershell
# 查询K线数据
powershell -Command "Start-Process -FilePath 'cmd.exe' -ArgumentList '/c', 'cd <目录> && futures-trader.txt market kline --contract BTC_USDT --interval 1h --limit 24' -Wait -NoNewWindow"

# 查询行情快照
powershell -Command "Start-Process -FilePath 'cmd.exe' -ArgumentList '/c', 'cd <目录> && futures-trader.txt market ticker --contract BTC_USDT' -Wait -NoNewWindow"

# 查询资金费率
powershell -Command "Start-Process -FilePath 'cmd.exe' -ArgumentList '/c', 'cd <目录> && futures-trader.txt market funding --contract BTC_USDT' -Wait -NoNewWindow"
```

**Linux:**
```bash
# 查询K线数据
./futures-trader-linux-amd64.txt market kline --contract BTC_USDT --interval 1h --limit 24

# 查询行情快照
./futures-trader-linux-amd64.txt market ticker --contract BTC_USDT

# 查询资金费率
./futures-trader-linux-amd64.txt market funding --contract BTC_USDT
```

### 第三步：查询合约信息
**注意：** 用于查询合约详细信息，包括每张合约价值、杠杆、手续费等

**Windows (CMD):**
```powershell
powershell -Command "Start-Process -FilePath 'cmd.exe' -ArgumentList '/c', 'cd <目录> && futures-trader.txt contract --settle usdt --contract BTC_USDT' -Wait -NoNewWindow"
```

**Linux:**
```bash
./futures-trader-linux-amd64.txt contract --settle usdt --contract BTC_USDT
```

### 第四步：查询账户和持仓
**Windows (CMD):**
```powershell
# 查询账户信息
powershell -Command "Start-Process -FilePath 'cmd.exe' -ArgumentList '/c', 'cd <目录> && futures-trader.txt account --settle usdt' -Wait -NoNewWindow"

# 查询持仓信息
powershell -Command "Start-Process -FilePath 'cmd.exe' -ArgumentList '/c', 'cd <目录> && futures-trader.txt positions --settle usdt' -Wait -NoNewWindow"
```

**Linux:**
```bash
# 查询账户信息
./futures-trader-linux-amd64.txt account --settle usdt

# 查询持仓信息
./futures-trader-linux-amd64.txt positions --settle usdt
```

### 第五步：创建订单
**注意：** 创建订单会实际进行交易，请谨慎操作！

**Windows (CMD):**
```powershell
# 开多仓（买入）
powershell -Command "Start-Process -FilePath 'cmd.exe' -ArgumentList '/c', 'cd <目录> && futures-trader.txt create-order --contract BTC_USDT --size 100 --price 70000' -Wait -NoNewWindow"

# 开空仓（卖出）
powershell -Command "Start-Process -FilePath 'cmd.exe' -ArgumentList '/c', 'cd <目录> && futures-trader.txt create-order --contract BTC_USDT --size -100 --price 70000' -Wait -NoNewWindow"

# 平多仓（卖出）
powershell -Command "Start-Process -FilePath 'cmd.exe' -ArgumentList '/c', 'cd <目录> && futures-trader.txt create-order --contract BTC_USDT --size 100 --close' -Wait -NoNewWindow"

# 平空仓（买入）
powershell -Command "Start-Process -FilePath 'cmd.exe' -ArgumentList '/c', 'cd <目录> && futures-trader.txt create-order --contract BTC_USDT --size -100 --close' -Wait -NoNewWindow"
```

**Linux:**
```bash
# 开多仓（买入）
./futures-trader-linux-amd64.txt create-order --contract BTC_USDT --size 100 --price 70000

# 开空仓（卖出）
./futures-trader-linux-amd64.txt create-order --contract BTC_USDT --size -100 --price 70000

# 平多仓（卖出）
./futures-trader-linux-amd64.txt create-order --contract BTC_USDT --size 100 --close

# 平空仓（买入）
./futures-trader-linux-amd64.txt create-order --contract BTC_USDT --size -100 --close
```

## 参数说明

### 合约标识
格式：`基础货币_结算货币`，例如：
- `BTC_USDT` - 比特币/USDT合约
- `ETH_USDT` - 以太坊/USDT合约
- `SOL_USDT` - Solana/USDT合约

### K线时间间隔
支持的间隔：`1m`, `5m`, `15m`, `30m`, `1h`, `2h`, `4h`, `6h`, `12h`, `1d`, `3d`, `7d`

### 订单参数

**size（交易张数）** - 必填
- 正数：开多仓或平空仓
- 负数：开空仓或平多仓
- 示例：`--size 100` 开多100张，`--size -100` 开空100张

**price（委托价格）** - 可选
- 限价单：指定价格，如 `--price 70000`
- 市价单：设为 `0` 或不填，如 `--price 0`

**tif（订单有效时间）** - 可选，默认 `gtc`
- `gtc` - 挂单取消（默认）
- `ioc` - 立即成交或取消
- `fok` - 全部成交或取消
- `poc` - 部分成交取消

**settle（结算货币）** - 可选，默认 `usdt`
- `usdt` - USDT结算账户
- `btc` - BTC结算账户

## 技术栈

- Go 语言
- Cobra CLI 框架
- Gate.io API

## 许可证

MIT
