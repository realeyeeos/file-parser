package ole

/*
Date：2023.03.02
Author：scl
Description：ole文件中的大小等
*/

import (
	"bytes"
	"encoding/binary"
)

var FREESECT = []byte{0xFF, 0xFF, 0xFF, 0xFF}

var ENDOFCHAIN = []byte{0xFE, 0xFF, 0xFF, 0xFF}

var FATSECT = []byte{0xFD, 0xFF, 0xFF, 0xFF}

var DIFSECT = []byte{0xFC, 0xFF, 0xFF, 0xFF}

var NOTAPP = []byte{0xFB, 0xFF, 0xFF, 0xFF}

var MAXREGSECT = []byte{0xFA, 0xFF, 0xFF, 0xFF}

type Sector struct {
	SectorSize uint32
	Data       []byte
}

func (s *Sector) getFATSectorLocations() []byte {
	return s.Data[0 : s.SectorSize-4]
}

func (s *Sector) getMiniFatFATSectorLocations() []byte {
	return s.Data[0:s.SectorSize]
}

func (s *Sector) getNextDIFATSectorLocation() []byte {
	return s.Data[s.SectorSize-4:]
}

// FAT
func NewSector(header *Header) Sector {
	return Sector{
		SectorSize: header.sectorSize(),
		Data:       make([]byte, header.sectorSize()),
	}

}

//MiniFat
func NewMiniFatSector(header *Header) Sector {
	return Sector{
		SectorSize: header.minifatsectorSize(),
		Data:       make([]byte, header.minifatsectorSize()),
	}
}

//读取固定大小数据
func (s *Sector) values(length int) (res []uint32) {
	res = make([]uint32, length)
	buf := bytes.NewBuffer(s.Data)

	_ = binary.Read(buf, binary.LittleEndian, res)

	return res
}
