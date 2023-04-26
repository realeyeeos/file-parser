package fileparser

/*
Date：2023.04.23
Author：scl
Description：解析html文件
*/
import (
	"errors"
	"os"

	"golang.org/x/net/html"
)

//获取文件数据
func GetHtmlDataFile(fileName string, callBack CallBackDataFunc) (err error) {
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
	err = GetHtmlData(f, callBack)
	return
}

//获取文件数据
func GetHtmlData(f *os.File, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	if f == nil {
		err = errors.New("os.File is nil")
		return
	}

	doc, err := html.Parse(f)
	if err != nil {
		return
	}

	var fun func(*html.Node)
	fun = func(n *html.Node) {
		if n.Type == html.TextNode && len(n.Data) != 0 {
			callBack(n.Data, "html")
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fun(c)
		}
	}
	fun(doc)
	return
}
