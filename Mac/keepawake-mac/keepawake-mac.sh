#!/bin/bash
#
############################################################################
# 作者：京城郭少
# 版本：v0.1
############################################################################
#
# 让 Mac 在「充电状态且未锁屏」时屏幕保持常亮、不锁屏;其余情况(电池供电,
# 或已被手动锁屏)完全交还 macOS 默认策略,不做任何干预。
#
# 设计:本脚本由 crontab 每分钟触发一次,每次运行只做判定与「续租」:
# 满足条件时启动 `caffeinate -d -t 90` 持有 90 秒显示器断言;90 > 60 保证
# 相邻两次 cron 执行(间隔 60s)的覆盖窗口无缝衔接,从而持续防熄屏。
# caffeinate -d 阻止显示器休眠,在 macOS 上会连带压制屏保与随之而来的锁屏
#(屏保与熄屏共享同一套 HID 空闲计时),因此屏幕常亮即不会锁屏,全程不改
# 任何系统设置。不满足条件时清理掉由本脚本启动的 caffeinate 并退出。

# cron 环境 PATH 极简,显式补齐以确保 pmset/ioreg/caffeinate 可用。
export PATH="/usr/bin:/bin:/usr/sbin:/sbin:$PATH"

# PID 是进程级临时状态,放用户临时目录。cron 环境可能没有 TMPDIR,兜底用 /tmp。
PIDFILE="${TMPDIR:-/tmp}/keepawake-mac.pid"
HOLD_SECONDS=90

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

# 已锁屏返回 0(true),未锁屏返回 1。
# ioreg 中锁屏时 CGSSessionScreenIsLocked = Yes,未锁屏时该键不出现。
is_locked() {
    ioreg -n Root -d1 -a 2>/dev/null | grep -q 'CGSSessionScreenIsLocked.*Yes'
}

# 已插电源适配器返回 0(true),电池供电返回 1。
is_on_ac() {
    pmset -g batt 2>/dev/null | grep -q "'AC Power'"
}

# 清理由本脚本启动的 caffeinate(校验确实是 caffeinate 进程,避免误杀),删 pidfile。
stop_caffeinate() {
    if [ -f "$PIDFILE" ]; then
        local pid
        pid=$(cat "$PIDFILE" 2>/dev/null)
        if [ -n "$pid" ] && ps -p "$pid" -o comm= 2>/dev/null | grep -qi caffeinate; then
            kill "$pid" 2>/dev/null
        fi
        rm -f "$PIDFILE"
    fi
}

if is_locked; then
    stop_caffeinate
    log "locked -> idle (released)"
elif ! is_on_ac; then
    stop_caffeinate
    log "on battery -> idle (released)"
else
    # 解锁 且 已插电:先清理旧实例避免并存,再续租。
    stop_caffeinate
    caffeinate -d -t "$HOLD_SECONDS" &
    mkdir -p "$(dirname "$PIDFILE")"
    echo "$!" > "$PIDFILE"
    log "unlocked+AC -> caffeinate renewed (${HOLD_SECONDS}s)"
fi
