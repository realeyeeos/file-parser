package fileparser

/*
Date：2023.03.22
Author：scl
Description：解析eml文件
*/
import (
	"collector/emlparser"
	"errors"
	"io"
	"os"
)

//打开文件
func GetEmlData(fileName string, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return
	}

	//数据
	m, err := emlparser.Parse(file, false)
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
