package publicfunc

/*
Date：2023.03.02
Author：scl
Description：公共函数
*/

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"unicode/utf16"

	"github.com/henrylee2cn/pholcus/common/mahonia"
)

type FILECONFIG struct {
	//扫描路径
	Path string
	//扫描文件类型
	Filetype string
	//扫描深度 -1 全部扫描
	Depth int32
}

// 格式化大小
func FormatSize(fileSize uint64) (size string) {
	if fileSize < 1024 {
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else {
		// if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

//获取字符串的md5
func GetFileMd5(filedata []byte) string {
	hash := md5.New()
	hash.Write(filedata)
	md5str := fmt.Sprintf("%x", hash.Sum(nil))
	return md5str
}

//获取文件md5
func GetFileMd5f(pfile *os.File) (string, error) {
	if pfile == nil {
		err := errors.New("md5 os.file is nil")
		return "", err
	}
	defer pfile.Seek(0, io.SeekStart)
	md5hash := md5.New()
	written, err := io.Copy(md5hash, pfile)
	if written == 0 || err != nil {
		if err == nil {
			err = errors.New("md5 iocopy len is 0")
		}
		return "", err
	}

	return hex.EncodeToString(md5hash.Sum(nil)), nil
}

//获取文件sha1
func GetFileSha1f(pfile *os.File) (string, error) {
	if pfile == nil {
		err := errors.New("md5 os.file is nil")
		return "", err
	}
	defer pfile.Seek(0, io.SeekStart)
	sha1hash := sha1.New()
	written, err := io.Copy(sha1hash, pfile)
	if written == 0 || err != nil {
		if err == nil {
			err = errors.New("sha1 iocopy len is 0")
		}
		return "", err
	}
	return hex.EncodeToString(sha1hash.Sum(nil)), nil
}

//获取文件sha256
func GetFileSha256f(pfile *os.File) (string, error) {
	if pfile == nil {
		err := errors.New("md5 os.file is nil")
		return "", err
	}
	defer pfile.Seek(0, io.SeekStart)
	sha256hash := sha256.New()
	written, err := io.Copy(sha256hash, pfile)
	if written == 0 || err != nil {
		if err == nil {
			err = errors.New("sha256 iocopy len is 0")
		}
		return "", err
	}
	return hex.EncodeToString(sha256hash.Sum(nil)), nil
}

func UTF16ToString(s []uint16) string {
	for i, v := range s {
		if v == 0 {
			s = s[:i]
			break
		}
	}
	return string(utf16.Decode(s))
}

//转码
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)

	srcResult := srcCoder.ConvertString(src)

	tagCoder := mahonia.NewDecoder(tagCode)

	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)

	result := string(cdata)

	return result
}
