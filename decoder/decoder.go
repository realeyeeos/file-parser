package decoder

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"encoding/base64"
	"mime/quotedprintable"

	"github.com/gogf/gf/crypto/gaes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

//转码
func EncodeReader(charset string, reader io.Reader) (trancreader *transform.Reader) {
	upercharset := strings.ToUpper(charset)
	if strings.Contains(upercharset, "UTF-16BE") { //utf16be
		enc := unicode.UTF16(unicode.BigEndian, unicode.UseBOM)
		trancreader = transform.NewReader(reader, enc.NewDecoder())
	} else if strings.Contains(upercharset, "UTF-16LE") { //utf16le
		enc := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
		trancreader = transform.NewReader(reader, enc.NewDecoder())
	} else { //gbk
		enc := simplifiedchinese.GBK
		trancreader = transform.NewReader(reader, enc.NewDecoder())
	}
	return trancreader
}

func UTF8(cs string, data []byte) ([]byte, error) {
	if strings.Contains(strings.ToUpper(cs), "UTF-8") {
		return data, nil
	}

	//转码
	trancreader := EncodeReader(cs, bytes.NewReader(data))

	if trancreader == nil {
		return nil, nil
	}
	return ioutil.ReadAll(trancreader)

	// r, err := charset.NewReader(bytes.NewReader(data), cs)
	// if err != nil {
	// 	return []byte{}, err
	// }

	// return ioutil.ReadAll(r)

}

func Parse(bstr []byte) ([]byte, error) {
	var err error
	var retbyte []byte
	slite := strings.Split(string(bstr), "?=")
	for _, v := range slite {
		strs := regexp.MustCompile(".*=\\?(.*?)\\?(.*?)\\?(.*)$").FindAllStringSubmatch(v, -1)
		if len(strs) > 0 && len(strs[0]) == 4 {
			c := strs[0][1]
			e := strs[0][2]
			dstr := strs[0][3]

			bstr, err = Decode(e, []byte(dstr))
			if err != nil {
				return bstr, err
			}

			utfbyte, err := UTF8(c, bstr)
			if err != nil {
				continue
			}
			retbyte = append(retbyte, utfbyte...)
		}

	}

	return retbyte, err

}

func Decode(e string, bstr []byte) ([]byte, error) {
	var err error
	switch strings.ToUpper(e) {
	case "Q", "QUOTED-PRINTABLE":
		bstr, err = ioutil.ReadAll(quotedprintable.NewReader(bytes.NewReader(bstr)))
	case "B", "BASE64":
		bstr, err = base64.StdEncoding.DecodeString(string(bstr))
	default:
		//not set encoding type

	}
	return bstr, err
}

// 加密
func aesCtrCrypt(plainText []byte, key []byte) ([]byte, error) {
	//1. 创建cipher.Block接口
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//2. 创建分组模式，在crypto/cipher包中
	iv := bytes.Repeat([]byte("1"), block.BlockSize())
	stream := cipher.NewCTR(block, iv)
	//3. 加密
	dst := make([]byte, len(plainText))
	stream.XORKeyStream(dst, plainText)
	return dst, nil
}

// 与Agent之间交互使用的加解密算法
func AgentEncrypt(plainText string) (string, error) {
	key := []byte{0xAE, 0x41, 0x6D, 0xE2, 0xA6, 0x81, 0xB2, 0xFC, 0x7d, 0x4F, 0x9B, 0xB9, 0x4C, 0xb7, 0x88, 0x66}
	plainTextBytes := []byte(plainText)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	plainTextBytes = gaes.PKCS5Padding(plainTextBytes, blockSize)
	plainTextBytes, err = aesCtrCrypt(plainTextBytes, []byte(key))
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encryptedData := make([]byte, len(plainTextBytes))
	blockMode.CryptBlocks(encryptedData, plainTextBytes)
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

// 与Agent之间交互使用的加解密算法
func AgentDecrypt(encryptedData string) (string, error) {
	key := []byte{0xAE, 0x41, 0x6D, 0xE2, 0xA6, 0x81, 0xB2, 0xFC, 0x7d, 0x4F, 0x9B, 0xB9, 0x4C, 0xb7, 0x88, 0x66}
	encryptedDataBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	decrypted := make([]byte, len(encryptedDataBytes))
	blockMode.CryptBlocks(decrypted, encryptedDataBytes)
	decrypted, err = aesCtrCrypt(decrypted, []byte(key))
	if err != nil {
		return "", err
	}
	decrypted, err = gaes.PKCS5UnPadding(decrypted, blockSize)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}
