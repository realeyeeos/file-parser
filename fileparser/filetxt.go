package fileparser

/*
Date：2023.03.10
Author：scl
Description：解析txt、csv文件
*/
import (
	"bufio"
	"collector/decoder"
	"errors"
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

	err = GetTxtData(file, callBack)
	return
}

//获取文件数据
func GetTxtData(file *os.File, callBack CallBackDataFunc) (err error) {
	if callBack == nil {
		err = errors.New("callback is nil")
		return
	}
	if file == nil {
		err = errors.New("os.File is nil")
		return
	}

	//读取100个字节判断文件编码
	var data [100]byte
	_, err = file.Read(data[:])
	//TODO
	if err != nil {
		if err == io.EOF {
			fi, err := file.Stat()
			if err != nil {
				return err
			}
			//文件小于100个字节，直接当utf8字节处理
			if fi.Size() < 100 && fi.Size() > 0 {
				filedata := make([]byte, fi.Size())
				_, err = file.Read(filedata[:])
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
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return
	}

	//gbk、csv文件
	if strings.Contains(charset.Charset, "GB") || strings.Contains(charset.Charset, "KOI8-R") ||
		strings.Contains(charset.Charset, "UTF-16LE") || strings.Contains(charset.Charset, "UTF-16BE") {
		//文件流编码解析
		trancreader := decoder.EncodeReader(charset.Charset, file)
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
		reader := bufio.NewReader(file)
		//行读取
		for {
			lineNum++
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				continue
			}

			if !callBack(string(line), "第"+strconv.Itoa(lineNum)+"行") {
				return nil
			}
		}
	}

	return
}
