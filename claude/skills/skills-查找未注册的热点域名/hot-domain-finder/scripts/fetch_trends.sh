#!/usr/bin/env bash
# fetch_trends.sh - 抓取 Google Trends 每日热搜(自带流量,便于排序)
# 用法:
#   ./fetch_trends.sh [geo] [top_n]
#   geo   地区码,默认 US。常用: US(英文热词) / HK(部分中文) / GB / JP / IN
#   top_n 取前 N 条,默认 20
# 输出(TSV,含表头): rank  traffic  keyword  context
#   traffic 已归一化为数字(如 "200+" -> 200),按降序排列。
#   context 取首条关联新闻标题,帮助理解热词语境(便于翻译/取关键词)。

set -uo pipefail

geo="${1:-US}"
top_n="${2:-20}"
TIMEOUT=12

url="https://trends.google.com/trending/rss?geo=${geo}"
xml="$(curl -s -m "$TIMEOUT" -A "Mozilla/5.0" "$url" 2>/dev/null)"

if [[ -z "$xml" ]] || echo "$xml" | grep -qi "Error 400"; then
  echo "ERROR: 无法获取 geo=${geo} 的热搜(该地区可能不支持,试试 US/HK/GB)" >&2
  exit 1
fi

tmp_xml="$(mktemp)"
printf '%s' "$xml" > "$tmp_xml"
trap 'rm -f "$tmp_xml"' EXIT

python3 - "$top_n" "$tmp_xml" <<'PY'
import sys, re, html
from xml.etree import ElementTree as ET

top_n = int(sys.argv[1])
data = open(sys.argv[2], encoding="utf-8").read()

# 处理带命名空间的 ht: 标签
data = data.replace('xmlns:ht="https://trends.google.com/trending/rss"', '')
data = re.sub(r'<ht:', '<ht_', data)
data = re.sub(r'</ht:', '</ht_', data)

def norm_traffic(s):
    if not s: return 0
    s = s.replace('+','').replace(',','').strip()
    m = re.match(r'([\d.]+)\s*([KkMm]?)', s)
    if not m: return 0
    n = float(m.group(1)); unit = m.group(2).lower()
    if unit == 'k': n *= 1000
    elif unit == 'm': n *= 1_000_000
    return int(n)

try:
    root = ET.fromstring(data)
except Exception as e:
    sys.stderr.write(f"XML parse error: {e}\n"); sys.exit(1)

rows = []
for item in root.iter('item'):
    title = (item.findtext('title') or '').strip()
    traffic = norm_traffic(item.findtext('ht_approx_traffic') or '')
    ctx = ''
    ni = item.find('ht_news_item')
    if ni is not None:
        ctx = (ni.findtext('ht_news_item_title') or '').strip()
    if title:
        rows.append((traffic, title, ctx))

rows.sort(key=lambda r: r[0], reverse=True)
print("rank\ttraffic\tkeyword\tcontext")
for i, (tr, kw, ctx) in enumerate(rows[:top_n], 1):
    ctx = html.unescape(ctx).replace('\t',' ').replace('\n',' ')
    kw = html.unescape(kw)
    print(f"{i}\t{tr}\t{kw}\t{ctx}")
PY
