package fileparser

/*
Date：2023.03.02
Author：scl
Description：解析ole文件
*/

import (
	"errors"
	"io"
	"os"

	"github.com/realeyeeos/file-parser/ole"
)

//获取文件数据
func GetOffice97DataFile(fileName string, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer f.Close()

	err = GetOffice97Data(f, 0, callBack)
	return
}

//获取文件数据
//fileSize 无用，为了统一格式
func GetOffice97Data(fileReadSeeker io.ReadSeeker, fileSize int64, callBack CallBackDataFunc) (err error) {
	if callBack == nil || fileReadSeeker == nil {
		err = errors.New("callBack is nil or io.ReadSeeker is nil")
		return
	}

	//获取oke文件句柄
	var oleInfo ole.OleInfo
	err = oleInfo.GetHandle(fileReadSeeker)
	if err != nil {
		return
	}

	//处理文件数据
	err = dealOffice97File(fileReadSeeker, &oleInfo, callBack)
	return
}

//处理office97文件
func dealOffice97File(fileReadSeeker io.ReadSeeker, ole *ole.OleInfo, callBack CallBackDataFunc) (err error) {
	err = ole.Read(fileReadSeeker)
	if err != nil {
		return
	}

	//获取数据
	err = ole.GetObjectData(func(str, position string) bool {
		if len(str) == 0 {
			return true
		}

		return callBack(str, position)
	})

	return
}
