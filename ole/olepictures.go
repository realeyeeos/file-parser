package ole

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

/*
Date：2023.05.16
Author：scl
Description：解析doc文件中的图片（嵌入/浮动）
*/

//解析图片（OfficeArtBStoreContainerFileBlock）结构数据
func (ole *OleInfo) getArtData(documentReader, reader io.ReadSeeker, isInline bool) (err error) {
	//PICFAndOfficeArtData结构
	if isInline {
		var picbyte [2]byte
		reader.Seek(6, io.SeekCurrent)
		reader.Read(picbyte[:])
		if binary.LittleEndian.Uint16(picbyte[:]) == 0x66 {
			reader.Seek(68, io.SeekCurrent)
			var tmpbyte [2]byte
			reader.Read(tmpbyte[0:1])
			reader.Seek(int64(binary.LittleEndian.Uint16(tmpbyte[:])), io.SeekCurrent)
		} else {
			reader.Seek(68, io.SeekCurrent)
		}
	}

	for {
		//============header
		var artRecordHeader [8]byte
		n, err := reader.Read(artRecordHeader[:])
		if err != nil || n != 8 {
			if err == nil {
				err = errors.New("read len is error")
			}
			break
		}

		//rh
		recInt := binary.LittleEndian.Uint16(artRecordHeader[:2])
		recVer := recInt & 0xF
		recInstance := recInt >> 4
		recType := binary.LittleEndian.Uint16(artRecordHeader[2:4])
		recLen := binary.LittleEndian.Uint32(artRecordHeader[4:])

		if recInt == 0 && recType == 0 && recLen == 0 {
			break
		}

		//fmt.Println(recVer, recInstance, recType, recLen)

		//OfficeArtDggContainer数据
		switch recType {
		case 0xF000:
			artData := make([]byte, recLen)
			n, err = reader.Read(artData[:])
			if err != nil || n != int(recLen) {
				if err == nil {
					err = errors.New("read len is error")
				}
				break
			}
			ole.getArtDggContainerData(documentReader, artData[:])
		case 0xF007: //OfficeArtBStoreContainerFileBlock(内联图片)
			if recVer != 0x2 {
				continue
			}
			artData := make([]byte, recLen)
			n, err = reader.Read(artData[:])
			if err != nil || n != int(recLen) {
				if err == nil {
					err = errors.New("read len is error")
				}
				break
			}

			//获取图片数据
			imageData, err := ole.getArtFBSEData(documentReader, artData[:])
			if err != nil {
				break
			}
			pathName := "F:\\project_git\\dsp-fileplugin\\tmpfile\\scl\\123"
			switch recInstance {
			case 0x02:
				pathName += ".emf"
				//("emf")
			case 0x03:
				pathName += ".wmf"
				//fmt.Println("wmf")
			case 0x04:
				pathName += ".pict"
				//fmt.Println("pict")
			case 0x05, 0x12:
				pathName += ".jpeg"
				//fmt.Println("jpeg")
			case 0x06:
				pathName += ".png"
				//fmt.Println("png")
			case 0x07:
				pathName += ".dib"
				//fmt.Println("dib")
			case 0x011:
				pathName += ".tiff"
				//fmt.Println("tiff")
			default:
				break
			}

			os.WriteFile(pathName, imageData[:], 0666)
		default: //OfficeArtBlip
			if recType >= 0xF018 && recType <= 0xF117 {
				artData := make([]byte, recLen)
				n, err = reader.Read(artData[:])
				if err != nil || n != int(recLen) {
					if err == nil {
						err = errors.New("read len is error")
					}
					break
				}
				imageData, err := ole.getImageByType(artData[:], recInstance, recType)
				if err != nil {
					break
				}

				pathName := "F:\\project_git\\dsp-fileplugin\\tmpfile\\scl\\123"
				switch recType {
				case 0xF01A: //emf
					pathName += ".emf"
				case 0xF01B: //wmf
					pathName += ".wmf"
				case 0xF01C: //pict
					pathName += ".pict"
				case 0xF01D, 0xF02A: //jpeg
					pathName += ".jpeg"
				case 0xF01E: //png
					pathName += ".png"
				case 0xF01F: //dib
					pathName += ".dib"
				case 0xF029: //tiff
					pathName += ".tiff"
				}

				os.WriteFile(pathName, imageData[:], 0666)
			} else {
				n, err := reader.Seek(int64(recLen), io.SeekCurrent)
				if err != nil || n != int64(recLen) {
					if err == nil {
						err = errors.New("read len is error")
					}
				}
			}
		}
	}

	return
}

