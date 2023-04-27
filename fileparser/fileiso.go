package fileparser

/*
Date：2023.04.25
Author：scl
Description：解析udf格式的iso文件
*/

import (
	"encoding/binary"
	"errors"
	"io"
	"os"

	"github.com/hooklift/iso9660"
	"github.com/mogaika/udf"
)

//获取文件数据
func GetIsoDataFile(fileName string, callBack ZipCallBack) (err error) {
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

	err = GetIsoData(f, callBack)
	return
}

//获取文件数据
func GetIsoData(f *os.File, callBack ZipCallBack) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	if f == nil {
		err = errors.New("os.File is nil")
		return
	}

	var tagData [2]byte
	_, err = f.ReadAt(tagData[:], 2048*256)
	if err != nil {
		return
	}
	f.Seek(0, io.SeekStart)

	//不是udf格式文件
	tagIdentifier := binary.LittleEndian.Uint16(tagData[:])
	if tagIdentifier != udf.DESCRIPTOR_ANCHOR_VOLUME_POINTER {
		//9660格式iso文件读取
		Read9660File(f, callBack)
		return
	}

	//udf格式iso文件读取
	udfReader := udf.NewUdfFromReader(f)
	//递归读取iso中的文件
	ReadUdfFile(udfReader, nil, callBack)

	return
}

//递归读取iso中的文件
func ReadUdfFile(udfReader *udf.Udf, fe *udf.FileEntry, callBack ZipCallBack) {
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
					//处理压缩包里的文件
					callBack(j.NewReader(), j.Name(), j.Size())
				}
			}
		} else {
			//大小：j.NewReader().Size()
			//数据：j.NewReader().Read()
			//处理压缩包里的文件
			callBack(v.NewReader(), v.Name(), v.Size())
		}
	}
}

//获取文件数据
func Read9660File(f *os.File, callBack ZipCallBack) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	if f == nil {
		err = errors.New("os.File is nil")
		return
	}

	isoReader, err := iso9660.NewReader(f)
	if err != nil {
		return
	}

	for {
		fs, err := isoReader.Next()
		if err != nil {
			break
		}

		if fs.Size() == 0 {
			continue
		}

		//读取文件数据
		//data := make([]byte, fs.Size())
		if fs.IsDir() {
			continue
		}
		reader := fs.Sys().(io.Reader)

		//处理压缩包里的文件
		callBack(reader, fs.Name(), fs.Size())
	}

	return
}
