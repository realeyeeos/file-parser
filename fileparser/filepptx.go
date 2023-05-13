package fileparser

/*
Date：2023.03.02
Author：scl
Description：解析pptx文件
*/

import (
	"errors"
	"io"
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
	fi, err := f.Stat()
	if err != nil || fi.Size() == 0 {
		if err == nil {
			err = errors.New("file size is nil")
		}
		return
	}

	err = GetPptxData(f, fi.Size(), callBack)
	return
}

//获取文件数据
func GetPptxData(fileReaderAt io.ReaderAt, fileSize int64, callBack CallBackDataFunc) (err error) {
	if callBack == nil || fileReaderAt == nil || fileSize == 0 {
		err = errors.New("callBack is nil or io.ReaderAt is nil or fileSize is 0")
		return
	}

	//获取pptx文件句柄
	pptx, err := presentation.Read(fileReaderAt, fileSize)
	if err != nil {
		return
	}

	//处理文件数据
	err = dealPptxFile(pptx, callBack)
	return
}

// 处理pptx文件
func dealPptxFile(ppt *presentation.Presentation, callBack CallBackDataFunc) (err error) {
	//图片
	// for _, v := range ppt.Images {
	// 	image, err := os.ReadFile(v.Path())
	// 	if err != nil {
	// 		continue
	// 	}

	// 	os.WriteFile("F:\\project_git\\dsp-fileplugin\\tmpfile\\scl\\ppt.jpeg", image, 0666)
	// 	//fmt.Println(v.Path())
	// }

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
