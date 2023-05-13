package ole

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

const (
	IMAGEEMF = iota
	IMAGEWMF
	IMAGEPICT
	IMAGEJPEG
	IMAGEPNG
	IMAGEDIB
	IMAGETIFF
)

func (ole *OleInfo) getInlinePicInfo(documentReader, dataDirReader io.ReadSeeker, object, root *Directory, callBack DataCallBackFunc) (err error) {
	// 读取FIBTable97中的fcPlcfBteChpx,lcPlcfBteChpx
	plcfBteChpxdata, err := ole.getFcLcPosiData(documentReader, FCPLCFBTECHPX)
	if err != nil {
		return
	}

	// PlcfBteChpx结构偏移量
	fcPlcfBteChpxSeek := binary.LittleEndian.Uint32(plcfBteChpxdata[:4])
	// PlcfBteChpx结构大小
	lcPlcfBteChpx := binary.LittleEndian.Uint32(plcfBteChpxdata[4:])

	if _, err = ole.fileReadSeeker.Seek(int64(fcPlcfBteChpxSeek+lcPlcfBteChpx), io.SeekStart); err != nil {
		err = errors.New("fcclx and size is error")
		return
	}

	// 获取PlcfBteChpx结构数据
	plcfBteChpx, err := ole.getFcData(fcPlcfBteChpxSeek, lcPlcfBteChpx, object, root)
	if err != nil {
		return
	}

	ole.getChpxBte(plcfBteChpx, documentReader, dataDirReader)

	plcfBteChpx = nil
	return
}

func (ole *OleInfo) getChpxBte(data []byte, documentReader, dataDirReader io.ReadSeeker) (err error) {
	if len(data)%4 != 0 {
		err = errors.New("data len is error")
		return
	}
	allNum := len(data) / 4
	//aPnBteChpx数量
	aPnNum := allNum / 2
	//aFC数量
	aFCNum := aPnNum + 1
	if aPnNum+aFCNum != allNum {
		err = errors.New("data len is error 2")
		return
	}

	//var aFCs, pns []uint32
	index := 0
	for i := 0; i < aFCNum; i++ {
		//保存的亦是文本数据
		//aFCs = append(aFCs, binary.LittleEndian.Uint32(data[index:index+4]))
		index += 4
	}

	pnMap := &sync.Map{}
	for i := 0; i < aPnNum; i++ {
		defer func() { index += 4 }()
		pnfkpchpx := binary.LittleEndian.Uint32(data[index : index+4])
		//取前22位
		pn := pnfkpchpx & 0x3FFFFFFF
		//pns = append(pns, pn)

		_, ok := pnMap.Load(pn)
		if ok {
			continue
		}
		pnMap.Store(pn, 1)
		_, err = documentReader.Seek(int64(pn*512), io.SeekStart)
		if err != nil {
			continue
		}
		var chpxFkpData [512]byte
		n, err := documentReader.Read(chpxFkpData[:])
		if err != nil || n != 512 {
			if err == nil {
				err = errors.New("read len is error")
			}
			continue
		}

		err = ole.getChpxs(chpxFkpData[:], documentReader, dataDirReader)
		if err != nil {
			continue
		}
	}

	return
}

