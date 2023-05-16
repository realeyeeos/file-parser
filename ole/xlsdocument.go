package ole

/*
Date：2023.03.02
Author：scl
Description：解析xls文件
*/

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	publicfun "github.com/realeyeeos/file-parser/publicfunc"
)

//xls结构中的几个关键值
const (
	XLS_BOF          = 0x809
	XLS_CODEPAGE     = 0x42
	XLS_DATE_1904    = 0x22
	XLS_FILEPASS     = 0x2F
	XLS_FORMAT       = 0x41E
	XLS_FORMULA      = 0x06
	XLS_INTEGER_CELL = 0x202
	XLS_RSTRING      = 0xD6

	//Globals
	XLS_SHEET    = 0x85
	XLS_CONTINUE = 0x3C
	XLS_SST      = 0xFC

	//Pic
	XLS_DRAWING = 0xEB

	XLS_STRING = 0x207
	XLS_XF     = 0xE0
	XLS_EOF    = 0x0A
)

//xls版本
const (
	BOF_BIFF_2       = 0x009
	BOF_BIFF_3       = 0x209
	BOF_BIFF_4       = 0x0409
	BOF_BIFF_5_AND_8 = XLS_BOF
)

//获取xls文件数据
func (ole *OleInfo) getXlsInfo(reader io.ReadSeeker, xls *Directory, callBack DataCallBackFunc) (err error) {
	var isSSTContinue, isPicContinue bool
	var picData []byte
	//获取sheet名字、位置 + SST中的数据
	for {

		if !isPicContinue && len(picData) > 8 {
			//TODO Pictures By Scl
			// picReader := bytes.NewReader(picData[:])
			// ole.getArtData(nil, picReader, false)
			// picData = nil
		}

		var urec_type, urec_len uint16
		//类型
		err = binary.Read(reader, binary.LittleEndian, &urec_type)
		if err != nil {
			return
		}
		//大小
		err = binary.Read(reader, binary.LittleEndian, &urec_len)
		if err != nil {
			return
		}

		//结束符
		if urec_type == XLS_EOF {
			break
		}

		if urec_len == 0 {
			continue
		}

		var data []byte
		if urec_type == XLS_BOF || urec_type == XLS_SHEET || urec_type == XLS_CONTINUE || urec_type == XLS_SST ||
			urec_type == XLS_DRAWING {
			data = make([]byte, urec_len)
			n, err := reader.Read(data[:])
			if err != nil || n != int(urec_len) {
				if err == nil {
					err = errors.New("read len is error")
				}
				return err
			}
		}

		//开始符
		if urec_type == XLS_BOF {
			isSSTContinue = false
			isPicContinue = false
			if urec_len < 2 {
				err = errors.New("ureclen is short")
				return
			}

			ole.xlsInfo.xlsVersion = data[0:2]
		} else if urec_type == XLS_SHEET {
			isSSTContinue = false
			isPicContinue = false
			//获取sheet信息（包括名字和偏移量）
			ole.getSheetInfo(data)
		} else if urec_type == XLS_CONTINUE {
			if isSSTContinue { //数据
				ole.getSSTInfo(data)
			} else if isPicContinue { //图片
				picData = append(picData, data...)
			}
		} else if urec_type == XLS_SST {
			isPicContinue = false
			if len(data) <= 8 {
				continue
			}
			cstTotal := data[0:4]
			//总共数量
			ucstTotal := binary.LittleEndian.Uint32(cstTotal[:])

			//获取SST数据
			ole.getSSTInfo(data[8:])
			if len(ole.xlsInfo.rgbStrings) < int(ucstTotal-1) {
				isSSTContinue = true
			}
		} else if urec_type == XLS_DRAWING { //图片数据
			//TODO Pictures By Scl
			// isSSTContinue = false
			// picData = append(picData, data...)
			// if urec_len == 8224 {
			// 	isPicContinue = true
			// } else {
			// 	isPicContinue = false
			// }
		} else {
			isSSTContinue = false
			isPicContinue = false
			_, err = reader.Seek(int64(urec_len), io.SeekCurrent)
			if err != nil {
				return
			}
		}
	}

	//解析sheet流
	for _, v := range ole.xlsInfo.sheetNames {
		err := ole.dealSheet(reader, v, callBack)
		if err != nil {
			continue
		}
	}
	//rgbstrings = nil

	return
}

//获取SST数据
func (ole *OleInfo) getSSTInfo(data []byte) (err error) {
	sstReader := bytes.NewReader(data[:])
	//rgbstrings = nil
	for {
		//字符数量
		var cchLen uint16
		err = binary.Read(sstReader, binary.LittleEndian, &cchLen)
		if err != nil {
			return
		}

		var flags [1]byte
		n, err := sstReader.Read(flags[:])
		if err != nil || n != 1 {
			if err == nil {
				err = errors.New("read len is error")
			}
			return err
		}

		//0-ansi	1-unicode
		grbit := flags[:1][0]
		if cchLen == 0 {
			continue
		}
		stringlen := cchLen
		var cRun uint16
		var cbExtRst uint32
		//fHighByte
		if grbit&1 == 1 {
			stringlen = cchLen * 2
		}

		//fRichSt
		if grbit&0x8 == 0x8 {
			err = binary.Read(sstReader, binary.LittleEndian, &cRun)
			if err != nil {
				return err
			}
		}
		//fExtSt
		if grbit&0x4 == 0x4 {
			err = binary.Read(sstReader, binary.LittleEndian, &cbExtRst)
			if err != nil {
				return err
			}
		}

		//每个结构保存的数据
		stringbyte := make([]byte, stringlen)
		n, err = sstReader.Read(stringbyte[:])
		if err != nil || n != int(stringlen) {
			if err == nil {
				err = errors.New("read len is error")
			}
			return err
		}

		//unicode数据
		var str string
		if grbit&1 == 1 {
			ioread := bytes.NewReader(stringbyte)
			uintmp := make([]uint16, cchLen)
			err = binary.Read(ioread, binary.LittleEndian, &uintmp)
			if err != nil {
				return err
			}
			str = publicfun.UTF16ToString(uintmp[:])
		} else {
			//ansi数据
			str = string(stringbyte)
		}

		if len(str) > 0 {
			ole.xlsInfo.rgbStrings = append(ole.xlsInfo.rgbStrings, str)
		}

		allOffset := int64(cRun) + int64(cbExtRst)
		if allOffset != 0 {
			_, err = sstReader.Seek(allOffset, io.SeekCurrent)
			if err != nil {
				return err
			}
		}
	}
}
