package fileparser

/*
Date：2023.03.02
Author：scl
Description：解析rtf文件
*/

import (
	"bufio"
	"errors"
	"os"
	"strings"

	"github.com/henrylee2cn/pholcus/common/mahonia"
)

//rtf文件字段
const (
	rtf_unknown = iota + 1
	rtf_b
	rtf_bin
	rtf_blue
	rtf_brdrnone
	rtf_bullet
	rtf_caps
	rtf_cb
	rtf_cell
	rtf_cellx
	rtf_cf
	rtf_clbrdrb
	rtf_clbrdrl
	rtf_clbrdrr
	rtf_clbrdrt
	rtf_clvertalb
	rtf_clvertalc
	rtf_clvertalt
	rtf_clvmgf
	rtf_clvmrg
	rtf_colortbl
	rtf_dn
	rtf_emdash
	rtf_emspace
	rtf_endash
	rtf_enspace
	rtf_fi
	rtf_field
	rtf_filetbl
	rtf_f
	rtf_fprq
	rtf_fcharset
	rtf_fnil
	rtf_froman
	rtf_fswiss
	rtf_fmodern
	rtf_fscript
	rtf_fdecor
	rtf_ftech
	rtf_fbidi
	rtf_fldrslt
	rtf_fonttbl
	rtf_footer
	rtf_footerf
	rtf_fs
	rtf_green
	rtf_header
	rtf_headerf
	rtf_highlight
	rtf_i
	rtf_info
	rtf_intbl
	rtf_ldblquote
	rtf_li
	rtf_line
	rtf_lquote
	rtf_margl
	rtf_object
	rtf_paperw
	rtf_par
	rtf_pard
	rtf_pict
	rtf_plain
	rtf_qc
	rtf_qj
	rtf_ql
	rtf_qmspace
	rtf_qr
	rtf_rdblquote
	rtf_red
	rtf_ri
	rtf_row
	rtf_rquote
	rtf_sa
	rtf_sb
	rtf_sect
	rtf_softline
	rtf_stylesheet
	rtf_sub
	rtf_super
	rtf_tab
	rtf_title
	rtf_trleft
	rtf_trowd
	rtf_trrh
	rtf_ul
	rtf_ulnone
	rtf_up
	rtf_unicode
	rtf_asterisk
	rtf_assi
)

//rtf部分关键字
var keyword_map map[string]int = map[string]int{
	"b":          rtf_b,
	"bin":        rtf_bin,
	"blue":       rtf_blue,
	"brdrnone":   rtf_brdrnone,
	"bullet":     rtf_bullet,
	"caps":       rtf_caps,
	"cb":         rtf_cb,
	"cell":       rtf_cell,
	"cellx":      rtf_cellx,
	"cf":         rtf_cf,
	"clbrdrb":    rtf_clbrdrb,
	"clbrdrl":    rtf_clbrdrl,
	"clbrdrr":    rtf_clbrdrr,
	"clbrdrt":    rtf_clbrdrt,
	"clvertalb":  rtf_clvertalb,
	"clvertalc":  rtf_clvertalc,
	"clvertalt":  rtf_clvertalt,
	"clvmgf":     rtf_clvmgf,
	"clvmrg":     rtf_clvmrg,
	"colortbl":   rtf_colortbl,
	"dn":         rtf_dn,
	"emdash":     rtf_emdash,
	"emspace":    rtf_emspace,
	"endash":     rtf_endash,
	"enspace":    rtf_enspace,
	"f":          rtf_f,
	"fprq":       rtf_fprq,
	"fcharset":   rtf_fcharset,
	"fnil":       rtf_fnil,
	"froman":     rtf_froman,
	"fswiss":     rtf_fswiss,
	"fmodern":    rtf_fmodern,
	"fscript":    rtf_fscript,
	"fdecor":     rtf_fdecor,
	"ftech":      rtf_ftech,
	"fbidi":      rtf_fbidi,
	"field":      rtf_field,
	"filetbl":    rtf_filetbl,
	"fldrslt":    rtf_fldrslt,
	"fonttbl":    rtf_fonttbl,
	"footer":     rtf_footer,
	"footerf":    rtf_footerf,
	"fs":         rtf_fs,
	"green":      rtf_green,
	"header":     rtf_header,
	"headerf":    rtf_headerf,
	"highlight":  rtf_highlight,
	"i":          rtf_i,
	"info":       rtf_info,
	"intbl":      rtf_intbl,
	"ldblquote":  rtf_ldblquote,
	"li":         rtf_li,
	"line":       rtf_line,
	"lquote":     rtf_lquote,
	"margl":      rtf_margl,
	"object":     rtf_object,
	"paperw":     rtf_paperw,
	"par":        rtf_par,
	"pard":       rtf_pard,
	"pict":       rtf_pict,
	"plain":      rtf_plain,
	"qc":         rtf_qc,
	"qj":         rtf_qj,
	"ql":         rtf_ql,
	"qr":         rtf_qr,
	"rdblquote":  rtf_rdblquote,
	"red":        rtf_red,
	"ri":         rtf_ri,
	"row":        rtf_row,
	"rquote":     rtf_rquote,
	"sa":         rtf_sa,
	"sb":         rtf_sb,
	"sect":       rtf_sect,
	"softline":   rtf_softline,
	"stylesheet": rtf_stylesheet,
	"sub":        rtf_sub,
	"super":      rtf_super,
	"tab":        rtf_tab,
	"title":      rtf_title,
	"trleft":     rtf_trleft,
	"trowd":      rtf_trowd,
	"trrh":       rtf_trrh,
	"ul":         rtf_ul,
	"ulnone":     rtf_ulnone,
	"up":         rtf_up,
	"u":          rtf_unicode,
	"*":          rtf_asterisk,
	"'":          rtf_assi,
}

