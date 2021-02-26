package network

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type Response struct {
	FirstLine string
	Protocol string
	Headers map[string]string
	Body []byte
}

func NewResponse() *Response {
	return &Response{
		Headers: make(map[string]string),
	}
}

func (r *Response) Parse(reader io.Reader) error {
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
			r.FirstLine = string(line)
			firstLineParsed = true
			continue
		}

		if err := r.parseHeader(line); err != nil {
			return err
		}
	}

	return r.getBody(bReader)
}



func (r *Response) parseHeader(buf []byte) error {
	idx := bytes.Index(buf, []byte(":"))
	if idx < 0 {
		return fmt.Errorf("can't parse header: " + string(buf))
	}

	r.Headers[string(buf[:idx])] = string(bytes.TrimSpace(buf[idx+1:]))
	return nil
}

func (r *Response) getBody(bRdr *bufio.Reader) error {
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

func (r *Response) Bytes() []byte {
	buff := bytes.Buffer{}
	buff.WriteString(fmt.Sprintf("%s\r\n", r.FirstLine))
	for h, val := range r.Headers {
		buff.WriteString(fmt.Sprintf("%s: %s \r\n", h, val))
	}
	buff.WriteString("\r\n")
	buff.Write(r.Body)

	return buff.Bytes()
}
