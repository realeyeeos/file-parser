package fileparser

/*
Date：2023.03.02
Author：scl
Description：解析ole文件
*/

import (
	"errors"
	"fileparser/ole"
	"os"
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

	err = GetOffice97Data(f, callBack)
	return
}

//获取文件数据
func GetOffice97Data(f *os.File, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	if f == nil {
		err = errors.New("os.File is nil")
		return
	}

	//获取oke文件句柄
	var oleInfo ole.OleInfo
	err = oleInfo.GetHandle(f)
	if err != nil {
		return
	}

	//处理文件数据
	err = dealOffice97File(f, &oleInfo, callBack)
	return
}

//处理office97文件
func dealOffice97File(fp *os.File, ole *ole.OleInfo, callBack CallBackDataFunc) (err error) {
	err = ole.Read(fp)
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
