// GetMapKeys 返回 map 的所有 key
func GetMapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}

// GetMapValues 返回 map 的所有 value
func GetMapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}
