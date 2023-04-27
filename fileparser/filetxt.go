package fileparser

/*
Date：2023.03.10
Author：scl
Description：解析txt、csv文件
*/
import (
	"bufio"
	"errors"
	"fileparser/decoder"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/saintfish/chardet"
)

//获取文件数据
func GetTxtDataFile(fileName string, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil || fi.Size() == 0 {
		if err == nil {
			err = errors.New("file size is nil")
		}
		return
	}

	err = GetTxtData(file, fi.Size(), callBack)
	return
}

//获取文件数据
func GetTxtData(fileReadSeeker io.ReadSeeker, fileSize int64, callBack CallBackDataFunc) (err error) {
	if callBack == nil || fileReadSeeker == nil || fileSize == 0 {
		err = errors.New("callBack is nil or io.ReadSeeker is nil or fileSize is 0")
		return
	}

	//读取100个字节判断文件编码
	var data [100]byte
	_, err = fileReadSeeker.Read(data[:])
	//TODO
	if err != nil {
		if err == io.EOF {
			//文件小于100个字节，直接当utf8字节处理
			if fileSize < 100 && fileSize > 0 {
				filedata := make([]byte, fileSize)
				_, err = fileReadSeeker.Read(filedata[:])
				if err != nil {
					return err
				}
				str := string(filedata)

				callBack(str, "第1行")

				return nil
			} else {
				return err
			}
		} else {
			return
		}
	}

	//获取编码信息
	detector := chardet.NewTextDetector()
	charset, err := detector.DetectBest(data[:])
	if err != nil {
		return
	}
	_, err = fileReadSeeker.Seek(0, io.SeekStart)
	if err != nil {
		return
	}

	//gbk、csv文件
	if strings.Contains(charset.Charset, "GB") || strings.Contains(charset.Charset, "KOI8-R") ||
		strings.Contains(charset.Charset, "UTF-16LE") || strings.Contains(charset.Charset, "UTF-16BE") {
		//文件流编码解析
		trancreader := decoder.EncodeReader(charset.Charset, fileReadSeeker)
		if trancreader == nil {
			err = errors.New("transreader is nil")
			return
		}

		linenum := 0
		sc := bufio.NewScanner(trancreader)
		//行读取
		for sc.Scan() {
			linenum++
			if err = sc.Err(); err != nil {
				break
			}

			if len(string(sc.Bytes())) == 0 {
				continue
			}

			if !callBack(string(sc.Bytes()), "第"+strconv.Itoa(linenum)+"行") {
				return nil
			}
		}

	} else /*if strings.Contains(charset.Charset, "UTF-8")*/ {
		lineNum := 0
		//utf8编码
		reader := bufio.NewReader(fileReadSeeker)
		//行读取
		for {
			lineNum++
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				continue
			}

			if len(line) > 0 && !callBack(string(line), "第"+strconv.Itoa(lineNum)+"行") {
				return nil
			}
			if err == io.EOF {
				break
			}
		}
	}

	return
}
