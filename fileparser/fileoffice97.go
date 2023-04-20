package fileparser

/*
Date：2023.03.02
Author：scl
Description：解析ole文件
*/

import (
	"collector/ole"
	"errors"
	"os"
)

//打开文件
func GetOffice97Data(fileName string, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer f.Close()

	var oleInfo ole.OleInfo
	err = oleInfo.GetHandle(f)
	if err != nil {
		return
	}

	err = DealOffice97File(f, &oleInfo, callBack)
	return
}

//处理office97文件
func DealOffice97File(fp *os.File, ole *ole.OleInfo, callBack CallBackDataFunc) (err error) {
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
