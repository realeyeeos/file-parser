package fileparser

/*
Date：2023.03.02
Author：scl
Description：解析xlsx文件
*/
import (
	"errors"
	"os"
	"strconv"

	"baliance.com/gooxml/spreadsheet"
)

//获取文件数据
func GetXlsxDataFile(fileName string, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer f.Close()

	//处理文件数据
	err = GetXlsxData(f, callBack)
	return
}

//获取文件数据
func GetXlsxData(f *os.File, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	if f == nil {
		err = errors.New("os.File is nil")
		return
	}

	fi, err := f.Stat()
	if err != nil {
		return
	}

	//获取xlsx文件句柄
	xlsx, err := spreadsheet.Read(f, fi.Size())
	if err != nil {
		return
	}

	//处理文件数据
	err = dealXlsxFile(xlsx, callBack)
	return
}

//处理xlsx文件
func dealXlsxFile(xlsx *spreadsheet.Workbook, callBack CallBackDataFunc) (err error) {
	defer xlsx.Close()

	stylesheet := xlsx.StyleSheet
	if stylesheet.X() == nil || stylesheet.X().CellXfs == nil || stylesheet.X().CellXfs.Xf == nil {
		return
	}

	//stysheetxfs := stylesheet.CellStyles()
	for _, sheet := range xlsx.Sheets() {
		name := sheet.Name()
		if len(name) == 0 {
			continue
		}

		if !callBack(string(name), "工作表名（"+name+")") {
			return nil
		}

		//行
		for k, row := range sheet.Rows() {
			str := ""
			for _, cell := range row.Cells() {
				if cell.IsEmpty() {
					continue
				}

				//获取一个单元格数据
				text := cell.GetFormattedValue()
				if len(text) == 0 {
					continue
				}

				str += text + "\t"

			}
			if len(str) == 0 {
				continue
			}

			if !callBack(str, "第"+strconv.Itoa(k+1)+"行") {
				return nil
			}

			str = ""
		}
	}

	return
}
