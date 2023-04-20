package ole

/*
Date：2023.03.02
Author：scl
Description：解析ole文件中每个目录
*/
import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"
)

//128字节
type Directory struct {
	EleName        [64]byte
	CbEleNameLenth [2]byte
	ObjectType     byte
	ColorFlag      byte
	LeftSiblingID  [4]byte
	RightSiblingID [4]byte
	ChildID        [4]byte
	CLSID          [16]byte
	UserFlags      [4]byte
	CreationTime   [8]byte
	ModifiedTime   [8]byte
	StartSect      [4]byte
	StreamSize     [8]byte
}

//获取direct数据
func (ole *OleInfo) getDirectories() (err error) {
	stream, err := ole.getDataFromFatChain(binary.LittleEndian.Uint32(ole.header.DirSect[:]), 0)

	if err != nil {
		return err
	}
	var section = make([]byte, 0)

	for _, value := range stream {
		section = append(section, value)
		//128个字节一个数据
		if len(section) == EntrySize {
			var dir Directory
			err = binary.Read(bytes.NewBuffer(section), binary.LittleEndian, &dir)
			if err == nil && (dir.ObjectType == 0x00 || dir.ObjectType == 0x01 || dir.ObjectType == 0x02 || dir.ObjectType == 0x05) &&
				(dir.ColorFlag == 0x00 || dir.ColorFlag == 0x01) {
				ole.dirs = append(ole.dirs, &dir)
			}

			section = make([]byte, 0)
		}
	}

	return
}

//获取FAT所有direct中data字节（一般数据量大于4096字节时调用此函数）
func (ole *OleInfo) getDataFromFatChain(offset uint32, size uint32) (data []byte, err error) {
	if size > 0 {
		datalen := size / ole.header.sectorSize()
		if size%ole.header.sectorSize() > 0 {
			datalen += 1
		}
		data = make([]byte, datalen*ole.header.sectorSize())
	}

	index := uint32(0)
	for {
		sector := NewSector(&ole.header)
		point := ole.sectorOffset(offset)

		err = ole.getData(point, &sector.Data)

		if err != nil {
			return data, err
		}

		if size > 0 {
			if index*ole.header.sectorSize() > uint32(len(data)) || (index+1)*ole.header.sectorSize() > uint32(len(data)) {
				break
			}

			copy(data[index*ole.header.sectorSize():], sector.Data[:])
			index++
		} else {
			data = append(data, sector.Data...)
		}

		offset = ole.difatPositions[offset]

		if offset == binary.LittleEndian.Uint32(ENDOFCHAIN) || offset == binary.LittleEndian.Uint32(FATSECT) ||
			offset == binary.LittleEndian.Uint32(DIFSECT) || offset == binary.LittleEndian.Uint32(NOTAPP) ||
			offset == binary.LittleEndian.Uint32(MAXREGSECT) || offset == binary.LittleEndian.Uint32(FREESECT) {
			break
		}
	}

	return data, err
}

//获取MiniFAT所有direct中data字节（一般数据量小于4096字节时调用此函数）
func (ole *OleInfo) getDataFromMiniFat(rootFatSectorLocation uint32, offset uint32, size uint32) (data []byte, err error) {
	//fat数量
	number_of_fat_sectors := binary.LittleEndian.Uint32(ole.header.NumberFATSectors[:])
	//fat中sector数量（固定128）
	sectors_count := number_of_fat_sectors * ole.header.sectorSize() / 4

	current_mini_sector_index := uint32(0)
	//512字节中包括minisector个数（8）
	mini_sectors_in_sector := ole.header.sectorSize() / ole.header.minifatsectorSize()
	//root在fat的sector坐标
	mini_sector_location := rootFatSectorLocation

	if size > 0 {
		datalen := size / ole.header.minifatsectorSize()
		if size%ole.header.minifatsectorSize() > 0 {
			datalen += 1
		}

		data = make([]byte, datalen*ole.header.minifatsectorSize())
	}

	index := uint32(0)
	//实际上结构是[rootfat + 8个table（512）]
	for {
		//root
		sector_index := offset / mini_sectors_in_sector
		if sector_index != current_mini_sector_index {
			current_mini_sector_index = sector_index
			mini_sector_location = rootFatSectorLocation
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
			return data, err
		}

		if size > 0 {
			if index*ole.header.minifatsectorSize() > uint32(len(data)) || (index+1)*ole.header.minifatsectorSize() > uint32(len(data)) {
				break
			}

			copy(data[index*ole.header.minifatsectorSize():], sector.Data[:])
			index++
		} else {
			data = append(data, sector.Data...)
		}

		offset = ole.miniFatPositions[offset]
		if offset == binary.LittleEndian.Uint32(ENDOFCHAIN) || offset == binary.LittleEndian.Uint32(FATSECT) ||
			offset == binary.LittleEndian.Uint32(DIFSECT) || offset == binary.LittleEndian.Uint32(NOTAPP) ||
			offset == binary.LittleEndian.Uint32(MAXREGSECT) || offset == binary.LittleEndian.Uint32(FREESECT) {
			break
		}
	}
	return
}

//获取direct偏移量
func (ole *OleInfo) sectorOffset(sid uint32) uint32 {
	return (sid + 1) * ole.header.sectorSize()
}

func (ole *OleInfo) calculateMiniFatOffset(sid uint32) (n uint32) {

	return sid * ole.header.minifatsectorSize()
}

//获取目录名
func (dir *Directory) getName() (name string) {
	size := binary.LittleEndian.Uint16(dir.CbEleNameLenth[:])
	if size > 0 {
		size = size - 1
	} else if size == 0 {
		size = 32
	}

	bytename := BytesToUints16(dir.EleName[:size])
	runes := utf16.Decode(bytename)
	return string(runes)
}