func (ole *OleInfo) getChpxs(data []byte, documentReader, dataDirReader io.ReadSeeker) (err error) {
	if len(data) != 512 {
		err = errors.New("data len is not 512")
		return
	}

	//获取crun字段
	crun := int64(data[511])
	if crun == 0 {
		err = errors.New("crun is 0")
		return
	}

	dataReader := bytes.NewReader(data[:])

	//移动过滤rgfc
	_, err = dataReader.Seek((crun+1)*4, io.SeekStart)
	//读取rgb数组
	var rgbs []uint32
	rgbsMap := &sync.Map{}
	for i := int64(0); i < crun; i++ {
		var rgbData [4]byte
		n, err := dataReader.Read(rgbData[0:1])
		if err != nil || n != 1 {
			if err == nil {
				err = errors.New("read len is error")
			}
			return err
		}
		rgb := binary.LittleEndian.Uint32(rgbData[:])
		_, ok := rgbsMap.Load(rgb)
		if !ok {
			rgbs = append(rgbs, rgb)
			rgbsMap.Store(rgb, 1)
		}
	}

	offsetMap := &sync.Map{}
	for _, v := range rgbs {
		_, err = dataReader.Seek(int64(v*2), io.SeekStart)
		if err != nil {
			return
		}
		//大小
		var cb [4]byte
		n, err := dataReader.Read(cb[0:1])
		if err != nil || n != 1 {
			if err == nil {
				err = errors.New("read len is error")
			}
			return err
		}

		//chpx数据
		chpxLen := binary.LittleEndian.Uint32(cb[:])
		chpxData := make([]byte, chpxLen)
		n, err = dataReader.Read(chpxData[:])
		if err != nil || n != int(chpxLen) {
			if err == nil {
				err = errors.New("read len is error")
			}
			return err
		}

		//sprm
		var sprm uint16
		chpxReader := bytes.NewReader(chpxData[:])

		//获取chpx结构（包括sprm、operand）
		var offset int64
		for {
			err = binary.Read(chpxReader, binary.LittleEndian, &sprm)
			if err != nil {
				break
			}

			if sprm != 0x6A03 {
				continue
			}
			// ispmd := sprm & 0x01FF
			// fSpec := (sprm / 512) & 0x1
			// sgc := (sprm / 1024) & 0x7
			spra := sprm / 8192
			//0-1, 1-1, 2-2, 3-4, 4-2, 5-2, 6-(不定，暂不处理) , 7-3
			switch spra {
			case 0, 1:
				operand := make([]byte, 2)
				n, err = chpxReader.Read(operand[0:1])
				if err != nil || n != 1 {
					if err == nil {
						err = errors.New("read len is error")
					}
					break
				}
				offset = int64(binary.LittleEndian.Uint16(operand[:]))
			case 2, 4, 5:
				operand := make([]byte, 2)
				n, err = chpxReader.Read(operand[:])
				if err != nil || n != 2 {
					if err == nil {
						err = errors.New("read len is error")
					}
					break
				}
				offset = int64(binary.LittleEndian.Uint16(operand[:]))
			case 3:
				operand := make([]byte, 4)
				n, err = chpxReader.Read(operand[:])
				if err != nil || n != 4 {
					if err == nil {
						err = errors.New("read len is error")
					}
					break
				}
				offset = int64(binary.LittleEndian.Uint32(operand[:]))
			case 7:
				operand := make([]byte, 4)
				n, err = chpxReader.Read(operand[0:3])
				if err != nil || n != 3 {
					if err == nil {
						err = errors.New("read len is error")
					}
					break
				}
				offset = int64(binary.LittleEndian.Uint32(operand[:]))
			default:
				err = errors.New("spra is not find")
			}
			if err == nil {
				_, ok := offsetMap.Load(offset)
				if !ok {
					offsetMap.Store(offset, 1)
					_, err = dataDirReader.Seek(offset, io.SeekStart)
					if err != nil {
						continue
					}
					//获取内联图片
					ole.getArtData(documentReader, dataDirReader, true)
				}
			}
		}
	}

	return
}

