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
func GetBz2DataFile(fileName string, callBack ZipCallBack) (err error) {
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

	err = GetBz2Data(f, fi.Size(), callBack)
	return
}

//获取文件数据
func GetBz2Data(fileReader io.Reader, fileSize int64, callBack ZipCallBack) (err error) {
	if callBack == nil || fileReader == nil || fileSize == 0 {
		err = errors.New("callBack is nil or io.Reader is nil or fileSize is 0")
		return
	}

	bz := bzip2.NewReader(fileReader)

	//处理压缩包里的文件
	callBack(bz, "", fileSize)

	return
}
