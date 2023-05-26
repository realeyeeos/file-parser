package fileparser

/*
Date：2023.04.24
Author：scl
Description：解析7z、tar、rar（4+）、zip文件
*/

import (
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/gen2brain/go-unarr"
)

//获取文件数据
func Get7zipDataFile(fileName string, callBack ZipCallBack) (err error) {
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

	err = Get7zipData(f, 0, callBack)
	return
}

//获取文件数据
//fileSize -无用，为了统一格式
func Get7zipData(fileReader io.Reader, fileSize int64, callBack ZipCallBack) (err error) {
	if callBack == nil || fileReader == nil {
		err = errors.New("callback is nil or io.Reader is nil")
		return
	}

	zipReader, err := unarr.NewArchiveFromReader(fileReader)
	if err != nil {
		return
	}
	defer zipReader.Close()
	//获取文件名
	contens, err := zipReader.List()
	if err != nil {
		return
	}
	//循环文件名
	for _, v := range contens {
		err = zipReader.EntryFor(v)
		if err != nil {
			continue
		}

		//读取所有数据，外边读取的话会有读取不全问题
		data, err := zipReader.ReadAll()
		if err != nil && err != io.EOF {
			continue
		}

		//处理压缩包中的文件
		callBack(bytes.NewReader(data), v)
	}

	return
}
