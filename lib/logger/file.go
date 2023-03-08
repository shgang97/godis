package logger

import (
	"fmt"
	"os"
)

/*
@author: shg
@since: 2023/2/24 12:12 AM
@mail: shgang97@163.com
*/

// 根据指定的目录和文件名，创建日志存储文件
func mustOpen(fileName, dir string) (*os.File, error) {
	if checkPermission(dir) {
		return nil, fmt.Errorf("permission denied dir: %s", dir)
	}

	if err := isNotExistMkDir(dir); err != nil {
		return nil, fmt.Errorf("error during make dir %s, err: %s", dir, err)
	}

	fp, err := os.OpenFile(dir+string(os.PathSeparator)+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to open file, err: %s", err)
	}

	return fp, nil
}

// 判断是否具有对目录 dir 的操作权限
func checkPermission(dir string) bool {
	_, err := os.Stat(dir)
	return os.IsPermission(err)
}

// 目录 dir 不存在则创建
func isNotExistMkDir(dir string) error {
	if checkNotExist(dir) {
		return mkDir(dir)
	}
	return nil
}

func checkNotExist(dir string) bool {
	_, err := os.Stat(dir)
	return os.IsNotExist(err)
}

func mkDir(dir string) error {
	return os.Mkdir(dir, os.ModePerm) // os.ModePerm Unix permission bits, 0o777
}
