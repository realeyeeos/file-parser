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
	"strings"
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

	err = GetBz2Data(f, fileName, callBack)
	return
}

//获取文件数据
func GetBz2Data(fileReader io.Reader, fileName string, callBack ZipCallBack) (err error) {
	if callBack == nil || fileReader == nil {
		err = errors.New("callBack is nil or io.Reader is nil ")
		return
	}

	//获取解析数据流
	bz := bzip2.NewReader(fileReader)

	//处理压缩包里的文件
	num := strings.LastIndex(fileName, ".")
	fileTypeName := fileName[0:num]
	callBack(bz, fileTypeName)

	return
}
