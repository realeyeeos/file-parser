// Package mail implements a parser for electronic mail messages as specified
// in RFC5322.
//
// We allow both CRLF and LF to be used in the input, possibly mixed.
package emlparser

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"net/textproto"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/realeyeeos/file-parser/decoder"
)

var benc = base64.URLEncoding

func mkId(s []byte) string {
	h := sha1.New()
	h.Write(s)
	hash := h.Sum(nil)
	ed := benc.EncodeToString(hash)
	return ed[0:20]
}

type HeaderInfo struct {
	FullHeaders []Header // all headers
	OptHeaders  []Header // unprocessed headers

	MessageId   string
	Id          string
	Date        time.Time
	From        []Address
	Sender      Address
	ReplyTo     []Address
	To          []Address
	Cc          []Address
	Bcc         []Address
	Subject     string
	Comments    []string
	Keywords    []string
	ContentType string

	InReply    []string
	References []string
}

type Message struct {
	HeaderInfo
	Body        []byte
	Text        []byte
	Html        []byte
	Attachments []Attachment
	//除了附件外的数据（图片等）
	OtherAttachments []Attachment
	Parts            []Part
}

type Attachment struct {
	Filename string
	Data     []byte
}

type Header struct {
	Key, Value string
}

//savefile:true-读取附件 false-不读取附件
func Parse(file *os.File, savefile bool) (m Message, e error) {
	reader := bufio.NewReader(file)
	r, e := ParseRaw(reader)
	if e != nil {
		return
	}
	return Process(r, savefile)
}