//解析图片相关结构
func (ole *OleInfo) getArtDggContainerData(documentReader io.ReadSeeker, data []byte) (err error) {
	artReader := bytes.NewReader(data[:])

	//==============drawingGroup(OfficeArtFDGGBlock) header
	var artFDGHeader [8]byte
	n, err := artReader.Read(artFDGHeader[:])
	if err != nil || n != 8 {
		if err == nil {
			err = errors.New("read len is error")
		}
		return
	}
	recIntFDG := binary.LittleEndian.Uint16(artFDGHeader[:2])
	recVerFDG := recIntFDG & 0xF
	recInstanceFDG := recIntFDG >> 4
	recTypeFDG := binary.LittleEndian.Uint16(artFDGHeader[2:4])
	recLenFDG := binary.LittleEndian.Uint32(artFDGHeader[4:])

	if recVerFDG != 0x0 || recInstanceFDG != 0x0 || recTypeFDG != 0xF006 {
		err = errors.New("recIntFDG is err")
		return
	}

	dGGBlockData := make([]byte, recLenFDG)
	n, err = artReader.Read(dGGBlockData[:])
	if err != nil || n != int(recLenFDG) {
		if err == nil {
			err = errors.New("read len is error")
		}
		return
	}
	dGGBlockReader := bytes.NewReader(dGGBlockData[:])

	//head(OfficeArtFDGG)
	type artFDGGStruct struct {
		SpidMax  uint32
		Cidcl    uint32
		CspSaved uint32
		CdgSaved uint32
	}
	var artFDGGstruct artFDGGStruct
	err = binary.Read(dGGBlockReader, binary.LittleEndian, &artFDGGstruct)
	if err != nil {
		return
	}

	//Rgidcl
	artDCL := make([]byte, 8*(artFDGGstruct.Cidcl-1))
	n, err = dGGBlockReader.Read(artDCL[:])
	if err != nil || n != int(8*(artFDGGstruct.Cidcl-1)) {
		if err == nil {
			err = errors.New("read len is error")
		}
		return
	}
	artDCLReader := bytes.NewReader(artDCL[:])
	type artDCLStruct struct {
		Dgid     uint32
		CspidCur uint32
	}
	artDCLstruct := make([]artDCLStruct, artFDGGstruct.Cidcl-1)
	err = binary.Read(artDCLReader, binary.LittleEndian, &artDCLstruct)
	if err != nil {
		return
	}

	//========blipStore(OfficeArtBStoreContainer) header
	var artBStoreHeader [8]byte
	n, err = artReader.Read(artBStoreHeader[:])
	if err != nil || n != 8 {
		if err == nil {
			err = errors.New("read len is error")
		}
		return
	}
	recIntBstore := binary.LittleEndian.Uint16(artBStoreHeader[:2])
	recVerBstore := recIntBstore & 0xF
	//BIP数量
	//recInstanceBstore := recIntBstore >> 4
	recTypeBstore := binary.LittleEndian.Uint16(artBStoreHeader[2:4])
	recLenBstore := binary.LittleEndian.Uint32(artBStoreHeader[4:])

	if recVerBstore != 0xF || recTypeBstore != 0xF001 {
		err = errors.New("recIntBstore is error")
		return
	}

	//fmt.Println(recInstanceBstore)

	//rgfb(OfficeArtBStoreContainer)数据
	artBstoreContainer := make([]byte, recLenBstore)
	n, err = artReader.Read(artBstoreContainer[:])
	if err != nil || n != int(recLenBstore) {
		if err == nil {
			err = errors.New("read len is error")
		}
		return
	}
	artBstoreReader := bytes.NewReader(artBstoreContainer[:])
	ole.getArtData(documentReader, artBstoreReader, false)
	return
}

