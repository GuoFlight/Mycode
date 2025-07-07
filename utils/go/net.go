package utils

import (
	"fmt"
	"net"
	"sort"
	"strconv"
)

// LocalIP 获取第一个非 loopback ip
func LocalIP() (net.IP, error) {
	tables, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, t := range tables {
		addrs, err := t.Addrs()
		if err != nil {
			return nil, err
		}
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() {
				continue
			}
			if v4 := ipnet.IP.To4(); v4 != nil {
				return v4, nil
			}
		}
	}
	return nil, fmt.Errorf("cannot find local IP address")
}

// GetLocalIps 获取非环回 ip
func GetLocalIps() ([]net.IP, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var ret []net.IP
	for _, netIf := range interfaces {
		addrs, err := netIf.Addrs()
		if err != nil {
			return nil, err
		}
		for _, a := range addrs {
			ipNet, ok := a.(*net.IPNet)
			if !ok || ipNet.IP.IsLoopback() {
				continue
			}
			if v4 := ipNet.IP.To4(); v4 != nil {
				ret = append(ret, v4)
			}
		}
	}
	return ret, nil
}

// IsPortListening 函数作用：判断某个端口是否被监听
func IsPortListening(port int) bool {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return true
	}
	defer ln.Close()
	return false
}

// IpToBin 将IPv4地址转换为二进制字符串
func IpToBin(ip string) (string, error) {
	// 解析IP地址
	addr := net.ParseIP(ip)
	if addr == nil {
		return ip, fmt.Errorf("无效的IP地址")
	}

	// 确保IP地址是IPv4
	ipv4 := addr.To4()
	if ipv4 == nil {
		return ip, fmt.Errorf("不支持的IP版本，仅支持IPv4")
	}

	// 将IP地址转换为二进制形式
	binaryIP := ""
	for _, byte := range ipv4 {
		// 将每个字节转换为8位二进制形式，并拼接
		binaryByte := fmt.Sprintf("%08b", byte)
		binaryIP += binaryByte
	}

	return binaryIP, nil
}

// SortIps 给ip排序
func SortIps(ips []string) []string {
	sort.Slice(ips, func(i, j int) bool {
		ip1, _ := IpToBin(ips[i])
		ip2, _ := IpToBin(ips[j])
		return ip1 < ip2
	})
	return ips
}

// IsIpInCIDRs 判断某个ip是否属于指定网段中(满足其一即可)
func IsIpInCIDRs(ip string, cidrs []string) (bool, error) {
	ipObj := net.ParseIP(ip)
	if ipObj == nil {
		return false, fmt.Errorf("ip异常: %s", ip)
	}
	for _, cidr := range cidrs {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return false, fmt.Errorf("网段解析失败: %s", err)
		}
		if ipNet.Contains(ipObj) {
			return true, nil
		}
	}
	return false, nil
}
// GetIpsFromStr 从文本中获取所有ip
func GetIpsFromStr(text string) []string {
	ipRegex := `(((2[0-5]{2}\.)|(1\d{2}\.)|(\d{1,2}\.)){3}((1\d{2})|(2[0-5]{2})|(\d{1,2})))`
	re := regexp.MustCompile(ipRegex)
	return re.FindAllString(text, -1)
}

// CheckConnTcp 检查TCP是否能连接
func CheckConnTcp(host string, port int, timeout time.Duration) bool {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// CheckConnTcpMulti 高并发探测TCP是否能连接
func CheckConnTcpMulti(hosts []string, ports []int, timeout time.Duration, maxConcurrency int) bool {
	ctx, cancel := context.WithCancel(context.Background()) // 用于取消多余的goroutine
	defer cancel()
	results := make(chan bool) // 用于保存所有goroutine的结果
	var wg sync.WaitGroup      // 等待所有goroutine执行完成
	chanConcurrencyLimit := make(chan struct{}, maxConcurrency)

	for _, host := range hosts {
		for _, port := range ports {
			wg.Add(1)
			go func(h string, p int) {
				defer wg.Done()
				chanConcurrencyLimit <- struct{}{}
				defer func() { <-chanConcurrencyLimit }()
				select {
				case <-ctx.Done():
					return // 如果已经取消，提前退出，避免无效探测
				default:
					ok := CheckConnTcp(h, p, timeout)
					// 防止无缓冲区channel被关闭时无法写入
					select {
					case results <- ok: // 等待写入(因为无缓冲区，所以可能要等待)
					case <-ctx.Done(): // 也有可能直接退出
					}
				}
			}(host, port)
		}
	}

	// 如果执行完了，就关闭channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// 消费结果
	for result := range results {
		if !result { // 只要遇到false就立即返回
			cancel()
			return false
		}
	}
	return true
}