//savefile:true-读取附件 false-不读取附件
func Process(r RawMessage, savefile bool) (m Message, e error) {
	m.FullHeaders = []Header{}
	m.OptHeaders = []Header{}
	for _, rh := range r.RawHeaders {
		h := Header{string(rh.Key), string(rh.Value)}
		m.FullHeaders = append(m.FullHeaders, h)
		switch string(rh.Key) {
		case `Content-Type`:
			m.ContentType = string(rh.Value)
		case `Message-ID`:
			v := bytes.Trim(rh.Value, `<>`)
			m.MessageId = string(v)
			m.Id = mkId(v)
		case `In-Reply-To`:
			ids := strings.Fields(string(rh.Value))
			for _, id := range ids {
				m.InReply = append(m.InReply, strings.Trim(id, `<> `))
			}
		case `References`:
			ids := strings.Fields(string(rh.Value))
			for _, id := range ids {
				m.References = append(m.References, strings.Trim(id, `<> `))
			}
		case `Date`:
			m.Date = ParseDate(string(rh.Value))
		case `From`:
			m.From, e = parseAddressList(rh.Value)
		case `Sender`:
			m.Sender, e = ParseAddress(rh.Value)
		case `Reply-To`:
			m.ReplyTo, e = parseAddressList(rh.Value)
		case `To`:
			m.To, e = parseAddressList(rh.Value)
		case `Cc`:
			m.Cc, e = parseAddressList(rh.Value)
		case `Bcc`:
			m.Bcc, e = parseAddressList(rh.Value)
		case `Subject`:
			subject, err := decoder.Parse(rh.Value)
			if err == nil {
				m.Subject = string(subject)
			}

		case `Comments`:
			m.Comments = append(m.Comments, string(rh.Value))
		case `Keywords`:
			ks := strings.Split(string(rh.Value), ",")
			for _, k := range ks {
				m.Keywords = append(m.Keywords, strings.TrimSpace(k))
			}
		default:
			m.OptHeaders = append(m.OptHeaders, h)
		}
		if e != nil {
			return
		}
	}
	if m.Sender == nil && len(m.From) > 0 {
		m.Sender = m.From[0]
	}

	if m.ContentType != `` {
		parts, er := parseBody(m.ContentType, r.Reader, savefile)
		if er != nil && er != io.EOF {
			//单数据的需要特殊处理
			if er.Error() == "multipart specified without boundary" {
				contenttype := regexp.MustCompile("(?is)charset=(.*)").FindStringSubmatch(m.ContentType)
				charset := "UTF-8"
				if len(contenttype) > 1 {
					charset = contenttype[1]
				}

				//删除\r\n
				var data []byte
				for {
					line, err := r.Reader.ReadBytes('\n')
					if err != nil {
						break
					}
					data = append(data, bytes.Replace(line, CRLF, nil, -1)...)
				}

				//data := bytes.Replace(r.Body, CRLF, nil, -1)
				part := Part{Type: m.ContentType, Charset: charset, Data: data}
				part.Headers = make(map[string][]string, len(m.OptHeaders))
				for _, v := range m.OptHeaders {
					part.Headers[v.Key] = []string{v.Value}
				}

				parts = append(parts, part)

			} else {
				e = er
				return
			}
		}

		for _, part := range parts {
			switch {
			case strings.Contains(part.Type, "text/plain"):
				if encoding, ok := part.Headers["Content-Transfer-Encoding"]; ok {
					m.Text, _ = decoder.Decode(encoding[0], part.Data)
				} else {
					m.Text = part.Data
				}

				m.Text, _ = decoder.UTF8(part.Charset, m.Text)

			case strings.Contains(part.Type, "text/html"):
				if encoding, ok := part.Headers["Content-Transfer-Encoding"]; ok {
					m.Html, _ = decoder.Decode(encoding[0], part.Data)
				} else {
					m.Html = part.Data
				}
				m.Html, _ = decoder.UTF8(part.Charset, m.Html)

			default:
				if !savefile {
					continue
				}
				if cd, ok := part.Headers["Content-Disposition"]; ok {
					//if strings.Contains(cd[0], "attachment") {
					filename := regexp.MustCompile("(?msi)name=\"(.*?)\"").FindStringSubmatch(cd[0]) //.FindString(cd[0])
					if len(filename) < 2 {
						//fmt.Println("failed get filename from header content-disposition")
						break
					}

					dfilename, err := decoder.Parse([]byte(filename[1]))
					if err != nil {
						//fmt.Println("Failed decode filename of attachment", err)
					} else {
						filename[1] = string(dfilename)
					}

					if encoding, ok := part.Headers["Content-Transfer-Encoding"]; ok {
						part.Data, _ = decoder.Decode(encoding[0], part.Data)
					}

					if strings.Contains(cd[0], "attachment") {
						m.Attachments = append(m.Attachments, Attachment{filename[1], part.Data})
					} else {
						m.OtherAttachments = append(m.Attachments, Attachment{filename[1], part.Data})
					}

					//}
				}
			}
		}

		m.Parts = parts
		//m.ContentType = parts[0].Type
		//m.Text, _ = decoder.Decode("B", parts[0].Data)
	} else {
		m.Text = r.Body
	}
	return
}

type RawHeader struct {
	Key, Value []byte
}

type RawMessage struct {
	RawHeaders []RawHeader
	Body       []byte
	Reader     *bufio.Reader
}

func ParseRaw(reader *bufio.Reader) (m RawMessage, e error) {
	// parser states
	m.RawHeaders = []RawHeader{}

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}

		if bytes.Equal(line, CRLF) {
			m.Reader = reader
			break
		}

		if line[0] == ' ' || line[0] == '\t' {
			rawlen := len(m.RawHeaders)
			if rawlen == 0 {
				continue
			}
			value := bytes.Replace(line, CRLF, nil, -1)
			m.RawHeaders[rawlen-1].Value = append(m.RawHeaders[rawlen-1].Value, value...)
		} else {
			r := textproto.NewReader(bufio.NewReader(bytes.NewReader(line)))
			header, err := r.ReadMIMEHeader()
			if err != nil && err != io.EOF {
				continue
			}

			for k, v := range header {
				value := bytes.Replace([]byte(v[0]), CRLF, nil, -1)
				hdr := RawHeader{[]byte(k), value}
				m.RawHeaders = append(m.RawHeaders, hdr)
			}
		}
	}
	return
}
