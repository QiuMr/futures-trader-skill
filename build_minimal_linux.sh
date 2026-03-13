# 编译最小体积的 Go 程序 (Linux)
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -extldflags '-static'" -o futures-trader-linux .

# 使用 UPX 压缩
if command -v upx &> /dev/null; then
    upx --best futures-trader-linux
fi

echo "Linux 编译完成"
