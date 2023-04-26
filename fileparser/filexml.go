package fileparser

/*
Date：2023.04.23
Author：scl
Description：解析docx文件
*/

import (
	"encoding/xml"
	"errors"
	"os"
)

//获取文件数据
func GetXmlDataFile(fileName string, callBack CallBackDataFunc) (err error) {
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

	err = GetXmlData(f, callBack)
	return
}

//获取文件数据
func GetXmlData(f *os.File, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	if f == nil {
		err = errors.New("os.File is nil")
		return
	}

	decoder := xml.NewDecoder(f)
	t, err := decoder.Token()
	for err == nil {
		var xmlstr string
		switch token := t.(type) {
		case xml.StartElement:
			for _, v := range token.Attr {
				xmlstr += v.Value + " "
			}

			if len(xmlstr) != 0 {
				callBack(xmlstr, token.Name.Local)
			}
		}
		t, err = decoder.Token()
	}

	return
}
