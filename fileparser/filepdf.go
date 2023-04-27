package fileparser

/*
Date：2023.03.02
Author：scl
Description：解析pdf文件
*/

import (
	"errors"
	"io"
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
	fi, err := f.Stat()
	if err != nil || fi.Size() == 0 {
		if err == nil {
			err = errors.New("file size is nil")
		}
		return err
	}

	err = GetPdfData(f, fi.Size(), callBack)
	return
}

//获取文件数据
func GetPdfData(fileReaderAt io.ReaderAt, fileSize int64, callBack CallBackDataFunc) (err error) {
	defer func() {
		if err := recover(); err != nil { // recover 捕获错误。
			return
		}
	}()
	if callBack == nil || fileReaderAt == nil || fileSize == 0 {
		err = errors.New("callBack is nil or io.ReaderAt is nil or fileSize is 0")
		return
	}

	//获取pdf文件句柄
	pdfReader, err := pdf.NewReader(fileReaderAt, fileSize)
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
