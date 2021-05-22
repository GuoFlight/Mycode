package tools

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

//函数作用：切片反转
func ReverseSlice(s []int) []int {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

//函数作用：域名反转
func ReverseDomain(domain string) string {
	domainSlice := strings.Split(domain,".")

	for i, j := 0, len(domainSlice)-1; i < j; i, j = i+1, j-1 {
		domainSlice[i], domainSlice[j] = domainSlice[j], domainSlice[i]
	}

	domainReverse := strings.Join(domainSlice,".")
	return domainReverse
}

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

