package ole

/*
Date：2023.03.02
Author：scl
Description：解析xls中sheet数据
*/

import (
	"bytes"
	"encoding/binary"
	publicfun "fileparser/publicfunc"
	"io"
	"math"
	"strconv"
)

const (
	//sheet
	XLS_LABEL_SST = 0xFD
	XLS_LABEL     = 0x204
	XLS_BLANK     = 0x201
	XLS_BOOLERR   = 0x205
	XLS_NUMBER    = 0x203
	XLS_MULBLANK  = 0xBE
	XLS_RK        = 0x27E
	XLS_MULRK     = 0xBD
)

//sheet结构（名字+偏移量）
type SHEETSTRUCT struct {
	name  string
	lpPos int64
}

//获取sheet相关信息（包括名字和偏移量）
func (ole *OleInfo) getSheetInfo(data []byte) (err error) {
	var sheetstruct SHEETSTRUCT

	sheetreader := bytes.NewReader(data[:])

	var lbPlyPos [4]byte
	_, err = sheetreader.Read(lbPlyPos[:])
	if err != nil {
		return
	}

	//每个sheet在总数据中的偏移量
	sheetstruct.lpPos = int64(binary.LittleEndian.Uint32(lbPlyPos[:]))
	_, err = sheetreader.Seek(2, io.SeekCurrent)
	if err != nil {
		return
	}

	//字符串数量
	var cchbyte, code [1]byte
	_, err = sheetreader.Read(cchbyte[:])
	if err != nil {
		return
	}

	//0-ansi	1-unicode
	_, err = sheetreader.Read(code[:])
	if err != nil {
		return
	}

	//字符串占用字节数
	var namelen uint
	namelen = uint(cchbyte[0])
	if code[0] == 1 {
		namelen *= 2
	}

	sheetname_byte := make([]byte, namelen)
	_, err = sheetreader.Read(sheetname_byte[:])
	if err != nil {
		return
	}

	var sheetname string
	if code[0] == 1 {
		ioread := bytes.NewReader(sheetname_byte)
		uintmp := make([]uint16, cchbyte[0])
		binary.Read(ioread, binary.LittleEndian, &uintmp)
		sheetname = publicfun.UTF16ToString(uintmp[:])
	} else {
		sheetname = string(sheetname_byte)
	}

	//保存sheet名和其偏移量
	sheetstruct.name = sheetname
	ole.xlsInfo.sheetNames = append(ole.xlsInfo.sheetNames, sheetstruct)

	return
}

