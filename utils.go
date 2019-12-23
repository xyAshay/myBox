package main

import (
	"fmt"
	"os"
)

func getFileSize(file os.FileInfo) string {
	if file.IsDir() {
		return "-"
	}
	size := file.Size()
	switch {
	case size > 1024*1024*1024:
		return fmt.Sprintf("%.2f GB", float64(size)/1024/1024/1024)
	case size > 1024*1024:
		return fmt.Sprintf("%.2f MB", float64(size)/1024/1024)
	case size > 1024:
		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	default:
		return fmt.Sprintf("%d Bytes", int(size))
	}
}
