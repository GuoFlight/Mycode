package utils

import "sync"

// SafeSlice 线程安全的Slice结构体
type SafeSlice[T any] struct {
	rwLock sync.RWMutex
	items  []T
}

// NewSafeSlice 创建一个新的线程安全Slice
func NewSafeSlice[T any]() *SafeSlice[T] {
	return &SafeSlice[T]{
		items: make([]T, 0),
	}
}

// Len 获取Slice长度
func (s *SafeSlice[T]) Len() int {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	return len(s.items)
}

// Append 向Slice追加元素
func (s *SafeSlice[T]) Append(items ...T) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	s.items = append(s.items, items...)
}

// Clear 清空Slice
func (s *SafeSlice[T]) Clear() {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	s.items = make([]T, 0)
}

// Items 获取Slice的副本
func (s *SafeSlice[T]) Items() []T {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	copied := make([]T, len(s.items))
	copy(copied, s.items)
	return copied
}

// Replace 全量替换Slice内容
func (s *SafeSlice[T]) Replace(newItems []T) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	s.items = make([]T, len(newItems))
	copy(s.items, newItems)
}

// GetByIndex 获取指定索引的元素
func (s *SafeSlice[T]) GetByIndex(index int) (T, bool) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	var zero T
	if index < 0 || index >= len(s.items) {
		return zero, false
	}
	return s.items[index], true
}

// DelByIndex 删除指定索引的元素
func (s *SafeSlice[T]) DelByIndex(index int) bool {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	if index < 0 || index >= len(s.items) {
		return false
	}

	s.items = append(s.items[:index], s.items[index+1:]...)
	return true
}

// UpdateByIndex 更新指定索引的元素
func (s *SafeSlice[T]) UpdateByIndex(index int, item T) bool {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	if index < 0 || index >= len(s.items) {
		return false
	}

	s.items[index] = item
	return true
}
