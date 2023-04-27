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

	err = GetRarData(f, callBack)
	return
}

func GetRarData(f *os.File, callBack ZipCallBack) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	if f == nil {
		err = errors.New("os.File is nil")
		return
	}

	rar := archiver.Rar{}
	// defer rar.Close()
	// rar.Walk(fileName, func(f archiver.File) error {
	// 	if f.Size() == 0 {
	// 		return errors.New("file size is 0")
	// 	}

	// 	//处理压缩包里的文件
	// 	callBack(f, f.Name(), f.Size())
	// 	return nil
	// })

	err = rar.Open(f, 0)
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

		//处理压缩包里的文件
		callBack(f, f.Name(), file.Size())

	}
	return
}
