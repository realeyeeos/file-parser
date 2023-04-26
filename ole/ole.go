package ole

/*
Date：2023.03.02
Author：scl
Description：解析ole文件
*/

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"strings"
)

//文件类型
const (
	NONE = iota
	WORDDOCUMENT
	XLS
	PPT
)

type XlsInfo struct {
	sheetNames []SHEETSTRUCT
	rgbStrings []string
	//版本
	xlsVersion []byte
}

type OleInfo struct {
	//文件头
	header Header
	//类型
	Type string
	//文件句柄
	file *os.File
	//文件属性
	fi os.FileInfo
	//fat索引
	difatPositions []uint32
	//minifat索引
	miniFatPositions []uint32
	//目录
	dirs []*Directory

	xlsInfo XlsInfo
}

type DataCallBackFunc func(string, string) bool

//判断是否是office97
func (ole *OleInfo) GetHandle(fp *os.File) error {
	var err error
	if fp == nil {
		err = errors.New("fp is nil")
		return err
	}

	//获取ole头
	err = ole.GetFileHeader(fp)
	if err != nil {
		return err
	}

	return nil
}

//打开文件
func (ole *OleInfo) OpenFile(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}

	err = ole.Read(f)
	if err != nil {
		return
	}
	return
}

//读取文件
func (ole *OleInfo) Read(fp *os.File) (err error) {
	if fp == nil {
		err = errors.New("fp is nil")
		return
	}

	//获取ole头
	err = ole.GetFileHeader(fp)
	if err != nil {
		return
	}

	//获取Fat结构
	err = ole.getFatSectors()
	if err != nil {
		return
	}

	//获取MiniFat结构
	err = ole.getMiniFATSectors()
	if err != nil {
		return
	}

	//获取所有Direcor结构
	err = ole.getDirectories()
	if err != nil {
		return
	}

	return
}

//获取结构
func (ole *OleInfo) GetObjectData(callBack DataCallBackFunc) (err error) {
	filetype := NONE

	var book *Directory
	var root *Directory
	var table0 *Directory
	var table1 *Directory
	for _, dir := range ole.dirs {
		fn := dir.getName()

		if strings.Contains(fn, "Root Entry") {
			root = dir
		}

		//xls
		if fn == "Book" {
			book = dir
			filetype = XLS
		}
		if fn == "Workbook" {
			if book == nil {
				book = dir
				filetype = XLS
			}
		}

		//doc
		if fn == "WordDocument" {
			book = dir
			filetype = WORDDOCUMENT
		}
		if fn == "0Table" {
			table0 = dir
		}
		if fn == "1Table" {
			table1 = dir
		}

		//ppt
		if fn == "PowerPoint Document" {
			book = dir
			filetype = PPT
		}

	}

	if book == nil {
		err = errors.New("book is nul")
		return
	}

	reader, err := ole.openObject(book, root)
	if err != nil {
		return
	}

	//获取doc数据
	if filetype == WORDDOCUMENT {
		//获取文档使用的table类型（1Table/0Table）
		table, err := ole.getTableReader(reader, table0, table1)
		if err != nil {
			return err
		}

		//获取table的数据流
		// table_reader, err := ole.openObject(table, root)
		// if err != nil {
		// 	return err
		// }

		//获取doc数据
		err = ole.getDocInfo(reader, table, root, callBack)
		if err != nil {
			return err
		}
		return nil
	} else if filetype == PPT {
		//获取ppt数据
		err = ole.getPptInfo(reader, book, callBack)
		if err != nil {
			return
		}
		return
	} else if filetype == XLS {
		err = ole.getXlsInfo(reader, book, callBack)
		if err != nil {
			return
		}
		return
	}

	return
}

