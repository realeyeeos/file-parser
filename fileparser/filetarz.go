package fileparser

/*
Date：2023.04.24
Author：scl
Description：解析tar.gz、tar.bz2、tar.xz文件
*/

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"io"
	"os"

	"github.com/ulikunitz/xz"
)

//获取文件数据(1-tar.gz 2-tar.bz2 3-tar.xz)
func GetTarzDataFile(fileName string, fileType int, callBack CallBackDataFunc) (err error) {
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

	err = GetTarzData(f, fileType, callBack)
	return
}

//获取文件数据
func GetTarzData(f *os.File, fileType int, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	if f == nil {
		err = errors.New("os.File is nil")
		return
	}

	var tarread *tar.Reader
	if fileType == 1 {
		gr, err := gzip.NewReader(f)
		if err != nil {
			return err
		}

		defer gr.Close()
		//使用tar解析其中的tar文件
		tarread = tar.NewReader(gr)
	} else if fileType == 2 {
		bz := bzip2.NewReader(f)
		tarread = tar.NewReader(bz)
	} else if fileType == 3 {
		xzReader, err := xz.NewReader(f)
		if err != nil {
			return err
		}
		tarread = tar.NewReader(xzReader)
	}

	if tarread == nil {
		err = errors.New("reader is nil")
		return
	}

	for {
		trgzfile, err := tarread.Next()
		if err != nil {
			break
		}

		//文件大小
		if trgzfile.Size == 0 {
			continue
		}
		//文件数据
		data := make([]byte, trgzfile.Size)
		_, err = tarread.Read(data)
		if err != nil && err != io.EOF {
			continue
		}
	}

	return
}
