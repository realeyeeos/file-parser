package file

/*
Date：2023.03.02
Author：scl
Description：解析pdf文件
*/

import (
	"bytes"
	"errors"
	"os"
	"strconv"

	"github.com/ledongthuc/pdf"
)

// ReadPdf 获取pdf文字内容
func ReadPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	// remember close file
	defer f.Close()
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer

	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	return buf.String(), nil
}

//阅读按行分组的文本
// func ReadPdfGroup(path string) (string, error) {
// 	f, r, err := pdf.Open(path)
// 	defer func() {
// 		_ = f.Close()
// 	}()
// 	if err != nil {
// 		return "", err
// 	}
// 	totalPage := r.NumPage()

// 	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
// 		p := r.Page(pageIndex)
// 		if p.V.IsNull() {
// 			continue
// 		}

// 		rows, _ := p.GetTextByRow()
// 		for _, row := range rows {
// 			println(">>>> row: ", row.Position)
// 			for _, word := range row.Content {
// 				fmt.Print(word.S)
// 			}
// 		}
// 	}
// 	return "", nil
// }

//打开文件
func GetPdfData(fileName string, callBack CallBackDataFunc) (err error) {
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
	if err != nil {
		return err
	}

	pdfReader, err := pdf.NewReader(f, fi.Size())
	if err != nil {
		return
	}

	err = DealPdfFile(pdfReader, callBack)

	return
}

//处理pdf
func DealPdfFile(pdfReader *pdf.Reader, callBack CallBackDataFunc) (err error) {
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
