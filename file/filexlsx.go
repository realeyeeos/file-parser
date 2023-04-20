package file

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

//打开文件
func GetXlsxData(fileName string, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return
	}

	xlsx, err := spreadsheet.Read(f, fi.Size())
	if err != nil {
		return
	}
	err = DealXlsxFile(xlsx, callBack)

	return
}

//处理xlsx文件
func DealXlsxFile(xlsx *spreadsheet.Workbook, callBack CallBackDataFunc) (err error) {
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
