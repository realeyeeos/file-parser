package fileparser

/*
Date：2023.03.02
Author：scl
Description：测试读取文件函数
*/

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gen2brain/go-fitz"
)

//go test -v -run ^TestDocx$ collector/file
//测试docx文件
func TestDocx(t *testing.T) {
	err := GetDocxDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\Agent组件文件命名及元数据字段规范.docx", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestScl(t *testing.T) {
	// file, err := os.Open("C:\\Users\\lenovo\\Desktop\\123.png")
	// if err != nil {
	// 	return
	// }
	// defer file.Close()
	// // 解码 PNG 图片文件
	// img, err := png.Decode(file)
	// if err != nil {
	// 	return
	// }

	// fmt.Println(img, str)

}

func TestPdfPng(t *testing.T) {
	document, err := fitz.New("F:\\project_git\\dsp-fileplugin\\tmpfile\\202304月考勤汇总公示版.pdf")
	if err != nil {
		return
	}

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
			os.WriteFile(filepath.Join("F:\\project_git\\dsp-fileplugin\\tmpfile\\scl", fmt.Sprintf("%03d.jpg", n)), buf.Bytes(), 0666)

			continue
		}

		fmt.Println(text)
	}

	return
}

//go test -v -run ^TestOffice97$ collector/file
//测试office97（doc、xls、ppt）文件
func TestOffice97(t *testing.T) {
	//F:\\project_git\\dsp-fileplugin\\tmpfile\\47304.doc
	//F:\\project_git\\dsp-fileplugin\\tmpfile\\Desktop\\测试doc.doc
	//F:\\project_git\\dsp-fileplugin\\tmpfile\\测试.ppt
	err := GetOffice97DataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\scl\\office.doc", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//go test -v -run ^TestPdf$ collector/file
//测试pdf文件
func TestPdf(t *testing.T) {
	defer func() {
		if err := recover(); err != nil { // recover 捕获错误。
			return
		}
	}()
	err := GetPdfDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\202304月考勤汇总公示版.pdf", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//go test -v -run ^TestPptx$ collector/file
//测试pptx文件
func TestPptx(t *testing.T) {
	err := GetPptxDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\研发部新员工转正_宋春良.pptx", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//go test -v -run ^TestRtf$ collector/file
//测试rtf文件
func TestRtf(t *testing.T) {
	err := GetRtfDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\测试.rtf", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//go test -v -run ^TestXlsx$ collector/file
//测试xlsx文件
func TestXlsx(t *testing.T) {
	err := GetXlsxDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\研发部-宋春良KPI.xlsx", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//go test -v -run ^TestXps$ collector/file
//测试xps文件
func TestXps(t *testing.T) {
	err := GetXpsDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\测试.xps", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestTxt(t *testing.T) {
	err := GetTxtDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\测试_le.txt", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// func TestEml(t *testing.T) {
// 	err := GetEmlDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\网上购票系统-用户支付通知.eml", CallBackData)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// }

func TestHtml(t *testing.T) {
	err := GetHtmlDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\108717.html", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestXml(t *testing.T) {
	err := GetXmlDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\sheet1.xml", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func Test7zip(t *testing.T) {
	err := Get7zipDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\压缩包\\excel.7z", ZipCallBackFun)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestBz2(t *testing.T) {
	err := GetBz2DataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\压缩包\\火绒终端安全管理系统V2.0产品使用说明.pdf.bz2", ZipCallBackFun)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestGz(t *testing.T) {
	err := GetGzDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\压缩包\\docrar.rar.gz", ZipCallBackFun)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestIso(t *testing.T) {
	err := GetIsoDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\压缩包\\20230504_114724.iso", ZipCallBackFun)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestLzma(t *testing.T) {
	err := GetLzmaDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\压缩包\\123.txt.lzma", ZipCallBackFun)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestOdt(t *testing.T) {
	err := GetOdtDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\测试ansi.odt", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestRar(t *testing.T) {
	err := GetRarDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\压缩包\\测试ansi.rar", ZipCallBackFun)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestTarz(t *testing.T) {
	err := GetTarzDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\压缩包\\test.tar.gz", ZipCallBackFun)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestXz(t *testing.T) {
	err := GetXzDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\压缩包\\测试ansi.doc.gz.xz", ZipCallBackFun)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func CallBackData(str, position string) bool {
	fmt.Println(position + "====" + str)
	return true
}

func ZipCallBackFun(zipReader io.Reader, fileName string) {
	splits := strings.Split(fileName, ".")
	if len(splits) == 0 {
		return
	}

	fileType := strings.ToLower(splits[len(splits)-1])

	var err error
	var buf bytes.Buffer
	//不使用io.copy会导致获取数据不全
	_, err = io.Copy(&buf, zipReader)
	if err != nil && err != io.EOF {
		return
	}
	byteReader := bytes.NewReader(buf.Bytes())

	if fileType == "doc" || fileType == "xls" || fileType == "ppt" {
		readSeeker := io.ReadSeeker(byteReader)
		err = GetOffice97Data(readSeeker, CallBackData)
	} else if fileType == "docx" {
		readerAt := io.ReaderAt(byteReader)
		err = GetDocxData(readerAt, int64(buf.Len()), CallBackData)
	} else if fileType == "xlsx" {
		readerAt := io.ReaderAt(byteReader)
		err = GetXlsxData(readerAt, int64(buf.Len()), CallBackData)
	} else if fileType == "pptx" {
		readerAt := io.ReaderAt(byteReader)
		err = GetPptxData(readerAt, int64(buf.Len()), CallBackData)
	} else if fileType == "pdf" {
		err = GetPdfData(byteReader, int64(buf.Len()), CallBackData)
	} else if fileType == "7z" || fileType == "tar" || fileType == "zip" {
		err = Get7zipData(byteReader, ZipCallBackFun)
	} else if fileType == "rar" {
		err = GetRarData(byteReader, int64(buf.Len()), ZipCallBackFun)
	} else if fileType == "bz2" {
		err = GetBz2Data(byteReader, fileName, ZipCallBackFun)
	} else if fileType == "gz" {
		err = GetGzData(byteReader, ZipCallBackFun)
	} else if fileType == "xz" {
		err = GetXzData(byteReader, fileName, ZipCallBackFun)
	} else {
		fmt.Println(fileName, ":", buf.String())
	}
	if err != nil {
		fmt.Println(err)
	}
}
