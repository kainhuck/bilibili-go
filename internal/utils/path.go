package utils

import "os"

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true // 文件存在
	}
	if os.IsNotExist(err) {
		return false // 文件不存在
	}
	return false // 发生了其他错误，可能无法确定文件是否存在
}
