package fileparser

/*
Date：2023.04.24
Author：scl
Description：解析gz文件
*/

import (
	"compress/gzip"
	"errors"
	"os"
)

//获取文件数据
func GetGzDataFile(fileName string, callBack ZipCallBack) (err error) {
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

	err = GetGzData(f, callBack)
	return
}

func GetGzData(f *os.File, callBack ZipCallBack) (err error) {
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

	gr, err := gzip.NewReader(f)
	if err != nil {
		return
	}
	defer gr.Close()
	//处理压缩包里的文件
	callBack(gr, fi.Name(), fi.Size())
	return
}
