package main

import (
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"sort"
	"sync"
	"time"
)

///////////////////////////////////////////////////////////////////////////
// 代码作用：解析域名列表，得到平均延迟、各分位延迟
// 版本：v0.1
// 最近一次更新时间：2023年5月22日
///////////////////////////////////////////////////////////////////////////

type DnsProbe struct {
	DomainList  []string
	Server      string
	Concurrency int
}

// 检查域名A记录是否正常解析
func checkDomainA(domain string, dnsServer string) (error, time.Duration) {
	c := dns.Client{}
	c.ReadTimeout = 3000 * time.Millisecond
	m := dns.Msg{}
	m.SetQuestion(domain, dns.TypeA)
	r, rtt, err := c.Exchange(&m, dnsServer)
	if err != nil {
		return err, rtt
	}
	// 从返回结果中得到A记录
	for _, ans := range r.Answer {
		_, ok := ans.(*dns.A)
		if ok {
			return nil, rtt
		}
	}
	return errors.New("查不到A记录"), 0
}

// 解析所有域名
func (d DnsProbe) CheckDomainList() {
	var errCount = 0
	var latencySum float64 = 0
	var latencys []float64

	var wg sync.WaitGroup
	wg.Add(len(d.DomainList))
	var lockErr sync.Mutex
	var lockLatencySum sync.Mutex
	countGoroutine := make(chan int, d.Concurrency)

	for _, domain := range d.DomainList {
		countGoroutine <- 1
		go func(domain string) {
			defer wg.Done()
			// 主要逻辑
			err, rtt := checkDomainA(domain, d.Server)
			if err != nil {
				lockErr.Lock()
				defer lockErr.Unlock()
				fmt.Printf("[debug] err:%s\n", domain+" "+err.Error())
				errCount++
			} else {
				lockLatencySum.Lock()
				defer lockLatencySum.Unlock()
				latency := float64(rtt.Microseconds()) / float64(1000)
				latencySum += latency
				latencys = append(latencys, latency)
			}
			<-countGoroutine
		}(domain)
	}
	wg.Wait()

	sum := len(d.DomainList)
	countSuccess := sum - errCount
	fmt.Printf("域名总数：%d\n", len(d.DomainList))
	if errCount > 0 {
		fmt.Printf("成功数：%d\n", countSuccess)
		fmt.Printf("失败数：%d\n", errCount)
	}
	if countSuccess == 0 {
		return
	}

	fmt.Printf("平均延迟：%.4f\n", latencySum/float64(countSuccess))

	//计算分位延迟
	sort.Float64s(latencys)
	index99 := int(0.99 * float64(len(latencys)-1))
	index95 := int(0.95 * float64(len(latencys)-1))
	index90 := int(0.90 * float64(len(latencys)-1))

	var tempSum float64
	i := 0
	// 90分位
	for ; i <= index90; i++ {
		tempSum += latencys[i]
	}
	fmt.Printf("90分位平均/最大延迟：%.2f  %.2f\n", tempSum/float64(index90+1), latencys[index90])

	// 95分位
	for ; i <= index95; i++ {
		tempSum += latencys[i]
	}
	fmt.Printf("95分位平均/最大延迟：%.2f  %.2f\n", tempSum/float64(index95+1), latencys[index95])

	//99 分位
	for ; i <= index99; i++ {
		tempSum += latencys[i]
	}
	fmt.Printf("99分位平均/最大延迟：%.2f  %.2f\n", tempSum/float64(index99+1), latencys[index99])
}

func main() {
	domainList := []string{"baidu.com.", "taobao.com."}
	probe := DnsProbe{
		DomainList:  domainList,
		Server:      "10.191.112.3:53",
		Concurrency: 5,
	}
	probe.CheckDomainList()
}
