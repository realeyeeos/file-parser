package ole

/*
Date：2023.03.02
Author：scl
Description：解析数据头
*/

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

//ole复合文件标识
var Olesign = []byte{0xd0, 0xcf, 0x11, 0xe0, 0xa1, 0xb1, 0x1a, 0xe1}

//主要版本为3，扇区大小为512
var MajorVersion3 = []byte{0x03, 0x00}

//主要版本为4，扇区大小为4096
var MajorVersion4 = []byte{0x04, 0x00}

//文件头 512
type Header struct {
	HeaderSignature            [8]byte
	HeaderCLSID                [16]byte
	MinorVersion               [2]byte
	MajorVersion               [2]byte
	ByteOrder                  [2]byte
	SectorShift                uint16
	MiniSectorShift            uint16
	Reserved                   [6]byte
	NumberDirectorySectors     [4]byte
	NumberFATSectors           [4]byte
	DirSect                    [4]byte
	TransactionSignatureNumber [4]byte
	MiniStreamCutoffSize       [4]byte
	MiniFatSect                [4]byte
	NumberMiniFATSectors       [4]byte
	DifatSect                  [4]byte
	NumberDIFATSectors         [4]byte
	DIFAT                      [436]byte
}

//读取文件头
func (ole *OleInfo) GetFileHeader(fileReadSeeker io.ReadSeeker) error {
	if ole.fileReadSeeker == nil {
		ole.fileReadSeeker = fileReadSeeker
	}
	_, err := fileReadSeeker.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	var headerdata [512]byte
	n, err := fileReadSeeker.Read(headerdata[:])
	if err != nil || n < 512 {
		errnew := errors.New("head read error")
		return errnew
	}

	err = binary.Read(bytes.NewBuffer(headerdata[:]), binary.LittleEndian, &ole.header)
	if err != nil {
		return err
	}

	//验证是否是ole复合结构文件
	if !bytes.Equal(ole.header.HeaderSignature[:], Olesign) {
		errnew := errors.New("head sign is error")
		return errnew
	}

	return nil
}

//fat大小
func (h *Header) sectorSize() (size uint32) {
	if h.SectorShift < 1 {
		return 512
	}

	size = 1 << h.SectorShift

	return size
}

//minifat大小
func (h *Header) minifatsectorSize() (size uint32) {
	if h.MiniSectorShift < 1 {
		return 64
	}

	size = 1 << h.MiniSectorShift
	return size
}

func (h *Header) getDIFATEntry(i uint32) []byte {
	return h.DIFAT[i*4 : (i*4)+4]
}
