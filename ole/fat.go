package ole

/*
Date：2023.03.02
Author：scl
Description：解析fat结构
*/

import (
	"bytes"
	"encoding/binary"
	"errors"
)

//每个fat大小
var EntrySize = 128
var DefaultDIFATEntries = uint32(109)

//获取Fat信息
func (ole *OleInfo) getFatSectors() (err error) {

	entries := DefaultDIFATEntries

	if binary.LittleEndian.Uint32(ole.header.NumberFATSectors[:]) < DefaultDIFATEntries {
		entries = binary.LittleEndian.Uint32(ole.header.NumberFATSectors[:])
	}

	//获取所有fat数据
	for i := uint32(0); i < entries; i++ {
		//4个字节一组
		position := ole.calculateOffset(ole.header.getDIFATEntry(i))
		sector := NewSector(&ole.header)

		err := ole.getData(position, &sector.Data)

		if err != nil {
			return err
		}

		ole.difatPositions = append(ole.difatPositions, sector.values(EntrySize)...)
	}

	if bytes.Equal(ole.header.DifatSect[:], ENDOFCHAIN) || bytes.Equal(ole.header.DifatSect[:], FATSECT) ||
		bytes.Equal(ole.header.DifatSect[:], DIFSECT) || bytes.Equal(ole.header.DifatSect[:], NOTAPP) ||
		bytes.Equal(ole.header.DifatSect[:], MAXREGSECT) || bytes.Equal(ole.header.DifatSect[:], FREESECT) {
		return
	}

	//后边基本用不到，如果文件很大可能会有
	// position := ole.calculateOffset(ole.header.DifatSect[:])
	// var section = make([]byte, 0)
	// for i := uint32(0); i < binary.LittleEndian.Uint32(ole.header.NumberDIFATSectors[:]); i++ {
	// 	sector := NewSector(&ole.header)
	// 	err := ole.getData(position, &sector.Data)

	// 	if err != nil {
	// 		return err
	// 	}

	// 	for _, value := range sector.getFATSectorLocations() {
	// 		section = append(section, value)
	// 		if len(section) == 4 {
	// 			position = ole.calculateOffset(section)
	// 			sectorF := NewSector(&ole.header)
	// 			err := ole.getData(position, &sectorF.Data)

	// 			if err != nil {
	// 				return err
	// 			}
	// 			ole.difatPositions = append(ole.difatPositions, sectorF.values(EntrySize)...)

	// 			section = make([]byte, 0)
	// 		}

	// 	}

	// 	position = ole.calculateOffset(sector.getNextDIFATSectorLocation())

	// }

	err = errors.New("read fat error")
	return
}
