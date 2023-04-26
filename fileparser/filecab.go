package fileparser

/*
Date：2023.04.24
Author：scl
Description：解析cab文件
*/

import (
	"errors"
	"io"
	"os"

	"github.com/secDre4mer/go-cab"
)

//获取文件数据
func GetCabDataFile(fileName string, callBack CallBackDataFunc) (err error) {
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

	err = GetCabData(f, callBack)
	return
}

//获取文件数据
func GetCabData(f *os.File, callBack CallBackDataFunc) (err error) {
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

	//打开cab文件
	cabinetFile, _ := cab.Open(f, fi.Size())
	for _, file := range cabinetFile.Files {
		reader, err := file.Open()
		if err != nil {
			continue
		}
		vfi := file.Stat()
		if vfi.Size() == 0 {
			continue
		}
		//读取文件
		data := make([]byte, vfi.Size())
		_, err = reader.Read(data)
		if err != nil && err != io.EOF {
			continue
		}
		callBack(string(data), file.Name)
	}

	return
}
