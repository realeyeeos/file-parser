package fileparser

/*
Date：2023.03.02
Author：scl
Description：测试读取文件函数
*/

import (
	"fmt"
	"io"
	"testing"
)

//go test -v -run ^TestDocx$ collector/file
//测试docx文件
func TestDocx(t *testing.T) {
	err := GetDocxDataFile("C:\\Users\\lenovo\\Desktop\\预警-WA-20230120国网-001WebLogic远程代码执行漏洞（CVE-2023-21839）风险预警\\预警-WA-20230120国网-001WebLogic远程代码执行漏洞（CVE-2023-21839）风险预警.docx", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//go test -v -run ^TestOffice97$ collector/file
//测试office97（doc、xls、ppt）文件
func TestOffice97(t *testing.T) {
	//F:\\project_git\\dsp-fileplugin\\tmpfile\\47304.doc
	//F:\\project_git\\dsp-fileplugin\\tmpfile\\Desktop\\测试doc.doc
	//F:\\project_git\\dsp-fileplugin\\tmpfile\\测试.ppt
	err := GetOffice97DataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\47304.doc", CallBackData)
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
	err := GetPdfDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\火绒终端安全管理系统V2.0产品使用说明.pdf", CallBackData)
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
	err := GetXlsxDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\Desktop\\测试excel.xlsx", CallBackData)
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
	err := Get7zipDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\压缩包\\压缩包.7z", ZipCallBackFun)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestBz2(t *testing.T) {
	err := GetBz2DataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\压缩包\\123.txt.bz2", ZipCallBackFun)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestGz(t *testing.T) {
	err := GetGzDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\压缩包\\456.txt.gz", ZipCallBackFun)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestIso(t *testing.T) {
	err := GetIsoDataFile("E:\\vm\\镜像\\sc_winxp_pro_with_sp2.iso", ZipCallBackFun)
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
	err := GetRarDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\压缩包\\压缩包.rar", ZipCallBackFun)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestTarz(t *testing.T) {
	err := GetTarzDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\压缩包\\test.tar.gz", 1, ZipCallBackFun)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestXz(t *testing.T) {
	err := GetXzDataFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\压缩包\\123.txt.xz", ZipCallBackFun)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func CallBackData(str, position string) bool {
	fmt.Println(position + "====" + str)
	return true
}

func ZipCallBackFun(zipReader io.Reader, fileName string, fileSize int64) {
	if fileSize == 0 {
		return
	}

	data := make([]byte, fileSize)
	zipReader.Read(data)

	fmt.Println(fileName, ":", string(data))
}
