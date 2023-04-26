package fileparser

/*
Date：2023.04.25
Author：scl
Description：解析udf格式的iso文件
*/

import (
	"errors"
	"io"
	"os"

	"github.com/hooklift/iso9660"
)

//获取文件数据
func Get9660IsoDataFile(fileName string, callBack CallBackDataFunc) (err error) {
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

	err = Get9660IsoData(f, callBack)
	return
}

//获取文件数据
func Get9660IsoData(f *os.File, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	if f == nil {
		err = errors.New("os.File is nil")
		return
	}

	isoReader, err := iso9660.NewReader(f)
	if err != nil {
		return
	}

	for {
		fs, err := isoReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}

		data := make([]byte, fs.Size())
		if fs.IsDir() {
			continue
		}
		reader := fs.Sys().(io.Reader)
		if reader == nil {
			continue
		}
		_, err = reader.Read(data)
		if err != nil && err != io.EOF {
			callBack(string(data), fs.Name())
		}
	}

	return
}
