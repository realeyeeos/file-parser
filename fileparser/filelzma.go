package fileparser

/*
Date：2023.04.24
Author：scl
Description：解析xz文件
*/

import (
	"errors"
	"os"

	"github.com/ulikunitz/xz/lzma"
)

//获取文件数据
func GetLzmaDataFile(fileName string, callBack ZipCallBack) (err error) {
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

	err = GetLzmaData(f, callBack)
	return
}

//获取文件数据
func GetLzmaData(f *os.File, callBack ZipCallBack) (err error) {
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

	lzmaReader, err := lzma.NewReader(f)
	if err != nil {
		return
	}

	callBack(lzmaReader, f.Name(), fi.Size())

	return
}
