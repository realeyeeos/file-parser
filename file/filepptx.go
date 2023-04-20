package file

/*
Date：2023.03.02
Author：scl
Description：解析pptx文件
*/

import (
	"errors"
	"io/fs"
	"os"
	"strconv"

	"baliance.com/gooxml/presentation"
)

//打开文件
func GetPptxData(fileName string, callBack CallBackDataFunc) (err error) {
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

	pptx, err := presentation.Read(f, fi.Size())
	if err != nil {
		return
	}

	err = DealPptxFile(pptx, callBack)

	return
}

//判断是否是pptx文件，获取句柄
func GetPptxHandle(fp *os.File, fi fs.FileInfo) (*presentation.Presentation, error) {
	ppt, err := presentation.Read(fp, fi.Size())
	if err != nil {
		return nil, err
	}
	return ppt, nil
}

// 处理pptx文件
func DealPptxFile(ppt *presentation.Presentation, callBack CallBackDataFunc) (err error) {
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
