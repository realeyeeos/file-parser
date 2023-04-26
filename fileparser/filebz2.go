package fileparser

/*
Date：2023.04.24
Author：scl
Description：解析bz2文件
*/

import (
	"compress/bzip2"
	"errors"
	"io"
	"os"
)

//获取文件数据
func GetBz2DataFile(fileName string, callBack CallBackDataFunc) (err error) {
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

	err = GetBz2Data(f, callBack)
	return
}

//获取文件数据
func GetBz2Data(f *os.File, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	if f == nil {
		err = errors.New("os.File is nil")
		return
	}

	fi, err := f.Stat()
	if err != nil || fi.Size() == 0 {
		if err == nil {
			err = errors.New("file size is nil")
		}
		return
	}

	bz := bzip2.NewReader(f)
	//读取数据
	data := make([]byte, fi.Size())
	_, err = bz.Read(data)
	if err != nil && err != io.EOF {
		return
	}
	callBack(string(data), "")

	return
}
