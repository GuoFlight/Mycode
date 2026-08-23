#!/usr/bin/env bash
# fetch_trends_cn.sh - 抓取中国大陆实时热搜(微博主源 + 百度备用)
# 用法:
#   ./fetch_trends_cn.sh [source] [top_n]
#   source  数据源: weibo(默认,带真实热度值) / baidu(按排名) / all(两源合并去重)
#   top_n   取前 N 条,默认 20
# 输出(TSV,含表头): rank  traffic  keyword  context
#   weibo 源:traffic 为微博热搜热度值(num,越大越热),按降序
#   baidu 源:无公开热度值,traffic 用 (52-排名)*20 归一化,便于统一排序
#   context: 数据源标记 + 榜单标签(如 热/新/沸),辅助判断话题性质
# 说明:热词为中文,需在 skill 内借助语境翻译成英文关键词再生成域名。

set -uo pipefail

source_type="${1:-weibo}"
top_n="${2:-20}"
TIMEOUT=12
UA="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0 Safari/537.36"

W_FILE="$(mktemp)"; B_FILE="$(mktemp)"
: > "$W_FILE"; : > "$B_FILE"
trap 'rm -f "$W_FILE" "$B_FILE"' EXIT

fetch_weibo() {
  local cj; cj="$(mktemp)"
  # 先访问首页拿匿名 cookie,否则接口返回 Forbidden
  curl -s -m "$TIMEOUT" -A "$UA" -c "$cj" "https://weibo.com/" -o /dev/null 2>/dev/null
  curl -s -m "$TIMEOUT" -A "$UA" -b "$cj" -H "Referer: https://weibo.com/" \
    "https://weibo.com/ajax/statuses/hot_band" 2>/dev/null > "$W_FILE"
  rm -f "$cj"
}

fetch_baidu() {
  curl -s -m "$TIMEOUT" -A "$UA" \
    "https://top.baidu.com/api/board?platform=wise&tab=realtime" 2>/dev/null > "$B_FILE"
}

case "$source_type" in
  weibo) fetch_weibo ;;
  baidu) fetch_baidu ;;
  all)   fetch_weibo; fetch_baidu ;;
  *) echo "ERROR: 未知数据源 '$source_type'(可选 weibo/baidu/all)" >&2; exit 1 ;;
esac

python3 - "$top_n" "$source_type" "$W_FILE" "$B_FILE" <<'PY'
import sys, json
top_n = int(sys.argv[1])
source_type = sys.argv[2]
weibo_raw = open(sys.argv[3], encoding="utf-8").read()
baidu_raw = open(sys.argv[4], encoding="utf-8").read()

rows = []  # (traffic:int, keyword:str, context:str)

def parse_weibo(raw):
    out = []
    try:
        d = json.loads(raw)
        bl = d["data"]["band_list"]
    except Exception:
        return out
    for b in bl:
        kw = (b.get("word") or b.get("note") or "").strip()
        if not kw:
            continue
        num = b.get("num") or 0
        try:
            num = int(num)
        except Exception:
            num = 0
        label = (b.get("label_name") or "").strip()
        ctx = "weibo" + (f"|{label}" if label else "")
        out.append((num, kw, ctx))
    return out

def parse_baidu(raw):
    out = []
    try:
        d = json.loads(raw)
        lst = d["data"]["cards"][0]["content"][0]["content"]
    except Exception:
        return out
    for c in lst:
        kw = (c.get("word") or "").strip()
        if not kw:
            continue
        idx = c.get("index")
        try:
            idx = int(idx)
        except Exception:
            idx = 50
        # 百度无公开热度值,用排名归一化(第1名≈1020,便于与微博量级共存排序)
        traffic = max(0, (52 - idx)) * 20
        # hotTag: 1=热 2=新 3=沸 4=商(数字标签映射为中文,便于判断话题性质)
        tagmap = {"1": "热", "2": "新", "3": "沸", "4": "商"}
        raw_tag = str(c.get("hotTag") or "").strip()
        tag = tagmap.get(raw_tag, "") or (c.get("newHotName") or "").strip()
        ctx = "baidu" + (f"|{tag}" if tag else "")
        out.append((traffic, kw, ctx))
    return out

if source_type in ("weibo", "all"):
    rows += parse_weibo(weibo_raw)
if source_type in ("baidu", "all"):
    rows += parse_baidu(baidu_raw)

if not rows:
    sys.stderr.write("ERROR: 未获取到任何中文热搜(可能被限流,稍后重试或换 source)\n")
    sys.exit(1)

# all 模式:同名热词去重,保留热度较高者
if source_type == "all":
    best = {}
    for tr, kw, ctx in rows:
        if kw not in best or tr > best[kw][0]:
            best[kw] = (tr, kw, ctx)
    rows = list(best.values())

rows.sort(key=lambda r: r[0], reverse=True)
print("rank\ttraffic\tkeyword\tcontext")
for i, (tr, kw, ctx) in enumerate(rows[:top_n], 1):
    kw = kw.replace("\t", " ").replace("\n", " ")
    ctx = ctx.replace("\t", " ").replace("\n", " ")
    print(f"{i}\t{tr}\t{kw}\t{ctx}")
PY
