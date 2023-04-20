package fileparser

/*
Date：2023.03.22
Author：scl
Description：解析eml文件
*/
import (
	"errors"
	"fileparser/emlparser"
	"io"
	"os"
)

func GetEmlDataFile(fileName string, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	err = GetEmlData(file, callBack)
	return
}

//获取文件数据
func GetEmlData(f *os.File, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	if f == nil {
		err = errors.New("os.File is nil")
		return
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return
	}

	//数据
	m, err := emlparser.Parse(f, false)
	if err != nil {
		return err
	}

	if !callBack(string(m.Text), "text") {
		return
	}

	if !callBack(string(m.Html), "html") {
		return
	}

	return
}
