package utils

import "os"

func IsDir(path string) bool {
	stat, err := os.Stat(path)
	return err == nil && stat.IsDir()
}