func (ole *OleInfo) getFloatingPicInfo(reader io.ReadSeeker, object, root *Directory, callBack DataCallBackFunc) (err error) {
	// 读取FIBTable97中的fcDggInfo,lccDggInfo
	plcfBteChpxdata, err := ole.getFcLcPosiData(reader, FCDGGINFO)
	if err != nil {
		return
	}

	// PlcfBteChpx结构偏移量
	fcDggInfoSeek := binary.LittleEndian.Uint32(plcfBteChpxdata[:4])
	// PlcfBteChpx结构大小
	lccDggInfo := binary.LittleEndian.Uint32(plcfBteChpxdata[4:])

	if _, err = ole.fileReadSeeker.Seek(int64(fcDggInfoSeek+lccDggInfo), io.SeekStart); err != nil {
		err = errors.New("fcclx and size is error")
		return
	}

	// 获取PlcfBteChpx结构数据
	fcDggInfo, err := ole.getFcData(fcDggInfoSeek, lccDggInfo, object, root)
	if err != nil {
		return
	}

	// plcSpaMomdata, err := ole.getFcLcPosiData(reader, FCSPAMOM)
	// if err != nil {
	// 	return
	// }
	// // PlcfBteChpx结构偏移量
	// fcSpaMomSeek := binary.LittleEndian.Uint32(plcSpaMomdata[:4])
	// // PlcfBteChpx结构大小
	// lccSpaMom := binary.LittleEndian.Uint32(plcSpaMomdata[4:])

	// if _, err = ole.fileReadSeeker.Seek(int64(fcSpaMomSeek+lccSpaMom), io.SeekStart); err != nil {
	// 	err = errors.New("fcclx and size is error")
	// 	return
	// }

	// // 获取PlcfBteChpx结构数据
	// fcSpaMom, err := ole.getFcData(fcSpaMomSeek, lccSpaMom, object, root)
	// if err != nil {
	// 	return
	// }

	// libs, err := ole.getPlcSpaInfo(fcSpaMom)
	// if err != nil {
	// 	return
	// }

	// fmt.Println(libs)

	ole.getArtData(reader, bytes.NewReader(fcDggInfo[:]), false)
	return
}

func (ole *OleInfo) getPlcSpaInfo(fcSpaMomInfo []byte) (lids []uint32, err error) {
	spaNum := (len(fcSpaMomInfo) - 4) / 30
	if spaNum == 0 {
		err = errors.New("spaNum is 0")
		return
	}
	fcSpaMomReader := bytes.NewReader(fcSpaMomInfo[:])

	// var aCPs []uint32
	// for i := 0; i <= spaNum; i++ {
	// 	var aCP uint32
	// 	err = binary.Read(fcSpaMomReader, binary.LittleEndian, &aCP)
	// 	if err != nil {
	// 		return
	// 	}
	// 	aCPs = append(aCPs, aCP)
	// }

	n, err := fcSpaMomReader.Seek(int64(4*(spaNum+1)), io.SeekStart)
	if err != nil || n != int64(4*(spaNum+1)) {
		return nil, err
	}

	for i := 0; i < spaNum; i++ {
		var lid uint32
		err = binary.Read(fcSpaMomReader, binary.LittleEndian, &lid)
		if err != nil {
			return
		}
		lids = append(lids, lid)
		n, err = fcSpaMomReader.Seek(22, io.SeekCurrent)
		if err != nil || n != 22 {
			return
		}
	}

	return
}

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
				fmt.Println("emf")
			case 0x03:
				pathName += ".wmf"
				fmt.Println("wmf")
			case 0x04:
				pathName += ".pict"
				fmt.Println("pict")
			case 0x05, 0x12:
				pathName += ".jpeg"
				fmt.Println("jpeg")
			case 0x06:
				pathName += ".png"
				fmt.Println("png")
			case 0x07:
				pathName += ".dib"
				fmt.Println("dib")
			case 0x011:
				pathName += ".tiff"
				fmt.Println("tiff")
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
				ole.getImageByType(artData[:], recInstance, recType)
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
		fmt.Println("wmf")
	case 0xF01C: //pict
		fmt.Println("pict")
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
		fmt.Println("dib")
	case 0xF029: //tiff
		fmt.Println("tiff")
	default:
		err = errors.New("type is not find")
		return nil, err
	}

	//图片数据
	var buf bytes.Buffer
	io.Copy(&buf, reader)
	return buf.Bytes(), nil
}
