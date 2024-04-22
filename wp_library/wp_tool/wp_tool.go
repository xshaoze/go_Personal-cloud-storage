package wp_tool

import (
	"os"
	"path/filepath"
)

func CreateDirRecursively(path string) error {
	// 检查路径是否已存在
	_, err := os.Stat(path)
	if err == nil {
		return nil // 如果路径已存在，则不做任何操作
	}

	// 如果路径不存在，则创建
	if os.IsNotExist(err) {
		// 递归创建上级目录
		parent := filepath.Dir(path)
		if err := CreateDirRecursively(parent); err != nil {
			return err
		}

		// 创建当前目录
		if err := os.Mkdir(path, 0755); err != nil {
			return err
		}

		return nil
	}

	// 其他错误情况
	return err
}