//解析图片信息 && 获取图片数据
func (ole *OleInfo) getArtFBSEData(documentReader io.ReadSeeker, data []byte) (imageData []byte, err error) {
	artReader := bytes.NewReader(data[:])
	_, err = artReader.Seek(20, io.SeekCurrent)
	if err != nil {
		return
	}
	// var btWin32 [1]byte
	// _, err = artReader.Read(btWin32[:])
	// if err != nil {
	// 	return
	// }
	// var btMacOS [1]byte
	// _, err = artReader.Read(btMacOS[:])
	// if err != nil {
	// 	return
	// }

	// var rgbUid [16]byte
	// _, err = artReader.Read(rgbUid[:])
	// if err != nil {
	// 	return
	// }

	// var tag uint16
	// err = binary.Read(artReader, binary.LittleEndian, &tag)
	// if err != nil {
	// 	return
	// }

	type sizeStruct struct {
		//数据大小
		DataSize uint32
		CRef     uint32
		FoDelay  uint32
	}

	var fileSizePosInfo sizeStruct
	//var sizeInfo sizeStruct
	err = binary.Read(artReader, binary.LittleEndian, &fileSizePosInfo)
	if err != nil {
		return
	}

	_, err = artReader.Seek(1, io.SeekCurrent)
	if err != nil {
		return
	}
	//图片名字长度
	var chName [2]byte
	n, err := artReader.Read(chName[0:1])
	if err != nil || n != 1 {
		if err == nil {
			err = errors.New("read len is error")
		}
		return
	}
	chNameSize := binary.LittleEndian.Uint16(chName[:])

	_, err = artReader.Seek(2, io.SeekCurrent)
	if err != nil {
		return
	}

	if chNameSize != 0 {
		//图片名字
		nameData := make([]byte, chNameSize)
		n, err = artReader.Read(nameData[:])
		if err != nil || n != int(chNameSize) {
			if err == nil {
				err = errors.New("read len is error")
			}
			return
		}
	}

	blipData := make([]byte, fileSizePosInfo.DataSize)
	if len(data) == 36 {
		if documentReader == nil {
			return
		}
		seekRes, err := documentReader.Seek(int64(fileSizePosInfo.FoDelay), io.SeekStart)
		if err != nil || seekRes != int64(fileSizePosInfo.FoDelay) {
			if err == nil {
				err = errors.New("read len is error")
			}
			return nil, err
		}

		n, err = documentReader.Read(blipData)
		if err != nil || n != int(fileSizePosInfo.DataSize) {
			if err == nil {
				err = errors.New("read len is error")
			}
			return nil, err
		}
	} else {
		//图片结构数据(blip在本结构中)
		n, err = artReader.Read(blipData)
		if err != nil || n != int(fileSizePosInfo.DataSize) {
			if err == nil {
				err = errors.New("read len is error")
			}
			return
		}
	}

	//获取图片数据
	imageData, err = ole.getImageData(blipData)
	if err != nil {
		return
	}
	return
}

//获取图片数据
func (ole *OleInfo) getImageData(blipData []byte) (imageData []byte, err error) {
	blipReader := bytes.NewBuffer(blipData[:])
	var artRecordHeader [8]byte
	n, err := blipReader.Read(artRecordHeader[:])
	if err != nil || n != 8 {
		return nil, err
	}

	//rh
	recInt := binary.LittleEndian.Uint16(artRecordHeader[:2])
	recVer := recInt & 0xF
	recInstance := recInt >> 4
	//类型
	recType := binary.LittleEndian.Uint16(artRecordHeader[2:4])
	recLen := binary.LittleEndian.Uint32(artRecordHeader[4:])

	if recVer != 0x0 {
		err = errors.New("recVer is err")
		return nil, err
	}

	//去除头的图片结构数据
	artBlipPngData := make([]byte, recLen)
	n, err = blipReader.Read(artBlipPngData[:])
	if err != nil || n != int(recLen) {
		if err == nil {
			err = errors.New("read data len is error")
		}
		return
	}

	imageData, err = ole.getImageByType(artBlipPngData[:], recInstance, recType)

	return
}

//根据文件类型获取图片数据
func (ole *OleInfo) getImageByType(data []byte, recInstance, recType uint16) ([]byte, error) {
	var err error
	reader := bytes.NewReader(data[:])
	//图片类型
	switch recType {
	case 0xF01A: //emf
		if recInstance == 0x3D4 {
			_, err = reader.Seek(50, io.SeekCurrent)
			if err != nil {
				return nil, err
			}
		} else if recInstance == 0x3D5 {
			_, err = reader.Seek(66, io.SeekCurrent)
			if err != nil {
				return nil, err
			}
		}
	case 0xF01B: //wmf
		//fmt.Println("wmf")
	case 0xF01C: //pict
		//fmt.Println("pict")
	case 0xF01D, 0xF02A: //jpeg
		if recInstance == 0x46A || recInstance == 0x6E2 {
			_, err = reader.Seek(17, io.SeekCurrent)
			if err != nil {
				return nil, err
			}
		} else if recInstance == 0x46B || recInstance == 0x6E3 {
			_, err = reader.Seek(33, io.SeekCurrent)
			if err != nil {
				return nil, err
			}
		}
	case 0xF01E: //png
		if recInstance == 0x6E0 {
			_, err = reader.Seek(17, io.SeekCurrent)
			if err != nil {
				return nil, err
			}
		} else if recInstance == 0x6E1 {
			_, err = reader.Seek(33, io.SeekCurrent)
			if err != nil {
				return nil, err
			}
		}
	case 0xF01F: //dib
		//fmt.Println("dib")
	case 0xF029: //tiff
		//fmt.Println("tiff")
	default:
		err = errors.New("type is not find")
		return nil, err
	}

	//图片数据
	var buf bytes.Buffer
	io.Copy(&buf, reader)
	return buf.Bytes(), nil
}
