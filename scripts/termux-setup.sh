#!/usr/bin/env bash
# Termux (Android) setup + run for Nova.
#
# Builds Nova AND its frontend NATIVELY inside Termux, then embeds the frontend
# (`-tags embedweb`) so a single ./nova binary serves everything — no separate
# web/ directory needed. Native build is required: a cross-compiled
# CGO_ENABLED=0 static binary breaks DNS (Go's netgo resolver reads
# /etc/resolv.conf, absent on Termux/Android), so every outbound LLM API call
# would fail. Building in Termux (CGO on → Bionic libc) gives working DNS.
#
# Run from anywhere:   bash scripts/termux-setup.sh
# Then open the phone browser at: http://127.0.0.1:8080
#
# Note: building the frontend (pnpm) on a phone is RAM/CPU heavy and can take
# several minutes the first time (dependency install). A modern arm64 phone
# handles it fine; re-runs reuse node_modules and are much faster.
set -e

PORT="${NOVA_PORT:-8080}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "==> 安装依赖（golang nodejs git）"
pkg update -y >/dev/null
pkg install -y golang nodejs git >/dev/null
corepack enable >/dev/null 2>&1 || true

echo "==> 构建前端（web/dist，首次较慢）"
pnpm --dir "$REPO_DIR/web" install
pnpm --dir "$REPO_DIR/web" build

echo "==> 准备内嵌前端资源（go:embed 源）"
rm -rf "$REPO_DIR/internal/webfs/dist"
cp -r "$REPO_DIR/web/dist" "$REPO_DIR/internal/webfs/dist"

RUN_DIR="$HOME/nova-run"
mkdir -p "$RUN_DIR"

echo "==> 在 Termux 内原生构建 Nova（内嵌前端，-tags embedweb）"
cd "$REPO_DIR"
go build -tags embedweb -ldflags "-s -w" -o "$RUN_DIR/nova" ./cmd/nova/

echo "==> 获取 wake-lock（防止息屏被系统杀死）"
if command -v termux-wake-lock >/dev/null 2>&1; then
  termux-wake-lock
else
  echo "  termux-wake-lock 不可用；建议 pkg install termux-api 并安装 termux:API 配套应用，"
  echo "  并在 系统设置→应用→Termux→电池 设为「无限制」，保持 Termux 通知不被划掉。"
fi

echo ""
echo "==> Nova 已就绪（单二进制，内嵌前端）"
echo "  手机浏览器打开: http://127.0.0.1:${PORT}"
echo "  停止服务: 在 Termux 按 Ctrl+C，再执行 termux-wake-unlock"
echo "  再次启动（无需重新构建）: cd \"${RUN_DIR}\" && ./nova --no-open --port ${PORT}"
echo ""
echo "------------------------------------------------------------"
cd "$RUN_DIR"
exec ./nova --no-open --port "$PORT"
