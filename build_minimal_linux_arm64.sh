# 编译最小体积的 Go 程序 (Linux ARM64)
GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -extldflags '-static'" -o futures-trader-linux-arm64 .

# 使用 UPX 压缩
if command -v upx &> /dev/null; then
    upx --best futures-trader-linux-arm64
fi

echo "Linux ARM64 编译完成"
