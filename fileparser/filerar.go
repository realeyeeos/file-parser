package fileparser

/*
Date：2023.04.24
Author：scl
Description：解析rar文件
*/

import (
	"errors"
	"io"

	"github.com/mholt/archiver"
)

func GetRarDataFile(fileName string, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	//打开文件

	rar := archiver.Rar{}
	defer rar.Close()
	rar.Walk(fileName, func(f archiver.File) error {
		if f.Size() == 0 {
			return errors.New("file size is 0")
		}
		data := make([]byte, f.Size())

		//读取文件数据
		_, err = f.Read(data)
		if err != nil && err != io.EOF {
			return err
		}

		callBack(string(data), f.Name())
		return nil
	})

	return
}
