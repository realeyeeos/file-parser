package fileparser

/*
Date：2023.04.23
Author：scl
Description：解析xml文件
*/

import (
	"encoding/xml"
	"errors"
	"io"
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
func GetXmlData(fileReader io.Reader, callBack CallBackDataFunc) (err error) {
	if callBack == nil || fileReader == nil {
		err = errors.New("callBack is nil or io.Reader is nil")
		return
	}

	decoder := xml.NewDecoder(fileReader)
	t, err := decoder.Token()
	for err == nil {
		switch token := t.(type) {
		case xml.StartElement:
			var xmlstr string
			for _, v := range token.Attr {
				xmlstr += v.Value + " "
			}

			if len(xmlstr) != 0 {
				callBack(xmlstr, token.Name.Space+token.Name.Local)
			}
		case xml.CharData:
			if len(token) > 0 {
				callBack(string(token), "")
			}
		}

		t, err = decoder.Token()
	}

	return
}
