package utils

import (
	"crypto/rand"
)

// randStrCapital 生成指定位数的随机大写英文字母串
func randStrCapital(strSize int) (string, error) {
	// 定义字母表
	dictionary := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// 生成一堆随机数
	var bytes = make([]byte, strSize)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	// 将随机数变成字母
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes), nil
}

func ConvByPage[T any](allData []T, page, size int) ([]T, int64) {
	// 计算偏移量
	offset := (page - 1) * size

	// 计算总量
	total := len(allData)

	// 数据边界检查
	if offset >= len(allData) {
		return nil, int64(total) // 超出范围，返回空
	}

	// 确定当前页的数据范围
	end := offset + size
	if end > len(allData) {
		end = len(allData)
	}

	// 返回分页后的数据
	return allData[offset:end], int64(total)
}
