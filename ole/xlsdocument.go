package ole

/*
Date：2023.03.02
Author：scl
Description：解析xls文件
*/

import (
	"bytes"
	"encoding/binary"
	publicfun "fileparser/publicfunc"
	"io"
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
	var SSTContinue bool
	//获取sheet名字、位置 + SST中的数据
	for {
		var rec_type, rec_len [2]byte
		//类型
		_, err = reader.Read(rec_type[:])
		if err != nil {
			return
		}

		//数据长度
		_, err = reader.Read(rec_len[:])
		if err != nil {
			return
		}

		//关键值
		urec_type := binary.LittleEndian.Uint16(rec_type[:])
		//大小
		urec_len := binary.LittleEndian.Uint16(rec_len[:])

		//结束符
		if urec_type == XLS_EOF {
			break
		}

		if urec_len == 0 {
			continue
		}
		data := make([]byte, urec_len)
		_, err = reader.Read(data[:])
		if err != nil {
			return
		}

		//开始符
		if urec_type == XLS_BOF {
			if urec_len < 2 {
				return
			}
			ole.xlsInfo.xlsVersion = data[0:2]
		} else if urec_type == XLS_SHEET {
			//获取sheet信息（包括名字和偏移量）
			ole.getSheetInfo(data)
		} else if urec_type == XLS_CONTINUE {
			if SSTContinue {

				// 	if len(wb.sst.RgbSrc) == 0  {
				// 		grbitOffset = 0
				// 	} else {
				// 		grbitOffset = 1
				// 	}

				// 	grbit = stream[sPoint]

				// 	wb.sst.RgbSrc = append(wb.sst.RgbSrc, stream[sPoint+grbitOffset:sPoint+recordDataLength]...)
				// 	wb.sst.Read(readType, grbit, prevLen)
			}

		} else if urec_type == XLS_SST {
			//获取SST数据
			ucstTotal, _ := ole.getSSTInfo(urec_type, data)

			if urec_len >= 8224 || len(ole.xlsInfo.rgbStrings) < int(ucstTotal-1) {
				SSTContinue = true
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
func (ole *OleInfo) getSSTInfo(urectype uint16, data []byte) (ucstTotal uint32, err error) {
	sstreader := bytes.NewReader(data[:])
	var cstTotal, cstUnique [4]byte
	_, err = sstreader.Read(cstTotal[:])
	if err != nil {
		return
	}
	//总共数量
	ucstTotal = binary.LittleEndian.Uint32(cstTotal[:])

	_, err = sstreader.Read(cstUnique[:])
	if err != nil {
		return
	}

	//rgbstrings = nil
	for {
		//字符数量
		var cchbyte [2]byte
		_, err = sstreader.Read(cchbyte[:])
		if err != nil {
			return
		}

		var flags [1]byte
		_, err = sstreader.Read(flags[:])
		if err != nil {
			return
		}

		//0-ansi	1-unicode
		grbit := flags[:1][0]

		cch := binary.LittleEndian.Uint16(cchbyte[:])
		if cch == 0 {
			continue
		}
		stringlen := cch
		if grbit&1 == 1 {
			stringlen = cch * 2
		}

		//每个结构保存的数据
		stringbyte := make([]byte, stringlen)
		_, err = sstreader.Read(stringbyte[:])
		if err != nil {
			return
		}

		//unicode数据
		var str string
		if grbit&1 == 1 {
			ioread := bytes.NewReader(stringbyte)
			uintmp := make([]uint16, cch)
			binary.Read(ioread, binary.LittleEndian, &uintmp)
			str = publicfun.UTF16ToString(uintmp[:])
		} else {
			//ansi数据
			str = string(stringbyte)
		}

		if len(str) > 0 {
			ole.xlsInfo.rgbStrings = append(ole.xlsInfo.rgbStrings, str)

		}
	}
}
