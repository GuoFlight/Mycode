# keepawake-mac

让 Mac 在**充电且未锁屏**时屏幕保持常亮、不锁屏;其余情况(电池供电,或已被手动锁屏)完全交还 macOS 默认策略,不做任何干预。

纯 Shell 脚本,无需编译、无依赖。

## 工作原理

脚本由 crontab 每分钟触发一次,每次运行只做判定与「续租」:

| 状态 | 行为 |
|------|------|
| 已锁屏 | 清理本脚本启动的 caffeinate,退出 |
| 未插电源适配器 | 清理本脚本启动的 caffeinate,退出 |
| 解锁 且 已插电 | 清理旧实例,启动 `caffeinate -d 90`,退出 |

`caffeinate` 靠一个常驻进程持有电源断言来阻止休眠,进程退出断言即释放。因此每次只持有 90 秒(`-t 90`);90 > cron 的 60 秒间隔,保证相邻两次执行的覆盖窗口无缝衔接,从而持续防熄屏。pid 记录在 `$TMPDIR/keepawake-mac.pid`(cron 环境无 TMPDIR 时兜底 `/tmp`),kill 前会校验目标确实是 caffeinate 进程,避免误杀。

判定用的都是系统自带命令:

- **锁屏**:`ioreg -n Root -d1 -a` 中的 `CGSSessionScreenIsLocked`,锁屏时为 `Yes`,未锁屏时该键不出现
- **充电**:`pmset -g batt` 第一行,插电为 `Now drawing from 'AC Power'`,电池为 `'Battery Power'`
- 使用 `caffeinate -d`(阻止显示器休眠)。在 macOS 上,持有显示器断言会连带压制屏保与随之而来的锁屏(屏保与熄屏共享同一套 HID 空闲计时),因此屏幕常亮即不会锁屏,全程不修改任何系统设置

## 系统限制

`caffeinate` **无法阻止合盖(clamshell)休眠**——那是固件级行为。合盖后若要不休眠,需同时外接电源与显示器。

## 配置步骤

```
chmod +x keepawake-mac.sh
crontab -e      # 添加(路径替换为脚本的实际绝对路径):
    * * * * * cd <脚本路径> && ./keepawake-mac.sh >> "$HOME/Library/Logs/keepawake-mac.log" 2>&1
```

## 手动测试

```bash
bash keepawake-mac.sh
```

## 排查

```bash
# 看判定日志
tail -f "$HOME/Library/Logs/keepawake-mac.log"

# 确认断言生效(应出现 "display sleep prevented by caffeinate")
pmset -g | grep -i display

# 看本脚本启动的 caffeinate
pgrep -lf 'caffeinate -d'
```

## 停用

删除 crontab 行即可。残留的 caffeinate 会在 90 秒内自动退出;如需立即停止:

```bash
kill "$(cat "${TMPDIR:-/tmp}/keepawake-mac.pid")" 2>/dev/null
```
