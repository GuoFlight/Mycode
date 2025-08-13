package utils

// IsInSlice 函数作用：判断是否在Slice中
func IsInSlice[T comparable](target T, slice []T) bool {
	for _, item := range slice {
		if target == item {
			return true
		}
	}
	return false
}

// ReverseSlice 函数作用：切片反转
func ReverseSlice[T any](s []T) []T {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// HasSameItems 判断两个 slice 是否拥有相同的元素（顺序无关）
func HasSameItems[T comparable](a, b []T) bool {
	// 创建一个 map 来统计元素出现次数
	mapA := make(map[T]struct{}, len(a))
	mapB := make(map[T]struct{}, len(b))
	for _, v := range a {
		mapA[v] = struct{}{}
	}
	for _, v := range b {
		mapB[v] = struct{}{}
	}
	if len(mapA) != len(mapB) {
		return false
	}

	// 检查第二个 slice 的元素
	for kA := range mapA {
		if _, ok := mapB[kA]; !ok {
			return false
		}
	}
	return true
}

// RemoveRepeat Slice去重
func RemoveRepeat[T comparable](s []T) []T {
	tmp := make(map[T]bool)
	for _, v := range s {
		tmp[v] = true
	}
	var ret []T
	for i, _ := range tmp {
		ret = append(ret, i)
	}
	return ret
}
