package utils

import (
	"strconv"
	"strings"
)

// RemoveRepeat 字符串去重
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

// StringToUnicode 函数作用：String类型的"\u91d1\u989d"转换成Unicode的"金额"
func StringToUnicode(strOrigin string) (string, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(strOrigin), `\\u`, `\u`, -1))
	if err != nil {
		return "", err
	}
	return str, nil
}
