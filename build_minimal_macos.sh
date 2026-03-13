# 编译最小体积的 Go 程序 (macOS)
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -extldflags '-static'" -o futures-trader-macos .

# 使用 UPX 压缩
if command -v upx &> /dev/null; then
    upx --best futures-trader-macos
fi

echo "macOS 编译完成"
