package pkg

import "os"

func Init() error {
	// 创建目录
	if err := CreateDirs(); err != nil {
		return err
	}

	return nil
}

// CreateDirs 创建目录
func CreateDirs() error {
	dirs := []string{
		"data",
	}

	for _, dir := range dirs {
		// 使用 os.MkdirAll 创建目录，如果目录已经存在，不会返回错误
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}
