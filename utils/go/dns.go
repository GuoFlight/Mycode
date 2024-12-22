package utils

import (
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// ReverseDomain 反转域名
func ReverseDomain(domain string) string {
	domainSlice := strings.Split(domain, ".")

	for i, j := 0, len(domainSlice)-1; i < j; i, j = i+1, j-1 {
		domainSlice[i], domainSlice[j] = domainSlice[j], domainSlice[i]
	}

	domainReverse := strings.Join(domainSlice, ".")
	return domainReverse
}

// ServernameToIP 将域名转化为ip
func ServernameToIP(servernames []string) (map[string]string, error) {
	ret := make(map[string]string)
	for _, hostname := range servernames {
		cmd := exec.Command("host", hostname)
		output, err := cmd.Output() // 执行程序，返回标准输出[]byte
		if err != nil {
			return nil, errors.New(fmt.Sprintf("获取%s ip失败,请检查域名列表", hostname))
		}

		r := regexp.MustCompile("[\\d]+.[\\d]+.[\\d]+.[\\d]+")
		RetRegexp := r.FindString(string(output))
		if RetRegexp == "" {
			return nil, errors.New(fmt.Sprintf("获取%s ip失败,请检查域名列表", hostname))
		}
		ret[hostname] = RetRegexp
	}
	return ret, nil
}

// IsContainByWildcardDn 判断是否属于某个泛域名
func IsContainByWildcardDn(dn, wildcardDn string) bool {
	if strings.HasPrefix(wildcardDn, "*.") {
		wildcardDn = wildcardDn[2:]
	} else {
		return false
	}
	if dn == wildcardDn {
		return false
	}
	return dns.IsSubDomain(wildcardDn, dn)
}

// GetSubZoneFromDn 通过域名得到zone
func GetSubZoneFromDn(dn string) (string, error) {
	indexDot := strings.Index(dn, ".")
	if indexDot == -1 {
		return "", errors.New("获取子域失败: " + dn)
	}
	zone := dns.Fqdn(dn[indexDot+1:])
	if zone == "." {
		return "", errors.New("获取子域失败: " + dn)
	}
	return zone, nil
}

// CheckDnA 检查域名是否有A记录
func CheckDnA(dn, dnsServer string) (bool, error) {
	dn = dns.Fqdn(dn)

	c := dns.Client{}
	c.ReadTimeout = 3000 * time.Millisecond
	m := dns.Msg{}
	m.SetQuestion(dn, dns.TypeA)
	r, _, err := c.Exchange(&m, dnsServer)
	if err != nil {
		return false, err
	}
	// 从返回结果中得到A记录
	for _, ans := range r.Answer {
		_, ok := ans.(*dns.A)
		// fmt.Println(recordA) // debug
		if ok {
			return true, nil
		}
	}
	return false, nil
}

func isSameBase(zone string, domains []string) bool {
	for _, domain := range domains {
		if zone == dns.Fqdn(domain) {
			continue
		}
		subZone, err := GetSubZoneFromDn(domain)
		if err != nil {
			return false
		}
		if subZone != zone {
			return false
		}
	}
	return true
}
func IsSameBase(domains []string) bool {
	if len(domains) == 0 {
		return true
	}

	baseDomain, err := GetSubZoneFromDn(domains[0])
	if err == nil {
		if isSameBase(baseDomain, domains) {
			return true
		}
	}
	baseDomain = dns.Fqdn(domains[0])
	return isSameBase(baseDomain, domains)
}
