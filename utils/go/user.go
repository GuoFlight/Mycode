package utils

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// GetUGid 获取指定用户的uid和gid
// 说明：无法使用os/user包来获取sso认证的用户信息
// 适用平台：Linux、Mac
func GetUGid(user string) (uid, gid int, err error) {
	uid, gid = -1, -1
	cmd := exec.Command("/bin/sh", "-c", "id "+user)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return -1, -1, err
	}
	cmd.Run() // 执行命令，阻塞直到执行完成

	// 获取uid
	re := regexp.MustCompile("uid=(?P<uid>[\\d]+)")
	res := re.FindStringSubmatch(string(output))
	if res != nil {
		uid, err = strconv.Atoi(res[1])
		if err != nil {
			return -1, -1, err
		}

	}

	// 获取gid
	re = regexp.MustCompile("gid=(?P<gid>[\\d]+)")
	res = re.FindStringSubmatch(string(output))
	if res != nil {
		gid, err = strconv.Atoi(res[1])
		if err != nil {
			return -1, -1, err
		}
	}

	return uid, gid, nil
}

// GetCurUser 获取当前用户
func GetCurUser() (string, error) {
	cmd := exec.Command("whoami")
	output, err := cmd.Output() // 执行程序，返回标准输出[]byte
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}
