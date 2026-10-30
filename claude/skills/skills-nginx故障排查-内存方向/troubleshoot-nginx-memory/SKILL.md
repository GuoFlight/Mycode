---
name: troubleshoot-nginx-memory
description: 排查 nginx 机器内存占用过高的问题,尤其是因大量"未退出的 nginx worker(worker process is shutting down)"长期驻留导致的内存不释放。当用户问"nginx 机器内存高/爆了"、"内存被谁占了"、"nginx 进程好多退不掉"、"shutting down 的 worker 一直在"等场景时触发此 skill。
tags: [nginx, memory, troubleshoot, process, ssh]
---

# 关于此SKILL

- 向用户问好,如: "你好,京城郭少!", 然后告知用户你要做的事情(排查 nginx 内存问题)。
- 作者: 京城郭少
- 版本: v0.1

# nginx 内存问题排查

排查 nginx 机器内存占用过高的问题。核心思路:**先定位内存是谁占用的 → 若是大量"未退出的 nginx worker"所致,则挑长期未退出的少数几个深挖它们为什么退不掉(通常是还挂着没结束的客户端连接)→ 顺着连接查到对端进程 → 整理产出对端服务器的进程详情。**

典型根因:nginx reload 后旧 worker 进入 `worker process is shutting down` 状态,本应处理完存量连接后自然退出,但因为仍挂着长连接迟迟不退,日积月累堆积几十个旧 worker,占用大量内存。

## 步骤 0:登录 nginx 机器

用 `ssh-login` 这个 skill 登录到出问题的 nginx 机器。

- 登录命令以 `ssh-login` skill 为准。**登录失败时**(如超时、Permission denied)向用户询问正确的 SSH 端口或完整登录命令。
- 后面排查对端机器时,**对端机器的 SSH 端口可能与 nginx 机器不同**,登录对端若失败,同样向用户询问端口。

## 步骤 1:定位内存是谁占用的

```bash
ps aux --sort=-%mem | head -10
```

判读:
- 若 Top 里大量是 `nginx: worker process is shutting down`,说明内存被**未退出的旧 nginx worker**占用 → 进入步骤 2。
- 若是其他进程(某个业务进程、某个 worker 单体内存异常大等)→ 则按常规内存排查(该进程是否泄漏、是否该限额/重启),不属于本 skill 的"旧 worker 堆积"场景。

## 步骤 2:确认旧 worker 的规模,挑少数几个来深挖

```bash
# 看有哪些"正在关闭"的 worker
ps aux --sort=-%mem | grep nginx | grep "is shutting down"
ps aux | grep "[n]ginx: worker process is shutting down" | wc -l   # 统计总数

# 按存活时长排序,挑"活得最久"的几个(etimes=已存活秒数,越大越久)
ps -eo pid,etimes,%mem,rss,lstart,cmd \
  | grep "[n]ginx: worker process is shutting down" \
  | sort -k2 -nr | head -5
```

要点(**务必遵守,否则排查会非常慢**):
- **不要排查所有 worker**。几十个全查既慢又没必要。
- **只挑长时间没退出的那几个**(存活时间最久、或占内存最大的),而且**限制数量**(建议 2~3 个)。这些"钉子户"最能说明问题。

## 步骤 3:看这些 worker 退不掉是因为还挂着谁的连接

对挑中的每个 `<pid>`:

```bash
# 看该进程的所有 TCP 连接
sudo lsof -nP -p <pid> | grep TCP
```

**关键:只看 nginx 作为"目的 IP"的连接(客户端连进来的那些),不要排查 nginx 主动发起的上游连接。**

- nginx 是代理,它的连接分两类:
  - **客户端侧(要查的)**:`本机IP:443->客户端IP:端口`(或 `:80`)。即**本机端口是 443/80**、nginx 是目的方。**正是这些"别人连进来、还没结束"的连接拖住了旧 worker 不让它退出**,所以要顺着对端(客户端 IP)去查。
  - **上游侧(不查的)**:`本机IP:随机端口->上游IP:上游端口`。这是 nginx **主动发起**到后端的连接,不是我们这次要追的对象。
- 用下面的写法只保留客户端侧连接(把 `<本机IP末段>` 换成本机 IP,或直接用 `:443->` / `:80->` 过滤):

  ```bash
  sudo lsof -nP -p <pid> | grep TCP | grep -E ':(443|80)->'
  ```

- 从结果里挑一条 `ESTABLISHED` 的客户端连接,记下**对端 IP** 和**对端源端口**(箭头右边的 `客户端IP:端口`),进入步骤 4。
- 连接是动态的:挑中的 worker 若此刻查不到客户端连接(已关闭),换下一个"钉子户",或重新 `lsof` 取最新连接。

## 步骤 4:登录对端机器,查是哪个进程建立的连接

从步骤 3 拿到**对端 IP** 和**对端源端口**后,用 `ssh-login` 登录对端(登录失败时向用户询问正确的 SSH 端口或完整登录命令):

```bash
# 在对端机器上,用"对端源端口" + "nginx 机器 IP" 定位这条连接的归属进程
sudo netstat -nalp | grep <对端源端口>
sudo netstat -nalp | grep <对端源端口> | grep <nginx机器IP>   # 进一步锁定这条连接
# 或:
sudo ss -tanp | grep <nginx机器IP>
```

正常情况下,从输出的 `PID/Program name` 列即可知道是对端的哪个进程连到了 nginx。

