/*
Date：2023.03.02
Author：scl
Description：总体文档
*/

package fileparser

import "io"

// 处理文件函数(数据，位置)
type CallBackDataFunc func(string, string) bool

type ZipCallBack func(io.Reader, string)
