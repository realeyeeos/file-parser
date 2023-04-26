package fileparser

/*
Date：2023.03.02
Author：scl
Description：解析docx文件
*/

import (
	"errors"
	"io"
	"os"
	"strconv"

	"baliance.com/gooxml/document"
	"github.com/beevik/etree"
)

//获取文件数据
func GetDocxDataFile(fileName string, callBack CallBackDataFunc) (err error) {
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

	err = GetDocxData(f, callBack)
	return
}

//获取文件数据
func GetDocxData(f *os.File, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	if f == nil {
		err = errors.New("os.File is nil")
		return
	}

	//获取文件属性
	fi, err := f.Stat()
	if err != nil || fi.Size() == 0 {
		if err == nil {
			err = errors.New("file size is nil")
		}
		return
	}

	//获取docx文件句柄
	docx, err := document.Read(f, fi.Size())
	if err != nil && err != io.EOF {
		return
	}

	//处理文件数据
	err = dealDocxFile(docx, callBack)
	return
}

//处理docx文件
func dealDocxFile(docx *document.Document, callBack CallBackDataFunc) (err error) {
	//批注
	for _, docfile := range docx.DocBase.ExtraFiles {
		if docfile.ZipPath != "word/comments.xml" {
			continue
		}

		file, err := os.Open(docfile.DiskPath)
		if err != nil {
			continue
		}
		defer file.Close()
		f, err := file.Stat()
		if err != nil {
			continue
		}
		size := f.Size()
		fileinfo := make([]byte, size)
		_, err = file.Read(fileinfo)
		if err != nil {
			continue
		}

		docment := etree.NewDocument()
		err = docment.ReadFromBytes(fileinfo)
		if err != nil {
			continue
		}

		root := docment.SelectElement("w:comments")

		for _, coment := range root.SelectElements("w:comment") {
			wp := coment.SelectElement("w:p")
			if wp == nil {
				continue
			}
			wr := wp.SelectElement("w:r")
			if wr == nil {
				continue
			}

			wt := wr.SelectElement("w:t")
			if wt == nil {
				continue
			}

			if len(wt.Text()) == 0 {
				continue
			}

			if !callBack(wt.Text(), "批注") {
				return nil
			}
		}
	}

	//书签
	for _, bookmark := range docx.Bookmarks() {
		bookname := bookmark.Name()
		if len(bookname) == 0 {
			continue
		}

		if !callBack(bookmark.Name(), "书签") {
			return nil
		}
	}

	//页眉
	for _, head := range docx.Headers() {
		var text string
		for _, para := range head.Paragraphs() {
			for _, run := range para.Runs() {
				text += run.Text()
			}
		}
		if len(text) == 0 {
			continue
		}

		if !callBack(text, "页眉") {
			return nil
		}
	}

	//页脚
	for _, footer := range docx.Footers() {
		for _, para := range footer.Paragraphs() {
			var text string
			for _, run := range para.Runs() {
				text += run.Text()
			}
			if len(text) == 0 {
				continue
			}

			if !callBack(text, "页脚") {
				return nil
			}
		}
	}

	//doc.Paragraphs()得到包含文档所有的段落的切片
	for k, para := range docx.Paragraphs() {
		var text string
		//run为每个段落相同格式的文字组成的片段
		for _, run := range para.Runs() {
			text += run.Text()
		}
		if len(text) == 0 {
			continue
		}

		if !callBack(text, "第"+strconv.Itoa(k+1)+"行") {
			return nil
		}
	}

	//获取表格
	for i, table := range docx.Tables() {
		for k, run := range table.Rows() {
			var text string
			for _, cell := range run.Cells() {
				if len(text) != 0 {
					text += "\t"
				}
				for _, para := range cell.Paragraphs() {
					for _, run := range para.Runs() {
						text += run.Text()
					}
				}
			}

			if len(text) == 0 {
				continue
			}

			if !callBack(text, "第"+strconv.Itoa(i+1)+"个表格第"+strconv.Itoa(k)+"行") {
				return nil
			}
		}
	}

	return
}
