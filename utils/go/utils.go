package tools

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

/////////////////////////////////////////////// 切片 ///////////////////////////////////////////////

// 函数作用：判断字符串是否在Slice中
func IsInSliceStr(str string, strs []string) bool {
	for _, item := range strs {
		if str == item {
			return true
		}
	}
	return false
}

//函数作用：切片反转
func ReverseSlice(s []int) []int {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

/////////////////////////////////////////////// DNS ///////////////////////////////////////////////

//函数作用：域名反转
func ReverseDomain(domain string) string {
	domainSlice := strings.Split(domain,".")

	for i, j := 0, len(domainSlice)-1; i < j; i, j = i+1, j-1 {
		domainSlice[i], domainSlice[j] = domainSlice[j], domainSlice[i]
	}

	domainReverse := strings.Join(domainSlice,".")
	return domainReverse
}

//将域名转化为ip
func ServernameToIP(servernames []string) (map[string]string, error) {
	ret := make(map[string]string)
	for _,hostname := range servernames{
		cmd := exec.Command("host", hostname)
		output, err :=cmd.Output()			//执行程序，返回标准输出[]byte
		if err!=nil{
			return nil,errors.New(fmt.Sprintf("获取%s ip失败,请检查域名列表", hostname))
		}

		r := regexp.MustCompile("[\\d]+.[\\d]+.[\\d]+.[\\d]+")
		RetRegexp := r.FindString(string(output))
		if RetRegexp==""{
			return nil,errors.New(fmt.Sprintf("获取%s ip失败,请检查域名列表", hostname))
		}
		ret[hostname] = RetRegexp
	}
	return ret, nil
}

// IsContainByWildcardDn 判断是否属于某个泛域名
// import "github.com/miekg/dns"
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
// GetZoneFromDn 通过域名得到zone
func GetZoneFromDn(dn string) (string, *gerror.Gerr) {
	indexDot := strings.Index(dn, ".")
	if indexDot == -1 {
		return "", gerror.NewErr("请检查域名：" + dn)
	}
	return dns.Fqdn(dn[indexDot+1:]), nil
}
// CheckDnA 检查域名是否有A记录
func CheckDnA(dn, dnsServer string) (bool,*gerror.Gerr) {
	dn = dns.Fqdn(dn)

	c := dns.Client{}
	c.ReadTimeout = 3000 * time.Millisecond
	m := dns.Msg{}
	m.SetQuestion(dn, dns.TypeA)
	r, _, err := c.Exchange(&m, dnsServer)
	if err != nil {
		return false,gerror.NewErr(err.Error())
	}
	// 从返回结果中得到A记录
	for _, ans := range r.Answer {
		_, ok := ans.(*dns.A)
		// fmt.Println(recordA) // debug
		if ok {
			return true,nil
		}
	}
	return false,nil
}

/////////////////////////////////////////////// 系统用户 ///////////////////////////////////////////////

//函数作用：获取指定用户的uid和gid
//说明：无法使用os/user包来获取sso认证的用户信息
//适用平台：Linux、Mac
func GetUGid(user string)(uid,gid int,err error){
	uid,gid=-1,-1
	cmd := exec.Command("/bin/sh", "-c", "id "+user)
	output,err:=cmd.CombinedOutput()
	if err!=nil{
		return -1,-1,err
	}
	cmd.Run()  //执行命令，阻塞直到执行完成

	//获取uid
	re:= regexp.MustCompile("uid=(?P<uid>[\\d]+)")
	res := re.FindStringSubmatch(string(output))
	if res != nil{
		uid,err = strconv.Atoi(res[1])
		if err!=nil{
			return -1,-1,err
		}

	}

	//获取gid
	re= regexp.MustCompile("gid=(?P<gid>[\\d]+)")
	res = re.FindStringSubmatch(string(output))
	if res != nil{
		gid,err = strconv.Atoi(res[1])
		if err!=nil{
			return -1,-1,err
		}
	}

	return uid,gid,nil
}

//查看某用户是否拥有免密切换到root的权限
func HasCurHostRootPermissionNoPasswd(user string)(ret bool){
	var cmd *exec.Cmd
	cmd = exec.Command("/bin/bash","-c","sudo -iu "+user+" sudo -i pwd")
	_,err := cmd.Output()
	if err!=nil{
		return false
	}
	return true
}

//获取当前用户
func GetCurUser()(string,error){
	cmd := exec.Command("whoami")
	output, err :=cmd.Output()			//执行程序，返回标准输出[]byte
	if err!=nil{
		return "", err
	}
	return strings.TrimSpace(string(output)),nil
}


/////////////////////////////////////////////// 其他 ///////////////////////////////////////////////

//字符串去重
func RemoveRepeat(str []string) []string {
	tmp := make(map[string]bool)
	for _, v := range str {
		tmp[v] = true
	}
	var ret []string
	for i, _ := range tmp {
		ret = append(ret, i)
	}
	return ret
}

// 获取第一个非 loopback ip
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

// 函数作用：判断某个端口是否被监听
func IsPortListening(port int) bool {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return true
	}
	defer ln.Close()
	return false
}



//函数作用：String类型的"\u91d1\u989d"转换成Unicode的"金额"
func StringToUnicode(strOrigin string) (string, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(strOrigin), `\\u`, `\u`, -1))
	if err != nil {
		return "", err
	}
	return str, nil
}

//函数作用：检查字符串中是否含有中文
func ContainChinese(str string)bool{
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			return true
		}
	}
	return false
}

// 函数作用：安全地递归删除文件/目录
func DelFile(path string) error {
	pathAbs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if pathAbs == "/" {
		return errors.New("无法删除/")
	}
	return os.RemoveAll(pathAbs)
}
