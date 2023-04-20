package ole

/*
Date：2023.03.02
Author：scl
Description：解析doc文件
*/

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	publicfun "collector/publicfunc"
)

const MaxFileNum = 3000

// WordDocument 结构
type FIB struct {
	WIdent [2]byte
	// 193
	NFib   uint16
	Unused [2]byte
	Lid    [2]byte
	PnNext [2]byte
	// A-是否是dot文件 C-是否为复杂模式	D-是否包含图片
	AE byte
	// F-是否加密 G-存储表（0/1）
	FM        byte
	NFibBack  [2]byte
	LKey      [4]byte
	Envr      byte
	NS        byte
	Reserved3 [2]byte
	Reserved4 [2]byte
	Fcmin     uint32
	Reserved6 [4]byte

	Csw       [2]byte
	FlibRgW97 [28]byte
	Cslw      [2]byte

	CbMac       [4]byte
	Reserved1_2 [8]byte
	// 正文字数
	CcpText uint32
	// 页脚字数
	CcpFtn uint32
	// 页眉字数
	CcpHdd uint32
	// 批注字数
	CcpMcr uint32
	// 尾注字数
	CcpAtn uint32
	CcpEdn uint32
	// 文本框字数
	CcpTxbx uint32
	// 页眉文本框字数
	CcpHdrTxbx uint32

	Reserved4_14 [44]byte
}

var m_docfib FIB

type PCD struct {
	Fnordr [2]byte
	Fc     [4]byte
	Prm    [2]byte
}

// 获取对应的table
func (ole *OleInfo) getTableReader(reader io.ReadSeeker, table0 *Directory, table1 *Directory) (table *Directory, err error) {
	err = binary.Read(reader, binary.LittleEndian, &m_docfib)
	if err != nil {
		return
	}

	if binary.LittleEndian.Uint16(m_docfib.WIdent[:]) != 0xA5EC {
		err = errors.New("this file is not doc")
		return
	}

	// 判断引用的table
	var wtb [2]byte
	wtb[0] = m_docfib.FM >> 1 & 0x1
	fWhichTbStm := binary.LittleEndian.Uint16(wtb[:])
	if fWhichTbStm == 1 && table1 != nil {
		table = table1
	} else if fWhichTbStm == 0 && table0 != nil {
		table = table0
	} else {
		err = errors.New("table is null")
		return
	}

	return
}

// 获取doc信息
func (ole *OleInfo) getDocInfo(reader io.ReadSeeker, object, root *Directory, callBack DataCallBackFunc) (err error) {
	// 读取FIBTable97中的fcClx、lcbClx
	_, err = reader.Seek(32+2+28+2+88+2+264, 0)
	if err != nil {
		return
	}
	var fibclxdata [8]byte
	_, err = reader.Read(fibclxdata[:])
	if err != nil {
		return
	}

	// clx结构偏移量
	fcClx := binary.LittleEndian.Uint32(fibclxdata[:4])
	// clx结构大小
	lcbClx := binary.LittleEndian.Uint32(fibclxdata[4:])

	if fcClx > uint32(ole.fi.Size()) || fcClx+lcbClx > uint32(ole.fi.Size()) {
		err = errors.New("fcclx and size is error")
		return
	}

	// 获取clx结构数据
	clxdata, err := ole.getClxData(fcClx, lcbClx, object, root)
	if err != nil {
		return
	}

	// pcdt结构
	if clxdata[0] == 0x02 {
		// 获取pcdt结构及数据
		err = ole.getPcdtStruct(reader, clxdata, callBack)
		clxdata = nil
		return
	} else {
		// TODO PRC结构 by scl
		// 0x01	PRC结构
	}

	clxdata = nil
	return
}

