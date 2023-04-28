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
	"strings"

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

	err = GetXzData(f, fileName, callBack)

	return
}

func GetXzData(fileReader io.Reader, fileName string, callBack ZipCallBack) (err error) {
	if callBack == nil || fileReader == nil {
		err = errors.New("callBack is nil or io.Reader is nil")
		return
	}

	xzReader, err := xz.NewReader(fileReader)
	if err != nil {
		return
	}

	//处理压缩包中的文件
	num := strings.LastIndex(fileName, ".")
	fileTypeName := fileName[0:num]

	callBack(xzReader, fileTypeName)
	return
}
