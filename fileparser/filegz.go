package fileparser

/*
Date：2023.04.24
Author：scl
Description：解析gz文件
*/

import (
	"compress/gzip"
	"errors"
	"io"
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
	fi, err := f.Stat()
	if err != nil || fi.Size() == 0 {
		if err == nil {
			err = errors.New("file size is nil")
		}
		return
	}

	err = GetGzData(f, fi.Size(), callBack)
	return
}

func GetGzData(fileReader io.Reader, fileSize int64, callBack ZipCallBack) (err error) {
	if callBack == nil || fileReader == nil || fileSize == 0 {
		err = errors.New("callBack is nil or io.Reader is nil or fileSize is 0")
		return
	}

	gr, err := gzip.NewReader(fileReader)
	if err != nil {
		return
	}
	defer gr.Close()
	//处理压缩包里的文件
	callBack(gr, "", fileSize)
	return
}