// 获取clx结构数据
func (ole *OleInfo) getClxData(fcClx, lcbClx uint32, object, root *Directory) (clxdata []byte, err error) {
	if lcbClx > 0 {
		// 从table流和fcclx偏移量找到clx结构位置（即从table的data中找到clx）

		clxtype, err := ole.getLenTableData(object, root, fcClx, 1)

		if err != nil {
			return nil, err
		}

		if clxtype[0] != 0x02 {
			err = errors.New("clxtype is error")
			return nil, err
		}

		clxdata, err = ole.getLenTableData(object, root, fcClx, lcbClx)
		if err != nil {
			clxdata = nil
			return nil, err
		}
	} else {
		// 大小为0的时候手动计算
		clxdata = make([]byte, 21)
		clxdata[0] = 0x02
		clxdata[1] = 0x10

		acp := m_docfib.CcpText + m_docfib.CcpFtn + m_docfib.CcpHdd + m_docfib.CcpMcr +
			m_docfib.CcpAtn + m_docfib.CcpEdn + m_docfib.CcpTxbx + m_docfib.CcpHdrTxbx
		binary.LittleEndian.PutUint32(clxdata[9:13], acp)

		fcMin := m_docfib.Fcmin << 1
		fcMin |= 0x40000000
		binary.LittleEndian.PutUint32(clxdata[15:19], fcMin)

	}
	return
}

// 获取Pcdt结构数据
func (ole *OleInfo) getPcdtStruct(reader io.ReadSeeker, clxdata []byte, callBack DataCallBackFunc) (err error) {
	//TODO
	if len(clxdata) < 5 {
		return
	}
	lcd := binary.LittleEndian.Uint32(clxdata[1:5])

	acpnums, apcds, err := ole.getPlcPcdInfo(clxdata[5:], lcd)
	if err != nil {
		return
	}

	// strs = make([]string, len(apcds))
	index := 0
	for i := 0; i < len(apcds); i++ {
		// var apdc [4]byte
		var fcCompressed uint32
		binary.Read(bytes.NewBuffer(apcds[i].Fc[:]), binary.LittleEndian, &fcCompressed)
		// 前30位是fc数据
		fc := fcCompressed & 0x3FFFFFFF

		// 第31位是fCompressed数据
		var isunicode bool
		if (fcCompressed & 0x40000000) == 0 {
			isunicode = true
		}

		// Unicode数据
		fcpoint := fc
		dataend_point := acpnums[i+1] * 2
		// Assi数据
		if !isunicode {
			fcpoint = fc / 2
			dataend_point /= 2
		}

		// 根据偏移量获取数据
		reader.Seek(int64(fcpoint), 0)

		if isunicode {
			datalen := int(dataend_point)
			for datalen > 0 {
				filenum := MaxFileNum * 2
				if datalen < MaxFileNum*2 {
					filenum = datalen
				}
				worddata := make([]byte, filenum)
				_, err := reader.Read(worddata)
				if err != nil {
					datalen -= MaxFileNum * 2
					continue
				}

				// unicode字符解析
				uintmp := make([]uint16, filenum/2)
				err = binary.Read(bytes.NewReader(worddata), binary.LittleEndian, &uintmp)
				if err != nil {
					continue
				}
				str2 := publicfun.UTF16ToString(uintmp[:])

				if !callBack(str2, "") {
					return nil
				}

				datalen -= filenum
			}

		} else {
			datalen := int(dataend_point)
			for datalen > 0 {
				filenum := MaxFileNum
				if datalen < MaxFileNum {
					filenum = datalen
				}
				worddata := make([]byte, filenum)
				_, err := reader.Read(worddata)
				if err != nil {
					datalen -= MaxFileNum
					continue
				}
				if !callBack(string(worddata), "") {
					return nil
				}

				datalen -= MaxFileNum
			}
		}
		index++
	}

	return
}

// 获取PlcPcd数据
func (ole *OleInfo) getPlcPcdInfo(data []byte, lcd uint32) (acpnums []uint32, apcds []PCD, err error) {
	pcd_num := (lcd - 4) / 12

	plcpcd_reader := bytes.NewReader(data)
	for i := uint32(0); i <= pcd_num; i++ {
		var acp [4]byte
		_, err = plcpcd_reader.Read(acp[:])
		if err != nil {
			break
		}

		acpnums = append(acpnums, binary.LittleEndian.Uint32(acp[:]))
	}

	for i := uint32(0); i < pcd_num; i++ {
		var pcd PCD
		err = binary.Read(plcpcd_reader, binary.BigEndian, &pcd)
		if err != nil {
			break
		}

		apcds = append(apcds, pcd)
	}

	return
}
