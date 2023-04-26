package fileparser

/*
Date：2023.04.24
Author：scl
Description：解析rar文件
*/

import (
	"errors"
	"io/fs"

	"github.com/mholt/archiver/v4"
)

//获取文件数据(1-tar.gz 2-tar.bz2 3-tar.xz)
func GetRarDataFile(fileName string, fileType int, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	//打开文件
	fsys, err := archiver.FileSystem(fileName)
	if err != nil {
		return
	}

	//处理文件
	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		//打开文件
		file, err := fsys.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		//文件属性
		fi, err := d.Info()
		if err != nil || fi.Size() == 0 {
			return err
		}

		//获取文件数据
		data := make([]byte, fi.Size())
		_, err = file.Read(data)
		if err != nil {
			return err
		}

		callBack(string(data), path)
		return nil
	})
	return
}
