package fileparser

/*
Date：2023.03.02
Author：scl
Description：解析pdf文件
*/

import (
	"errors"
	"os"
	"strconv"

	"github.com/ledongthuc/pdf"
)

//获取文件数据
func GetPdfDataFile(fileName string, callBack CallBackDataFunc) (err error) {
	defer func() {
		if err := recover(); err != nil { // recover 捕获错误。
			return
		}
	}()
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	err = GetPdfData(f, callBack)
	return
}

//获取文件数据
func GetPdfData(f *os.File, callBack CallBackDataFunc) (err error) {
	defer func() {
		if err := recover(); err != nil { // recover 捕获错误。
			return
		}
	}()
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
		if err == nil {
			err = errors.New("file size is nil")
		}
		return err
	}

	//获取pdf文件句柄
	pdfReader, err := pdf.NewReader(f, fi.Size())
	if err != nil {
		return
	}

	//处理文件数据
	err = dealPdfFile(pdfReader, callBack)
	return
}

//处理pdf
func dealPdfFile(pdfReader *pdf.Reader, callBack CallBackDataFunc) (err error) {
	for pageIndex := 1; pageIndex <= pdfReader.NumPage(); pageIndex++ {
		p := pdfReader.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		text, err := p.GetPlainText(nil)
		if err != nil || len(text) == 0 {
			continue
		}

		if !callBack(text, "第"+strconv.Itoa(pageIndex+1)+"页") {
			return nil
		}
	}

	return
}
