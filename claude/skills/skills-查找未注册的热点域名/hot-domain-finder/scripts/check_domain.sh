#!/usr/bin/env bash
# check_domain.sh - 判断单个域名是否已被注册
# 用法: ./check_domain.sh example.com
# 输出: <domain>\t<REGISTERED|AVAILABLE|UNKNOWN>\t<method>
#
# 策略:
#   1. 通过 IANA RDAP bootstrap(data.iana.org/rdap/dns.json)动态解析该 TLD 的
#      权威 RDAP 端点(结果缓存到本地,24h 有效),避免硬编码端点失效。
#      HTTP 200=已注册,404=未注册。
#   2. 该 TLD 无 RDAP 端点或 RDAP 请求失败时,回退到 whois 关键字匹配。

set -uo pipefail

domain="${1:-}"
if [[ -z "$domain" ]]; then
  echo "usage: $0 <domain>" >&2
  exit 2
fi

domain="$(echo "$domain" | tr '[:upper:]' '[:lower:]' | xargs)"

# 基本格式校验:必须形如 label.tld,只允许字母数字与连字符,含至少一个点
if ! echo "$domain" | grep -qE '^[a-z0-9]([a-z0-9-]*[a-z0-9])?(\.[a-z0-9]([a-z0-9-]*[a-z0-9])?)+$'; then
  printf "%s\t%s\t%s\n" "$domain" "INVALID" "format"
  exit 0
fi

tld="${domain##*.}"

TIMEOUT=8
CACHE_DIR="${TMPDIR:-/tmp}/hot-domain-finder-cache"
BOOTSTRAP_CACHE="$CACHE_DIR/rdap-dns.json"
mkdir -p "$CACHE_DIR" 2>/dev/null || true

# --- 获取/刷新 IANA RDAP bootstrap 缓存 ---
fetch_bootstrap() {
  # 缓存存在且 24h 内则复用
  if [[ -f "$BOOTSTRAP_CACHE" ]]; then
    local age
    age=$(( $(date +%s) - $(stat -f %m "$BOOTSTRAP_CACHE" 2>/dev/null || stat -c %Y "$BOOTSTRAP_CACHE" 2>/dev/null || echo 0) ))
    if [[ "$age" -lt 86400 ]]; then
      return 0
    fi
  fi
  curl -s -m 10 https://data.iana.org/rdap/dns.json -o "$BOOTSTRAP_CACHE.tmp" 2>/dev/null \
    && mv "$BOOTSTRAP_CACHE.tmp" "$BOOTSTRAP_CACHE" 2>/dev/null
}

# 解析指定 TLD 的 RDAP base url(取第一个)
resolve_rdap_base() {
  local t="$1"
  [[ -f "$BOOTSTRAP_CACHE" ]] || return 1
  python3 - "$t" "$BOOTSTRAP_CACHE" <<'PY' 2>/dev/null
import json,sys
tld,path=sys.argv[1],sys.argv[2]
try:
    d=json.load(open(path))
except Exception:
    sys.exit(1)
for svc in d.get("services",[]):
    tlds,urls=svc[0],svc[1]
    if tld in tlds and urls:
        print(urls[0].rstrip("/"))
        break
PY
}

check_rdap() {
  local url="$1"
  local code
  code="$(curl -sL -m "$TIMEOUT" -o /dev/null -w "%{http_code}" "$url" 2>/dev/null)"
  case "$code" in
    200) echo "REGISTERED"; return 0 ;;
    404) echo "AVAILABLE";  return 0 ;;
    *)   return 1 ;;
  esac
}

check_whois() {
  local out
  out="$(whois "$domain" 2>/dev/null)"
  [[ -z "$out" ]] && return 1
  # 未注册的常见提示(放前面优先判断)
  if echo "$out" | grep -qiE "No match for|NOT FOUND|No Data Found|Domain not found|is available for|No entries found|Status: *free|Status: *AVAILABLE"; then
    echo "AVAILABLE"; return 0
  fi
  # 已注册的常见字段
  if echo "$out" | grep -qiE "Registry Domain ID|Creation Date|Registrar:|Registrar WHOIS|Domain Status|Name Server|Updated Date"; then
    echo "REGISTERED"; return 0
  fi
  return 1
}

result=""

# --- 1. RDAP(动态端点) ---
fetch_bootstrap
base="$(resolve_rdap_base "$tld")"
if [[ -n "$base" ]]; then
  if result="$(check_rdap "${base}/domain/${domain}")"; then
    printf "%s\t%s\t%s\n" "$domain" "$result" "rdap"
    exit 0
  fi
fi

# --- 2. whois 兜底 ---
if result="$(check_whois)"; then
  printf "%s\t%s\t%s\n" "$domain" "$result" "whois"
  exit 0
fi

printf "%s\t%s\t%s\n" "$domain" "UNKNOWN" "none"
exit 0
