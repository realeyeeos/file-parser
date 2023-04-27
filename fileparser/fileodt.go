package fileparser

/*
Date：2023.04.27
Author：scl
Description：解析odt文件
*/

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"io"
	"os"
)

//获取文件数据
func GetOdtDataFile(fileName string, callBack CallBackDataFunc) (err error) {
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

	err = GetOdtData(f, callBack)
	return
}

//获取文件数据
func GetOdtData(f *os.File, callBack CallBackDataFunc) (err error) {
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
		return
	}

	zipReader, err := zip.NewReader(f, fi.Size())
	if err != nil {
		return
	}

	//处理压缩文件，获取其中的数据文件
	var filaReader io.Reader
	for _, v := range zipReader.File {
		if v.Name == "content.xml" {
			filaReader, err = v.Open()
			if err != nil && err != io.EOF {
				break
			}
		}
	}

	if filaReader == nil {
		return
	}

	//解析xml文件中的数据
	decoder := xml.NewDecoder(filaReader)
	t, err := decoder.Token()
	for err == nil {
		switch token := t.(type) {
		case xml.CharData:
			if len(token) > 0 {
				callBack(string(token), "")
			}
		}

		t, err = decoder.Token()
	}

	return
}
