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
	fi, err := f.Stat()
	if err != nil || fi.Size() == 0 {
		if err == nil {
			err = errors.New("file size is nil")
		}
		return
	}

	err = GetLzmaData(f, fi.Size(), callBack)
	return
}

//获取文件数据
func GetLzmaData(fileReader io.Reader, fileSize int64, callBack ZipCallBack) (err error) {
	if callBack == nil || fileReader == nil {
		err = errors.New("callBack is nil or io.Reader is nil or fileSize is 0")
		return
	}

	lzmaReader, err := lzma.NewReader(fileReader)
	if err != nil {
		return
	}

	callBack(lzmaReader, "")

	return
}
