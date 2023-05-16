package ole

/*
Date：2023.05.13
Author：scl
Description：解析doc文件中的图片（嵌入/浮动）
*/
import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"sync"
)

//获取内联图片
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

//获取ChpxBte结构
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

//解析ChpxBte结构
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

//获取浮动图片
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

//解析PlcSpa结构
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
