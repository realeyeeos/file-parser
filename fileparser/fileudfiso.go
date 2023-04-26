package fileparser

/*
Date：2023.04.25
Author：scl
Description：解析udf格式的iso文件
*/

import (
	"errors"
	"os"

	"github.com/mogaika/udf"
)

//获取文件数据
func GetUdfIsoDataFile(fileName string, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	//打开文件
	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer f.Close()

	err = GetUdfIsoData(f, callBack)
	return
}

//获取文件数据
func GetUdfIsoData(f *os.File, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	if f == nil {
		err = errors.New("os.File is nil")
		return
	}

	// var tagData [2]byte
	// _, err = f.ReadAt(tagData[:], 2048*256)
	// if err != nil {
	// 	return
	// }

	// //不是udf格式文件
	// tagIdentifier := binary.LittleEndian.Uint16(tagData[:])
	// if tagIdentifier != udf.DESCRIPTOR_ANCHOR_VOLUME_POINTER {
	// 	return
	// }

	// f.Seek(0, io.SeekStart)

	udfReader := udf.NewUdfFromReader(f)

	//递归读取iso中的文件
	ReadUdfFile(udfReader, nil, callBack)

	return
}

//递归读取iso中的文件
func ReadUdfFile(udfReader *udf.Udf, fe *udf.FileEntry, callBack CallBackDataFunc) {
	defer func() {
		if err := recover(); err != nil { // recover 捕获错误。
			return
		}
	}()
	for _, v := range udfReader.ReadDir(fe) {
		if v.Fid.FileCharacteristics&0x2 == 0x2 {
			//subPathName := pathName + "\\" + v.Name()
			for _, j := range v.ReadDir() {
				if j.Fid.FileCharacteristics&0x2 == 0x2 {
					//ReadUdfFile(udfReader, j.FileEntry(), subPathName+"\\"+j.Name())
					ReadUdfFile(udfReader, j.FileEntry(), callBack)
				} else {
					//大小：j.NewReader().Size()
					//数据：j.NewReader().Read()

					//callBack(string(data), "")
				}
			}
		} else {
			//大小：j.NewReader().Size()
			//数据：j.NewReader().Read()
			//callBack(string(data), "")
		}
	}
}
