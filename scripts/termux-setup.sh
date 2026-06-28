#!/usr/bin/env bash
# Termux (Android) setup + run for Nova.
#
# Builds Nova NATIVELY inside Termux (CGO on by default → links Bionic libc →
# working DNS for outbound LLM API calls). A cross-compiled CGO_ENABLED=0
# static binary would serve the page but break every outbound API call, because
# Go's netgo resolver reads /etc/resolv.conf which Termux/Android lacks. So we
# build here, not on a laptop.
#
# Run from anywhere:   bash scripts/termux-setup.sh
# Then open the phone browser at: http://127.0.0.1:8080
set -e

PORT="${NOVA_PORT:-8080}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "==> 安装依赖（golang）"
pkg update -y >/dev/null
pkg install -y golang >/dev/null

RUN_DIR="$HOME/nova-run"
mkdir -p "$RUN_DIR"

echo "==> 在 Termux 内原生构建 Nova（CGO 开启 → DNS 正常）"
cd "$REPO_DIR"
go build -ldflags "-s -w" -o "$RUN_DIR/nova" ./cmd/nova/

echo "==> 获取 wake-lock（防止息屏被系统杀死）"
if command -v termux-wake-lock >/dev/null 2>&1; then
  termux-wake-lock
else
  echo "  termux-wake-lock 不可用；建议 pkg install termux-api 并安装 termux:API 配套应用，"
  echo "  并在 系统设置→应用→Termux→电池 设为「无限制」，保持 Termux 通知不被划掉。"
fi

echo ""
echo "==> Nova 已就绪"
echo "  手机浏览器打开: http://127.0.0.1:${PORT}"
echo "  停止服务: 在 Termux 按 Ctrl+C，再执行 termux-wake-unlock"
echo "  再次启动: cd \"${RUN_DIR}\" && ./nova --no-open --port ${PORT}"
echo ""
echo "------------------------------------------------------------"
cd "$RUN_DIR"
exec ./nova --no-open --port "$PORT"
