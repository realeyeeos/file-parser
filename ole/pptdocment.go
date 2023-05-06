package ole

/*
Date：2023.03.02
Author：scl
Description：解析ppt文件
*/
import (
	"bytes"
	"encoding/binary"
	"io"

	publicfun "github.com/realeyeeos/file-parser/publicfunc"
)

const (
	//有数据的类型
	RT_TEXT_CHARS_ATOM = uint16(4000)
	RT_TEXT_BYTES_ATOM = uint16(4008)
	RT_CSTRING         = uint16(4026)

	//不需要移动的类型
	RT_DOCUMENT               = uint16(1000)
	RT_SLIDE_BASE             = uint16(1004)
	RT_SLIDE                  = uint16(1006)
	RT_DRAWING                = uint16(1036)
	RT_LIST                   = uint16(2000)
	RT_SLIDE_LIST_WITH_TEXT   = uint16(4080)
	OFFICE_ART_DG_CONTAINER   = uint16(61442)
	OFFICE_ART_SPGR_CONTAINER = uint16(61443)
	OFFICE_ART_SP_CONTAINER   = uint16(61444)
	OFFICE_ART_CLIENT_TEXTBOX = uint16(0xF00D) //61453
)

//获取ppt信息
func (ole *OleInfo) getPptInfo(reader io.ReadSeeker, ppt *Directory, callBack DataCallBackFunc) (err error) {
	var str string
	//循环获取数据
	for {
		if len(str) >= 500 {
			if !callBack(str, "") {
				return nil
			}
			str = ""
		}

		var buf [8]byte
		_, err = reader.Read(buf[:])
		if err != nil {
			break
		}

		rec_type := binary.LittleEndian.Uint16(buf[2:4])
		rec_len := binary.LittleEndian.Uint32(buf[4:8])

		tmpstr, err := ole.getRecData(rec_type, reader, rec_len)
		if len(tmpstr) > 0 && err == nil {
			str += tmpstr
		}
		rec_type = uint16(0)
	}

	if len(str) > 0 {
		callBack(str, "")
		str = ""
	}

	return
}

//获取数据流中的数据
func (ole *OleInfo) getRecData(rec_type uint16, ppt_reader io.ReadSeeker, reclen uint32) (str string, err error) {
	if rec_type == RT_DOCUMENT || rec_type == RT_SLIDE_BASE || rec_type == RT_SLIDE || rec_type == RT_DRAWING || rec_type == RT_LIST ||
		rec_type == RT_SLIDE_LIST_WITH_TEXT || rec_type == OFFICE_ART_CLIENT_TEXTBOX || rec_type == OFFICE_ART_DG_CONTAINER ||
		rec_type == OFFICE_ART_SPGR_CONTAINER || rec_type == OFFICE_ART_SP_CONTAINER {
		return
	}

	//unicode存储
	if rec_type == RT_TEXT_CHARS_ATOM || rec_type == RT_CSTRING {
		bufdata := make([]byte, reclen)
		_, err = ppt_reader.Read(bufdata[:])
		if err != nil {
			return
		}
		textlen := reclen / 2
		//unicode字符解析
		ioread := bytes.NewReader(bufdata[:textlen*2])
		uintmp := make([]uint16, textlen)
		binary.Read(ioread, binary.LittleEndian, &uintmp)
		str = publicfun.UTF16ToString(uintmp[:])

		//打印数据
		//fmt.Println(str2)
		return
	}

	//ansi存储
	if rec_type == RT_TEXT_BYTES_ATOM {
		bufdata := make([]byte, reclen)
		_, err = ppt_reader.Read(bufdata[:])
		if err != nil {
			return
		}

		str = string(bufdata)
		//打印数据
		//fmt.Println(string(bufdata))
	}

	_, err = ppt_reader.Seek(int64(reclen), 1)
	if err != nil {
		return
	}

	return
}
