package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// DelFile 删除文件
// 文件不存在时不报错
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

// DelSubFiles 删除目录下的所有文件
func DelSubFiles(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if dir == path {
			return nil
		}
		err = DelFile(path)
		if err != nil {
			return err
		}
		return nil
	})
}
func WalkFilesBySuffix(dir, suffix string) ([]os.FileInfo, error) {
	var ret []os.FileInfo
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if dir == path {
			return nil
		}
		if strings.HasSuffix(info.Name(), suffix) {
			ret = append(ret, info)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}
// WalkFiles 得到目录下的所有文件(不包括目录)
func WalkFiles(dir string) ([]os.FileInfo, error) {
	var ret []os.FileInfo
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if dir == path {
			return nil
		}
		if !info.IsDir() {
			ret = append(ret, info)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}