//根据需要获取的长度获取数据
func (ole *OleInfo) getLenTableData(object, root *Directory, bgpos, datalen uint32) (data []byte, err error) {
	limitLen := int(datalen)
	//minifat
	if binary.LittleEndian.Uint32(object.StreamSize[:]) < binary.LittleEndian.Uint32(ole.header.MiniStreamCutoffSize[:]) {
		if root == nil {
			err = errors.New("root is nil")
			return nil, err
		}
		data = make([]byte, limitLen)

		offset := binary.LittleEndian.Uint32(object.StartSect[:])
		//fat数量
		number_of_fat_sectors := binary.LittleEndian.Uint32(ole.header.NumberFATSectors[:])
		//fat中sector数量（固定128）
		sectors_count := number_of_fat_sectors * ole.header.sectorSize() / 4

		current_mini_sector_index := uint32(0)
		//512字节中包括minisector个数（8）
		mini_sectors_in_sector := ole.header.sectorSize() / ole.header.minifatsectorSize()
		//root在fat的sector坐标
		mini_sector_location := binary.LittleEndian.Uint32(root.StartSect[:])

		index := uint32(0)
		isDataEmpty := true
		//实际上结构是[rootfat + 8个table（512）]
		for limitLen > 0 {
			//root
			sector_index := offset / mini_sectors_in_sector
			if sector_index != current_mini_sector_index {
				current_mini_sector_index = sector_index
				mini_sector_location = binary.LittleEndian.Uint32(root.StartSect[:])
				for sector_index > 0 {
					if mini_sector_location >= sectors_count {
						return
					}

					mini_sector_location = ole.difatPositions[mini_sector_location]
					sector_index--
				}
			}
			//sector偏移个数（minifat相对于fat）
			mini_sector_offset := offset - current_mini_sector_index*mini_sectors_in_sector
			//每组数据的偏移量（root在fat偏移量+table在minifat偏移量）
			point := (1+mini_sector_location)*ole.header.sectorSize() + mini_sector_offset*ole.header.minifatsectorSize()

			//根据偏移量获取数据
			sector := NewMiniFatSector(&ole.header)
			err := ole.getData(point, &sector.Data)
			if err != nil {
				return nil, err
			}

			if bgpos <= (index+1)*ole.header.minifatsectorSize() && bgpos >= index*ole.header.minifatsectorSize() {
				pos := bgpos - index*ole.header.minifatsectorSize()
				if limitLen+int(pos) > int(ole.header.minifatsectorSize()) {
					copy(data, sector.Data[pos:])
					limitLen -= int(ole.header.minifatsectorSize() - pos)
				} else {
					copy(data[int(datalen)-limitLen:], sector.Data[pos:int(pos)+limitLen])
					limitLen = 0
				}
				isDataEmpty = false
			} else {
				//是否已经存放数据
				if !isDataEmpty && limitLen > 0 {
					if limitLen > int(ole.header.minifatsectorSize()) {
						copy(data, sector.Data[:])
						limitLen -= int(ole.header.minifatsectorSize())
					} else {
						copy(data[int(datalen)-limitLen:], sector.Data[:limitLen])
						limitLen = 0
					}
				}
			}

			offset = ole.miniFatPositions[offset]
			if offset == binary.LittleEndian.Uint32(ENDOFCHAIN) || offset == binary.LittleEndian.Uint32(FATSECT) ||
				offset == binary.LittleEndian.Uint32(DIFSECT) || offset == binary.LittleEndian.Uint32(NOTAPP) ||
				offset == binary.LittleEndian.Uint32(MAXREGSECT) || offset == binary.LittleEndian.Uint32(FREESECT) {
				break
			}
			index++
		}
	} else { //fat
		data = make([]byte, limitLen)

		index := uint32(0)
		offset := binary.LittleEndian.Uint32(object.StartSect[:])
		isdataempty := true
		for limitLen > 0 {
			point := ole.sectorOffset(offset)

			sector := NewSector(&ole.header)
			err = ole.getData(point, &sector.Data)
			if err != nil {
				return nil, err
			}

			if bgpos <= (index+1)*ole.header.sectorSize() && bgpos >= index*ole.header.sectorSize() {
				pos := bgpos - index*ole.header.sectorSize()
				if limitLen+int(pos) > int(ole.header.sectorSize()) {
					copy(data, sector.Data[pos:])
					limitLen -= int(ole.header.sectorSize() - pos)
				} else {
					copy(data, sector.Data[pos:int(pos)+limitLen])
					limitLen = 0
				}
				isdataempty = false
			}

			//数据是否已经开始获取
			if !isdataempty {
				if limitLen > int(ole.header.sectorSize()) {
					copy(data, sector.Data[:])
					limitLen -= int(ole.header.sectorSize())
				} else {
					copy(data, sector.Data[:limitLen])
				}
			}

			offset = ole.difatPositions[offset]

			if offset == binary.LittleEndian.Uint32(ENDOFCHAIN) || offset == binary.LittleEndian.Uint32(FATSECT) ||
				offset == binary.LittleEndian.Uint32(DIFSECT) || offset == binary.LittleEndian.Uint32(NOTAPP) ||
				offset == binary.LittleEndian.Uint32(MAXREGSECT) || offset == binary.LittleEndian.Uint32(FREESECT) {
				break
			}
			index++
		}

	}
	return
}

//获取数据
func (ole *OleInfo) openObject(object, root *Directory) (reader io.ReadSeeker, err error) {
	if binary.LittleEndian.Uint32(object.StreamSize[:]) < binary.LittleEndian.Uint32(ole.header.MiniStreamCutoffSize[:]) {
		//return reader, nil
		if root == nil {
			err = errors.New("root is nil")
			return nil, err
		}
		data, err := ole.getDataFromMiniFat(binary.LittleEndian.Uint32(root.StartSect[:]), binary.LittleEndian.Uint32(object.StartSect[:]),
			binary.LittleEndian.Uint32(object.StreamSize[:]))

		if err != nil {
			return nil, err
		}

		reader = bytes.NewReader(data)
	} else {
		data, err := ole.getDataFromFatChain(binary.LittleEndian.Uint32(object.StartSect[:]), binary.LittleEndian.Uint32(object.StreamSize[:]))
		if err != nil {
			return nil, err
		}

		reader = bytes.NewReader(data)
	}

	return
}

//获取偏移量
func (ole *OleInfo) calculateOffset(sectorID []byte) (n uint32) {

	if len(sectorID) == 4 {
		n = binary.LittleEndian.Uint32(sectorID)
	}
	if len(sectorID) == 2 {
		n = uint32(binary.LittleEndian.Uint16(sectorID))
	}
	return (n * ole.header.sectorSize()) + ole.header.sectorSize()
}

//根据偏移量获取具体数据
func (ole *OleInfo) getData(offset uint32, data *[]byte) (err error) {

	_, err = ole.file.Seek(int64(offset), 0)
	if err != nil {

		return
	}

	_, err = ole.file.Read(*data)

	if err != nil {
		return
	}
	return
}

func BytesToUints16(b []byte) (res []uint16) {
	var section = make([]byte, 0)
	for _, value := range b {
		section = append(section, value)
		if len(section) == 2 {
			res = append(res, binary.LittleEndian.Uint16(section))

			section = make([]byte, 0)
		}
	}
	return
}
