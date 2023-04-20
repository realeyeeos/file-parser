// Handle multipart messages.

package emlparser

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"regexp"
	"strings"
)

type Part struct {
	Type    string
	Charset string
	Data    []byte
	Headers map[string][]string
}

const (
	CR = '\r'
	LF = '\n'
)

var CRLF = []byte{CR, LF}

//修改成使用 bufio.Reader 读取数据 add by scl
func parseBody(ct string, reader *bufio.Reader, savefile bool) (parts []Part, err error) {
	_, ps, err := mime.ParseMediaType(ct)
	if err != nil {
		return
	}

	boundary, ok := ps["boundary"]
	if !ok {
		return nil, errors.New("multipart specified without boundary")
	}
	r := multipart.NewReader(reader, boundary)
	p, err := r.NextRawPart()
	for err == nil {
		var subparts []Part
		//modify by scl
		subparts, err = parseBody(p.Header["Content-Type"][0], bufio.NewReader(p), savefile)
		//if err == nil then body have sub multipart, and append him
		if err == nil {
			parts = append(parts, subparts...)
		} else {
			contenttype := regexp.MustCompile("(?is)charset=(.*)").FindStringSubmatch(p.Header["Content-Type"][0])
			if !savefile {
				if !strings.Contains(string(p.Header["Content-Type"][0]), "text/plain") &&
					!strings.Contains(string(p.Header["Content-Type"][0]), "text/html") {
					p, err = r.NextRawPart()
					continue
				}
			}

			charset := "UTF-8"
			if len(contenttype) > 1 {
				charset = contenttype[1]
			}
			data, _ := ioutil.ReadAll(p) // ignore error
			//删除\r\n
			data = bytes.Replace(data, CRLF, nil, -1)
			part := Part{p.Header["Content-Type"][0], charset, data, p.Header}
			parts = append(parts, part)
		}
		p, err = r.NextRawPart()
	}
	if err == io.EOF {
		err = nil
	}
	return
}
