package fileparser

/*
Date：2023.03.02
Author：scl
Description：解析pptx文件
*/

import (
	"errors"
	"os"
	"strconv"

	"baliance.com/gooxml/presentation"
)

//获取文件数据
func GetPptxDataFile(fileName string, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer f.Close()

	err = GetPptxData(f, callBack)
	return
}

//获取文件数据
func GetPptxData(f *os.File, callBack CallBackDataFunc) (err error) {
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

	//获取pptx文件句柄
	pptx, err := presentation.Read(f, fi.Size())
	if err != nil {
		return
	}

	//处理文件数据
	err = dealPptxFile(pptx, callBack)
	return
}

// 处理pptx文件
func dealPptxFile(ppt *presentation.Presentation, callBack CallBackDataFunc) (err error) {
	for k, slide := range ppt.Slides() {
		str := ""
		for _, choice := range slide.X().CSld.SpTree.Choice {
			if choice.Sp == nil {
				continue
			}
			for _, sp := range choice.Sp {
				if sp.TxBody == nil {
					continue
				}
				for _, p := range sp.TxBody.P {
					textrun := p.EG_TextRun
					var text string
					for _, run := range textrun {
						if run.R != nil {
							text += run.R.T
						}
					}
					if len(text) == 0 {
						continue
					}

					str += text
				}
			}
		}
		if len(str) == 0 {
			continue
		}

		if !callBack(str, "第"+strconv.Itoa(k+1)+"页") {
			return nil
		}
	}

	return
}
