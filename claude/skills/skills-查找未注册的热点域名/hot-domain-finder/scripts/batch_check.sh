#!/usr/bin/env bash
# batch_check.sh - 批量检测域名注册状态(并发)
# 用法:
#   ./batch_check.sh domains.txt          # 从文件读,每行一个域名
#   echo -e "a.com\nb.io" | ./batch_check.sh -   # 从 stdin 读
# 输出(TSV,含表头): domain  status  method
# 并发度默认 6,可用环境变量 CONCURRENCY 覆盖。

set -uo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CHECK="$SCRIPT_DIR/check_domain.sh"
CONCURRENCY="${CONCURRENCY:-6}"

src="${1:-}"
if [[ -z "$src" ]]; then
  echo "usage: $0 <domains.txt|->" >&2
  exit 2
fi

if [[ "$src" == "-" ]]; then
  input="$(cat)"
else
  input="$(cat "$src")"
fi

# 去重、去空行、小写
domains="$(echo "$input" | tr '[:upper:]' '[:lower:]' | sed '/^[[:space:]]*$/d' | awk '{$1=$1};1' | sort -u)"

echo -e "domain\tstatus\tmethod"

# 用 xargs 并发调用单域名脚本
echo "$domains" | xargs -P "$CONCURRENCY" -I {} bash "$CHECK" {}
