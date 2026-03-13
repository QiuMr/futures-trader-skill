# Gate.io 期货交易 CLI 工具

一个功能完整的命令行工具，用于与 Gate.io 期货交易平台进行交互。

## 功能特性

- ✅ API 密钥管理（保存/清除）
- ✅ 账户信息查询（USDT/BTC 结算账户）
- ✅ 持仓信息查询
- ✅ 合约信息查询（最新价、持仓量、杠杆、手续费等）
- ✅ 订单创建（开仓/平仓）
- ✅ 行情数据查询（K线、行情快照、资金费率）

## 安装

```bash
npx skills add [你的GitHub用户名]/[仓库名]
```

## 使用示例

### 1. 保存 API 密钥
```bash
futures-trader save-key --api-key YOUR_API_KEY --api-secret YOUR_API_SECRET
```

### 2. 查询合约信息
```bash
futures-trader contract --settle usdt --contract BTC_USDT
```

### 3. 查询账户余额
```bash
futures-trader account --settle usdt
```

### 4. 创建订单
```bash
futures-trader create-order --contract BTC_USDT --size 100 --price 70000
```

## 平台支持

- **Windows**: `futures-trader.txt` (约 2.7MB，UPX 压缩)
- **Linux**: `futures-trader-linux-amd64.txt` (约 2.6MB，UPX 压缩)

## 技术栈

- Go 语言
- Cobra CLI 框架
- Gate.io API

## 许可证

MIT
