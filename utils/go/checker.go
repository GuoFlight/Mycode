package utils

import (
	"net"
	"os"
	"os/exec"
	"unicode"
)

func IsValidPort(port int) bool {
	if port <= 0 || port >= 65536 {
		return false
	}
	return true
}

func IsValidIp(ip string) bool {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return false
	} else {
		return true
	}
}

// IsValidNetSeg 判断是否为合理的网段
func IsValidNetSeg(netSeg string) bool {
	_, _, err := net.ParseCIDR(netSeg)
	if err != nil {
		return false
	}
	return true
}

// HasChinese 函数作用：检查字符串中是否含有中文
func HasChinese(str string) bool {
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			return true
		}
	}
	return false
}

func CheckFileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	// 如果返回错误为 nil，说明文件存在
	if err == nil {
		return true, nil
	}
	// 如果错误类型为文件不存在，则返回 false
	if os.IsNotExist(err) {
		return false, nil
	}
	// 其他错误，无法确定文件是否存在，通常返回 false
	return false, err
}

// CheckHasChinese 函数作用：检查字符串中是否含有中文
func CheckHasChinese(str string) bool {
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			return true
		}
	}
	return false
}

// HasCurHostRootPermissionNoPasswd 查看某用户是否拥有免密切换到root的权限
func HasCurHostRootPermissionNoPasswd(user string) (ret bool) {
	var cmd *exec.Cmd
	cmd = exec.Command("/bin/bash", "-c", "sudo -iu "+user+" sudo -i pwd")
	_, err := cmd.Output()
	if err != nil {
		return false
	}
	return true
}
