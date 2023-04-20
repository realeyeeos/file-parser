/*
Date：2023.03.07
Author：scl
Description：获取Windows下文件时间
*/
package file

import (
	"os"
	"syscall"
	"time"
)

// 获取文件创建时间、最后访问时间
func GetFileTime(fi os.FileInfo) (string, string) {
	winFileAttr := fi.Sys().(*syscall.Win32FileAttributeData)

	//文件创建时间：
	creatTime := time.Unix(winFileAttr.CreationTime.Nanoseconds()/1e9, 0)
	creatTimeStr := creatTime.Format("2006-01-02 15:04:05")
	//最后访问时间
	lastAccessTime := time.Unix(winFileAttr.LastAccessTime.Nanoseconds()/1e9, 0)
	lastAccessTimeStr := lastAccessTime.Format("2006-01-02 15:04:05")

	return creatTimeStr, lastAccessTimeStr
}
