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
	fi, err := f.Stat()
	if err != nil || fi.Size() == 0 {
		if err == nil {
			err = errors.New("file size is nil")
		}
		return
	}

	err = GetXzData(f, fi.Size(), callBack)

	return
}

func GetXzData(fileReader io.Reader, fileSize int64, callBack ZipCallBack) (err error) {
	if callBack == nil || fileReader == nil || fileSize == 0 {
		err = errors.New("callBack is nil or io.Reader is nil or fileSize is 0")
		return
	}

	xzReader, err := xz.NewReader(fileReader)
	if err != nil {
		return
	}

	//处理压缩包中的文件
	callBack(xzReader, "", fileSize)
	return
}
