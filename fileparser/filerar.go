package fileparser

/*
Date：2023.04.24
Author：scl
Description：解析rar文件
*/

import (
	"errors"
	"io"
	"os"

	"github.com/mholt/archiver"
)

func GetRarDataFile(fileName string, callBack ZipCallBack) (err error) {
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

	// rar := archiver.Rar{}
	// rar.Unarchive(fileName, "C:\\Users\\lenovo\\Desktop\\123")

	err = GetRarData(f, fi.Size(), callBack)
	return
}

func GetRarData(fileReader io.Reader, fileSize int64, callBack ZipCallBack) (err error) {
	if callBack == nil || fileReader == nil {
		err = errors.New("callBack is nil or io.Reader is nil")
		return
	}

	rar := archiver.Rar{}
	defer rar.Close()

	err = rar.Open(fileReader, fileSize)
	if err != nil {
		return
	}
	defer rar.Close()

	for {
		file, err := rar.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		if file.IsDir() {
			continue
		}

		// 	//处理压缩包里的文件
		callBack(file, file.Name())
	}
	return
}
