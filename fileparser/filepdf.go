package fileparser

/*
Date：2023.03.02
Author：scl
Description：解析pdf文件
*/

import (
	"bytes"
	"errors"
	"image/jpeg"
	"io"
	"os"
	"strconv"

	"github.com/gen2brain/go-fitz"
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
func GetPdfData(fileReader io.Reader, fileSize int64, callBack CallBackDataFunc) (err error) {
	defer func() {
		if err := recover(); err != nil { // recover 捕获错误。
			return
		}
	}()
	if callBack == nil || fileReader == nil || fileSize == 0 {
		err = errors.New("callBack is nil or io.ReaderAt is nil or fileSize is 0")
		return
	}

	//获取pdf文件句柄
	// pdfReader, err := pdf.NewReader(fileReaderAt, fileSize)
	// if err != nil {
	// 	return
	// }

	document, err := fitz.NewFromReader(fileReader)
	if err != nil {
		return
	}

	//处理文件数据
	err = dealPdfFile(document, callBack)
	return
}

//处理pdf
func dealPdfFile(document *fitz.Document, callBack CallBackDataFunc) (err error) {
	// for pageIndex := 1; pageIndex <= pdfReader.NumPage(); pageIndex++ {
	// 	p := pdfReader.Page(pageIndex)
	// 	if p.V.IsNull() {
	// 		continue
	// 	}

	// 	text, err := p.GetPlainText(nil)
	// 	if err != nil || len(text) == 0 {
	// 		continue
	// 	}

	// 	if !callBack(text, "第"+strconv.Itoa(pageIndex)+"页") {
	// 		return nil
	// 	}
	// }

	//提取文字
	for n := 0; n < document.NumPage(); n++ {
		text, err := document.Text(n)
		if err != nil || len(text) == 0 {
			//提取图片（可以直接识别文字的也会提取出来）
			img, err := document.Image(n)
			if err != nil {
				continue
			}

			var buf bytes.Buffer
			err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
			if err != nil {
				continue
			}
			//os.WriteFile(filepath.Join("F:\\project_git\\dsp-fileplugin\\tmpfile\\scl", fmt.Sprintf("%03d.jpg", n)), buf.Bytes(), 0666)

			continue
		}

		if !callBack(text, "第"+strconv.Itoa(n+1)+"页") {
			return nil
		}
	}

	return
}