//处理sheet流获取每个sheet的数据
func (ole *OleInfo) dealSheet(sheetReader io.ReadSeeker, sheetStruct SHEETSTRUCT, callBack DataCallBackFunc) error {
	//sheetreader := bytes.NewReader(sheetdata[:])
	_, err := sheetReader.Seek(sheetStruct.lpPos, io.SeekStart)
	if err != nil {
		return err
	}

	//var str string
	for {
		var rec_type, rec_len [2]byte
		//类型
		_, err = sheetReader.Read(rec_type[:])
		if err != nil {
			return err
		}

		//数据长度
		_, err = sheetReader.Read(rec_len[:])
		if err != nil {
			return err
		}

		urec_type := binary.LittleEndian.Uint16(rec_type[:])
		urec_len := binary.LittleEndian.Uint16(rec_len[:])

		//结束
		if urec_type == XLS_EOF {
			break
		}

		if urec_len == 0 {
			continue
		}

		if urec_type != XLS_LABEL_SST && urec_type != XLS_LABEL && urec_type != XLS_BOOLERR && urec_type != XLS_NUMBER &&
			urec_type != XLS_RK {
			sheetReader.Seek(int64(urec_len), io.SeekCurrent)
			continue
		}

		//当前流的数据
		data := make([]byte, urec_len)
		_, err = sheetReader.Read(data[:])
		if err != nil {
			data = nil
			return err
		}

		datareader := bytes.NewReader(data[:])

		//字符串
		if urec_type == XLS_LABEL_SST {

			type LabelSSt struct {
				Rw   uint16
				Col  uint16
				Ixfe uint16
				Isst uint32
			}

			var lablesst LabelSSt
			err = binary.Read(datareader, binary.LittleEndian, &lablesst)
			if err != nil {
				data = nil
				return err
			}

			if len(ole.xlsInfo.rgbStrings) >= int(lablesst.Isst)+1 {
				if !callBack(ole.xlsInfo.rgbStrings[lablesst.Isst], "工作表（"+sheetStruct.name+"）：第"+strconv.Itoa(int(lablesst.Rw+1))+"行，第"+strconv.Itoa(int(lablesst.Col+1))+"列") {
					return nil
				}
				//str += rgbstrings[lablesst.Isst] + "\t"
				//fmt.Println(lablesst.Rw, ",", lablesst.Col, ":", rgbstrings[lablesst.Isst])
			}
		} else if urec_type == XLS_LABEL {
			//字符串
			type LabelBIFF8 struct {
				Rw    uint16
				Col   uint16
				Ixfe  uint16
				Cch   uint16
				Grbit [2]byte
			}

			var lablbif8 LabelBIFF8
			err = binary.Read(datareader, binary.LittleEndian, &lablbif8)
			if err != nil {
				data = nil
				return err
			}
			if lablbif8.Cch == 0 {
				data = nil
				continue
			}
			//LabelBIFF8
			labledata := make([]byte, lablbif8.Cch)

			if bytes.Equal(ole.xlsInfo.xlsVersion, []byte{0x00, 0x06}) && lablbif8.Grbit[0] == 1 {
				labledata = make([]byte, lablbif8.Cch*2)
			}

			_, err = datareader.Read(labledata)
			if err != nil {
				labledata = nil
				data = nil
				return err
			}

			if bytes.Equal(ole.xlsInfo.xlsVersion, []byte{0x00, 0x06}) && lablbif8.Grbit[0] == 1 {
				ioread := bytes.NewReader(labledata)
				uintmp := make([]uint16, lablbif8.Cch)
				binary.Read(ioread, binary.LittleEndian, &uintmp)

				if !callBack(publicfun.UTF16ToString(uintmp[:]), "工作表（"+sheetStruct.name+"）：第"+strconv.Itoa(int(lablbif8.Rw+1))+"行，第"+strconv.Itoa(int(lablbif8.Col+1))+"列") {
					return nil
				}
				//str += publicfun.UTF16ToString(uintmp[:]) + "\t"
			} else {
				if !callBack(string(labledata), "工作表（"+sheetStruct.name+"）：第"+strconv.Itoa(int(lablbif8.Rw+1))+"行，第"+strconv.Itoa(int(lablbif8.Col+1))+"列") {
					return nil
				}
				//str += string(labledata) + "\t"
			}

			//一个空的列
			//} else if urec_type == XLS_BLANK {

			// type Blank struct {
			// 	rw   [2]byte
			// 	col  [2]byte
			// 	ixfe [2]byte
			// }

		} else if urec_type == XLS_BOOLERR {
			//bool类型数据
			type BoolErr struct {
				Rw       uint16
				Col      uint16
				Ixfe     uint16
				BBoolErr uint
				FError   uint
			}

			var boolerr BoolErr
			err = binary.Read(datareader, binary.LittleEndian, &boolerr)
			if err != nil {
				data = nil
				return err
			}

			var errorstr string
			if boolerr.FError == 1 {
				switch boolerr.BBoolErr {
				case 0:
					errorstr = "#NULL!"
				case 7:
					errorstr = "#DIV/0!"
				case 15:
					errorstr = "#VALUE!"
				case 23:
					errorstr = "#REF!"
				case 29:
					errorstr = "#NAME?"
				case 36:
					errorstr = "#NUM!!"
				case 42:
					errorstr = "#N/A"
				}
			} else {
				if boolerr.BBoolErr == 1 {
					errorstr = "TRUE"
				} else {
					errorstr = "FALSE"
				}
			}
			if !callBack(errorstr, "工作表（"+sheetStruct.name+"）：第"+strconv.Itoa(int(boolerr.Rw+1))+"行，第"+strconv.Itoa(int(boolerr.Col+1))+"列") {
				return nil
			}
			//str += errorstr

		} else if urec_type == XLS_NUMBER {
			//数字
			type Number struct {
				Rw   uint16
				Col  uint16
				Ixfe uint16
				Num  uint64
			}

			var numberinfo Number
			err = binary.Read(datareader, binary.LittleEndian, &numberinfo)
			if err != nil {
				data = nil
				return err
			}

			if !callBack(strconv.Itoa(int(numberinfo.Num)), "工作表（"+sheetStruct.name+"）：第"+strconv.Itoa(int(numberinfo.Rw+1))+"行，第"+strconv.Itoa(int(numberinfo.Col+1))+"列") {
				return nil
			}
			//str += strconv.Itoa(int(numberinfo.Num)) + "\t"

			// } else if urec_type == XLS_MULBLANK {

		} else if urec_type == XLS_RK {
			//数字
			type Rk struct {
				Rw   uint16
				Col  uint16
				Ixfe uint16
				Rk   uint32
			}

			var rk Rk
			err = binary.Read(datareader, binary.LittleEndian, &rk)
			if err != nil {
				return err
			}

			if rk.Rk == 0 {
				return nil
			}
			inum, fnum, isFloat := ole.getrk_number(rk.Rk)
			if isFloat {
				if !callBack(strconv.FormatFloat(fnum, 'f', -1, 64), "工作表（"+sheetStruct.name+"）：第"+strconv.Itoa(int(rk.Rw+1))+"行，第"+strconv.Itoa(int(rk.Col+1))+"列") {
					return nil
				}
				//str += strconv.FormatFloat(fnum, 'f', -1, 64) + "\t"
			} else {
				if !callBack(strconv.FormatInt(inum, 10), "工作表（"+sheetStruct.name+"）：第"+strconv.Itoa(int(rk.Rw+1))+"行，第"+strconv.Itoa(int(rk.Col+1))+"列") {
					return nil
				}
				//str += strconv.FormatInt(inum, 10) + "\t"

			}

			// } else if urec_type == XLS_MULRK {

		}
		data = nil
	}
	// if len(str) > 0 {
	// 	callback(str)
	// 	str = ""
	// }

	return nil
}

//获取保存的数据
func (ole *OleInfo) getrk_number(rk uint32) (intNum int64, floatNum float64, isFloat bool) {
	val := uint64(rk >> 2)
	rkType := uint(rk << 30 >> 30)

	var fn float64
	switch rkType {
	case 0:
		fn = math.Float64frombits(uint64(rk&0xfffffffc) << 32)
		isFloat = true
	case 1:
		fn = math.Float64frombits(uint64(rk&0xfffffffc)<<32) / 100
		isFloat = true
	case 3:
		fn = float64(val) / 100
		isFloat = true
	}

	return int64(val), float64(fn), isFloat
}
