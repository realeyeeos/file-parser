package fileparser

/*
Date：2023.03.02
Author：scl
Description：解析xps文件
*/
import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

//总体xml
type FIXEDDOCUMENTSEQUENCE struct {
	Document struct {
		Source string `xml:"Source,attr"`
	} `xml:"DocumentReference"`
}

//保存数据文件的xml
type FIXEDDOCUMENT struct {
	PageContent struct {
		Source string `xml:"Source,attr"`
	} `xml:"PageContent"`
}

//数据文件
type FPAGE struct {
	Glyphs []struct {
		UnicodeString string `xml:"UnicodeString,attr"`
	} `xml:"Glyphs"`
}

//获取zip中文件数据
func GetZipFileData(zipfile *zip.File) (fileData []byte, err error) {
	fixedfile, err := zipfile.Open()
	if err != nil {
		err = errors.New("FixedDocumentSequence open error")
		return
	}
	defer fixedfile.Close()

	fileinfo := zipfile.FileInfo()
	fixedlen := fileinfo.Size()
	fileData = make([]byte, fixedlen)

	_, err = fixedfile.Read(fileData)
	if err != nil && err != io.EOF {
		return
	}
	return fileData, nil
}

//获取xps的数据
func GetFpageInfo(mapvalues map[string]*zip.File, callBack CallBackDataFunc) (err error) {
	fixedDocument := mapvalues["/FixedDocumentSequence.fdseq"]
	if fixedDocument == nil {
		err = errors.New("FixedDocumentSequence is not find")
		return
	}

	//获取zip中最外层xml的数据（FixedDocumentSequence.fdseq）
	fixed_sequence_data, err := GetZipFileData(fixedDocument)
	if err != nil {
		return
	}

	//反序列化读取主xml文件内容
	var fixed_sequence_xmldata FIXEDDOCUMENTSEQUENCE
	xml.Unmarshal(fixed_sequence_data, &fixed_sequence_xmldata)

	//获取保存数据的xml文件（Documents\1\FixedDocument.fdoc）
	pageFixedDocument := mapvalues[fixed_sequence_xmldata.Document.Source]
	if pageFixedDocument == nil {
		err = errors.New("FixedDocumentSequence Document.Source is not find")
		return
	}

	//获取 FixedDocument.fdoc 文件数据
	fixedData, err := GetZipFileData(pageFixedDocument)
	if err != nil {
		return
	}

	//反序列化 FixedDocument.fdoc 文件
	var fixedXmlData FIXEDDOCUMENT
	xml.Unmarshal(fixedData, &fixedXmlData)

	//获取保存数据的文件名（Documents\1\Pages\1.fpage）
	fpage_zipfile := mapvalues[fixedXmlData.PageContent.Source]
	if fpage_zipfile == nil {
		err = errors.New("FixedDocument PageContent.Source is not find")
		return
	}

	//获取保存数据文件的数据（Documents\1\Pages\1.fpage）
	ffpageData, err := GetZipFileData(fpage_zipfile)
	if err != nil {
		return
	}

	//反序列化"1.fpage"文件
	var fpage FPAGE
	xml.Unmarshal(ffpageData, &fpage)

	//获取所有数据
	for k, v := range fpage.Glyphs {
		if len(v.UnicodeString) == 0 {
			continue
		}

		if !callBack(v.UnicodeString, "第"+strconv.Itoa(k+1)+"行") {
			return nil
		}
	}

	return
}

//获取文件数据
func GetXpsDataFile(fileName string, callBack CallBackDataFunc) (err error) {
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
	err = GetXpsData(f, callBack)
	return
}

//获取文件数据
func GetXpsData(f *os.File, callBack CallBackDataFunc) (err error) {
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
		return
	}

	//创建一个zip的reader
	zipReader, err := zip.NewReader(f, fi.Size())
	if err != nil {
		return
	}

	var isXps bool
	for _, v := range zipReader.File {
		if strings.Contains(v.Name, "FixedDocumentSequence.fdseq") {
			isXps = true
			break
		}
	}

	if !isXps {
		err = errors.New("is not xps")
		return
	}

	//处理文件数据
	err = dealXpsFile(zipReader, callBack)
	return
}

//处理xps文件
func dealXpsFile(zipreader *zip.Reader, callBack CallBackDataFunc) (err error) {
	//所有文件及其指针保存到map中，方便后续查找
	var mapValues map[string]*zip.File = make(map[string]*zip.File)
	for _, v := range zipreader.File {
		mapValues["/"+v.Name] = v
	}

	//获取文档中的数据
	err = GetFpageInfo(mapValues, callBack)
	if err != nil {
		return
	}

	return
}
