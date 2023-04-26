package fileparser

/*
Date：2023.04.24
Author：scl
Description：解析xz文件
*/

import (
	"errors"
	"io"
	"os"

	"github.com/ulikunitz/xz"
)

//获取文件数据
func GetXzDataFile(fileName string, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	//打开文件
	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil || fi.Size() == 0 {
		if err == nil {
			err = errors.New("file size is nil")
		}
		return
	}

	xzReader, err := xz.NewReader(f)
	if err != nil {
		return
	}
	//读取数据
	data := make([]byte, fi.Size())
	_, err = xzReader.Read(data)
	if err != nil && err != io.EOF {
		return
	}
	callBack(string(data), fileName)
	return
}
