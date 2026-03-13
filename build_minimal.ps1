# 编译最小体积的 Go 程序 (Windows)
go build -ldflags "-s -w -extldflags '-static'" -o futures-trader.exe .

# 使用 UPX 压缩
if (Get-Command upx -ErrorAction SilentlyContinue) {
    upx --best futures-trader.exe
}

Write-Host "Windows 编译完成" -ForegroundColor Green
