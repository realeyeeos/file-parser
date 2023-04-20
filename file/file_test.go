package file

/*
Date：2023.03.02
Author：scl
Description：测试读取文件函数
*/

import (
	"fmt"
	"testing"
)

//go test -v -run ^TestDocx$ collector/file
//测试docx文件
func TestDocx(t *testing.T) {
	err := GetDocxData("C:\\Users\\lenovo\\Desktop\\预警-WA-20230120国网-001WebLogic远程代码执行漏洞（CVE-2023-21839）风险预警\\预警-WA-20230120国网-001WebLogic远程代码执行漏洞（CVE-2023-21839）风险预警.docx", CallBackData)
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
	err := GetOffice97Data("F:\\project_git\\dsp-fileplugin\\tmpfile\\Desktop\\测试doc.doc", CallBackData)
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
	err := GetPdfData("C:\\Users\\lenovo\\Desktop\\预警-WA-20230120国网-001WebLogic远程代码执行漏洞（CVE-2023-21839）风险预警\\预警-WA-20230120国网-001WebLogic远程代码执行漏洞（CVE-2023-21839）风险预警.pdf", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//go test -v -run ^TestPptx$ collector/file
//测试pptx文件
func TestPptx(t *testing.T) {
	err := GetPptxData("", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//go test -v -run ^TestRtf$ collector/file
//测试rtf文件
func TestRtf(t *testing.T) {
	err := GetRtfData("", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//go test -v -run ^TestXlsx$ collector/file
//测试xlsx文件
func TestXlsx(t *testing.T) {
	err := GetXlsxData("F:\\project_git\\dsp-fileplugin\\tmpfile\\Desktop\\测试elsx.xlsx", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//go test -v -run ^TestXps$ collector/file
//测试xps文件
func TestXps(t *testing.T) {
	err := GetXpsData("", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestTxt(t *testing.T) {
	err := GetTxtData("F:\\project_git\\dsp-fileplugin\\tmpfile\\测试_le.txt", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestEml(t *testing.T) {
	err := GetEmlData("F:\\project_git\\dsp-fileplugin\\tmpfile\\网上购票系统-用户支付通知.eml", CallBackData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func CallBackData(str, position string) bool {
	fmt.Println(position + "====" + str)
	return true
}
