package fileparser

/*
Date：2023.04.24
Author：scl
Description：解析7z、tar、rar（4+）、zip文件
*/

import (
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

	err = Get7zipData(f, callBack)
	return
}

//获取文件数据
func Get7zipData(fileReader io.Reader, callBack ZipCallBack) (err error) {
	if callBack == nil || fileReader == nil {
		err = errors.New("callback is nil or io.Reader is nil")
		return
	}

	zipreader, err := unarr.NewArchiveFromReader(fileReader)
	if err != nil {
		return
	}
	defer zipreader.Close()
	//获取文件名
	contens, err := zipreader.List()
	if err != nil {
		return
	}
	//循环文件名
	for _, v := range contens {
		err = zipreader.EntryFor(v)
		if err != nil {
			continue
		}

		//处理压缩包中的文件
		callBack(zipreader, v, int64(zipreader.Size()))
	}

	return
}