> ⚠️ **对端查不到本地进程的两种常见情况(真实遇到过,务必会识别):**
>
> 1. **对端是 NAT / LB / 转发网关**:对端 `netstat`/`ss` 里**找不到**这条连接对应的本地进程(甚至 `grep <nginx IP>` 为空)。排查佐证:
>    ```bash
>    cat /proc/sys/net/ipv4/ip_forward   # 为 1 说明开了转发
>    ss -s                               # estab 不为 0,但按 nginx IP 却 grep 不到
>    ```
>    说明连接是**经这台设备转发/SNAT**过来的,真正的发起方在它上游,需要继续往上游追(或结合该设备的转发/NAT 规则)。
>
> 2. **对端端口是 K8s NodePort**(常见如 `3xxxx`):节点上 `netstat`/`ss` 往往只看到该端口的 `LISTEN`(属于 `kube-proxy`),看不到 ESTABLISHED——真实连接被 iptables/IPVS DNAT 进了后端 Pod。需结合 `kubectl get svc/endpoints` 由 NodePort 反查 Service 和后端 Pod,或进 Pod 内再查。
>
> 两种情况都要把"对端是转发设备 / NodePort、真实后端在何处"作为结论的一部分记录下来,而不是停在"查不到"。

## 步骤 5:整理进程的详细信息

对需要留档的进程(对端定位到的进程,或 nginx 机器上的旧 worker),整理 `/proc/<pid>` 下的详情:

```bash
sudo ls -l /proc/<pid>/cwd      # 工作目录
sudo ls -l /proc/<pid>/exe      # 可执行文件路径(确认是哪套二进制)
ps -o pid,lstart,etimes,%mem,rss,stat -p <pid>   # 启动时间 / 已存活时长 / 内存 / 状态
sudo cat /proc/<pid>/cmdline | tr '\000' ' '; echo   # 完整启动命令行
sudo ls /proc/<pid>/fd | wc -l                        # 打开的 fd(含 socket)数量
```

关注:启动时间(存活多久)、exe 指向的二进制、cwd、打开 fd 数量(旧 worker 若还握着大量 socket,说明客户端连接确实没清干净)。

## 步骤 6:产出结论

最终**整理产出"这些对端服务器的进程的详细信息"**,把排查链路串成一条可读结论,例如:

| 旧 worker PID | 存活时长 | 占用内存(RSS) | 挂着的客户端连接 | 对端机器 | 对端进程(PID/名称) | 对端进程详情(exe/cwd/启动时间) |
|---|---|---|---|---|---|---|
| 35xxxxx | ~30h | ~1.5G | `:443-><对端IP>:<端口>` | <对端IP> | <进程名/PID>(如为转发设备则注明) | <exe>/<cwd>/<启动时间> |

结论要说明:**内存被 N 个未退出的旧 worker 占用 → 它们退不掉是因为仍挂着来自 <对端> 的客户端长连接 → 对端持有连接的进程是 <进程>(或:对端是转发设备/NodePort,真实发起方在其上游/某 Pod)**,把"内存高"归因到具体连接/进程链路上。

# 多层 SSH / 命令要点

- **变量被本地 shell 吞掉的陷阱**:`sudo -iu op ssh ... '... /proc/$pid/...'` 里的 `$pid` 会被**本地** shell 先解析成空,导致路径变成 `/proc//...` 报错。规避:
  - 远程命令里**直接写死 PID 数字**(推荐,排查时 PID 本就是确定的);
  - 或对 `$` 转义(`\$pid`)并确保远程侧有定义;
  - 同理 `awk '{...$1...}'` 的 `$` 也会被吞,必要时转义或把数据抓到本地再处理。
- `ps -p` 必须紧跟 PID(`ps -o ... -p <pid>`),PID 为空会报 "list of process IDs must follow -p"——这通常就是上面变量被吞的征兆。
- `/proc/<pid>`、`lsof`、`netstat -nalp` 的进程/文件信息多数需要 `sudo` 才能看全(否则 PID/Program、链接目标会缺失)。

# 注意事项

- **只读排查**:全程只做读操作(ps/lsof/netstat/看 /proc)。**不要**随意 `kill` 旧 worker 或 reload/重启 nginx;确需清理(如 kill 残留 worker、调整 keepalive/优雅退出超时)时**先向用户确认**。
- **只查客户端侧连接**:nginx 是目的 IP(本机 `:443`/`:80`)的连接才是拖住旧 worker 的元凶;nginx 主动发起到上游的连接不在本次排查范围。
- **限量深挖**:不要遍历全部旧 worker;挑存活最久/占内存最大的 2~3 个即可,否则排查极慢。
- **端口要确认**:登录(nginx 机器或对端)失败时,向用户询问正确的 SSH 端口或完整登录命令。
- **连接是动态的**:ESTABLISHED 连接可能在排查过程中关闭,查不到就换一个"钉子户",或重新 `lsof` 取最新连接。
- **清理**:若抓了临时文件到本地,用完删掉。

# skill自我进化

- 当 skill 出现问题,通过探索、用户引导等方式找到正确方式后,询问用户并给出完善建议,优化此 skill,同时更新版本号。

# 资源

- 依赖 `ssh-login` skill 完成登录。
- 相关:`troubleshoot-nginx`(排查 nginx 代理的访问慢/超时/报错)。
