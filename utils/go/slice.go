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
