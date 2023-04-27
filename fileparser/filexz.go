package fileparser

/*
Date：2023.04.24
Author：scl
Description：解析xz文件
*/

import (
	"errors"
	"os"

	"github.com/ulikunitz/xz"
)

//获取文件数据
func GetXzDataFile(fileName string, callBack ZipCallBack) (err error) {
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

	err = GetXzData(f, callBack)

	return
}

func GetXzData(f *os.File, callBack ZipCallBack) (err error) {
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
	xzReader, err := xz.NewReader(f)
	if err != nil {
		return
	}

	//处理压缩包中的文件
	callBack(xzReader, fi.Name(), fi.Size())
	return
}
