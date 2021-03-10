package http

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
)

type Request struct {
	Method string
	Url *url.URL
	Protocol string
	Headers map[string]string
	Body []byte
}

func NewRequest() *Request {
	return &Request{
		Headers: make(map[string]string),
	}
}

func (r *Request) SetFirstLine(method, u, protocol string) error {
	r.Method = method
	r.Protocol = protocol

	uu, err := url.Parse(u)
	if err != nil {
		return fmt.Errorf("can't parse url: " + u)
	}

	r.Url = uu
	return nil
}

func (r *Request) SetHeader(header, value string) {
	r.Headers[header] = value
}


func (r *Request) Parse(reader io.Reader) error {
	bReader := bufio.NewReader(reader)
	started := false
	firstLineParsed := false

	for {
		line, err := bReader.ReadBytes('\n')
		if err != nil {
			return err
		}

		line = bytes.Trim(line, "\n\r")
		if len(line) == 0 {
			if started {
				break
			} else {
				continue
			}
		}
		started = true

		if !firstLineParsed {
			if err := r.parseFirstLine(line); err != nil {
				return err
			}
			firstLineParsed = true
			continue
		}

		if err := r.parseHeader(line); err != nil {
			return err
		}
	}

	return r.getBody(bReader)
}

func (r *Request) parseFirstLine(buf []byte) error {
	fields := bytes.Fields(buf)
	if len(fields) != 3 {
		return fmt.Errorf("can't parse first line: " + string(buf))
	}

	r.Method = string(fields[0])
	r.Protocol = string(fields[2])

	urlField := string(fields[1])
	addedPrefix := false
	if !strings.HasPrefix(urlField, "http") {
		urlField = "http://" + urlField
		addedPrefix = true
	}

	u, err := url.ParseRequestURI(urlField)
	if err != nil {
		return fmt.Errorf("can't parse url line: " + string(buf))
	}

	if addedPrefix {
		u.Scheme = ""
		u.String()
	}
	r.Url = u

	return nil
}

func (r *Request) parseHeader(buf []byte) error {
	idx := bytes.Index(buf, []byte(":"))
	if idx < 0 {
		return fmt.Errorf("can't parse header: " + string(buf))
	}

	r.Headers[string(buf[:idx])] = string(bytes.TrimSpace(buf[idx+1:]))
	return nil
}

func (r *Request) getBody(bRdr *bufio.Reader) error {
	length, ok := r.Headers["Content-Length"]
	if ok {
		l, err := strconv.Atoi(length)
		if err != nil {
			return err
		}

		r.Body  = make([]byte, l)
		n, err := io.ReadFull(bRdr, r.Body)
		r.Body = r.Body[:n]

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Request) Bytes() []byte {
	buff := bytes.Buffer{}
	buff.WriteString(fmt.Sprintf("%s %s %s\r\n", r.Method, r.Url.String(), r.Protocol))
	for h, val := range r.Headers {
		buff.WriteString(fmt.Sprintf("%s: %s \r\n", h, val))
	}
	buff.WriteString("\r\n")
	buff.Write(r.Body)

	return buff.Bytes()
}