//判断是否为英文
func Isalpha(data []byte) bool {
	if (data[0] >= 0x41 && data[0] <= 0x5a) || (data[0] >= 0x61 && data[0] <= 0x7a) {
		return true
	}
	return false
}

//判断是否为数字
func Isdigit(data []byte) bool {
	if data[0] >= 0x30 && data[0] <= 0x39 {
		return true
	}
	return false
}

//下一个关键数组
func NextChar(reader *bufio.Reader) (bufdata []byte, lastdata [1]byte) {
	for {
		var data [1]byte
		_, err := reader.Read(data[:])
		if err != nil {
			return
		}

		if !Isalpha(data[:]) {
			if Isdigit(data[:]) || data[0] == '-' {
				continue
			} else if data[0] == '*' {
				bufdata = append(bufdata, data[0])
				break
			} else {
				//由于跟指针不一样，取出来放不回去，所以保存一下最后一个字节，外边会用到
				lastdata[0] = data[0]
				break
			}
		}
		bufdata = append(bufdata, data[0])
	}
	return
}

//byte字符转16进制（2个字符转换成1个字节）
func Strtol(chardata []byte) (strbyte []byte) {
	for i := 0; i <= len(chardata)/2; i += 2 {
		ret := 0
		for j := 0; j < 2; j++ {
			char := chardata[i+j]
			bytedata := []byte{char}
			ch := int(char)
			if Isdigit(bytedata) {
				ch -= '0'
			} else if Isalpha(bytedata) {
				if char >= 'A' && char <= 'Z' {
					ch -= 'A' - 10
				} else {
					ch -= 'a' - 10
				}
			} else {
				break
			}

			if ch >= 16 {
				break
			}

			ret *= 16
			ret += ch
		}
		strbyte = append(strbyte, byte(ret))
	}

	return strbyte
}

//打开文件
func GetRtfData(fileName string, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	f, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	var rtfdata [5]byte
	_, err = reader.Read(rtfdata[:])
	if err != nil {
		return err
	}

	if !strings.Contains(string(rtfdata[:]), "{\\rtf") {
		err = errors.New("is not rtf")
		return err
	}

	reader.Reset(f)

	err = DealRtfFile(reader, callBack)

	return
}

//处理rtf文件
func DealRtfFile(reader *bufio.Reader, callBack CallBackDataFunc) (err error) {
	var lastbyte [1]byte
	var str string
	var data [1]byte
	for {
		if lastbyte[0] == '\\' {
			data[0] = '\\'
			lastbyte[0] = 0x00
		} else {
			_, err := reader.Read(data[:])
			if err != nil {
				break
			}
		}

		//var keyword string
		switch data[0] {
		//"{","}",13,10,"<",">"
		case '{', '}', 13, 10, '<', '>':
		//"\"
		case '\\':
			if len(str) > 500 {
				if !callBack(str, "") {
					return nil
				}

				str = ""
			}

			var chardata []byte
			chardata, lastbyte = NextChar(reader)
			if len(chardata) == 0 && lastbyte[0] != '\'' {
				break
			}

			lastkey := keyword_map[string(lastbyte[0])]
			//解析汉字
			switch lastkey {
			case rtf_assi:
				var assi_byte []byte
				count := 0
				//2个字节一个汉字
				for count < 4 {
					var charbyte [1]byte
					_, err = reader.Read(charbyte[:])
					if err != nil {
						break
					}
					if !Isalpha(charbyte[:]) && !Isdigit(charbyte[:]) {
						continue
					}
					assi_byte = append(assi_byte, charbyte[0])
					count++
				}
				//num,_ := strconv.ParseUint(string(assi_byte),16,8)
				hexbyte := Strtol(assi_byte)

				ec := mahonia.NewDecoder("gbk")
				ecstr := ec.ConvertString(string(hexbyte))
				str += ecstr
			}

			num := keyword_map[string(chardata)]
			switch num {
			//过滤fonttbl
			case rtf_fonttbl:
				var font_data [1]byte
				_, err = reader.Read(font_data[:])
				if err != nil {
					return
				}

				var isfont bool
				for {
					//"}"
					if isfont && font_data[0] == '}' {
						break
					}

					if font_data[0] == 0x7d {
						isfont = true
					} else {
						isfont = false
					}

					_, err = reader.Read(font_data[:])
					if err != nil {
						return
					}
				}
			//过滤colortbl
			case rtf_colortbl:
				var color_data [1]byte
				//"}"
				for color_data[0] != '}' {
					_, err := reader.Read(color_data[:])
					if err != nil {
						break
					}
				}
			//过滤不需要的
			case rtf_asterisk, rtf_filetbl, rtf_stylesheet, rtf_header, rtf_footer, rtf_headerf, rtf_footerf, rtf_pict, rtf_object, rtf_info:
				var rightcut int64 = 1
				for rightcut > 0 {
					var asterisk_data [1]byte
					if lastbyte[0] == '{' {
						asterisk_data[0] = '{'
						lastbyte[0] = 0x00
					} else {
						_, err := reader.Read(asterisk_data[:])
						if err != nil {
							break
						}
					}

					switch asterisk_data[0] {
					case '{':
						rightcut++
					case '}':
						rightcut--
					}
				}
			case rtf_par, rtf_sect:
				str += "\n"
			}

		default:
			str += string(data[0])
		}
	}

	if len(str) > 0 {
		if !callBack(str, "") {
			return nil
		}
	}

	return
}
