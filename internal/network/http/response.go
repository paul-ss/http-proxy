package http

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

type Response struct {
	Status int
	Message string
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
			r.parseFirstLine(line)
			firstLineParsed = true
			continue
		}

		if err := r.parseHeader(line); err != nil {
			return err
		}
	}

	return r.getBody(bReader)
}


func (r *Response) parseFirstLine(buf []byte) error {
	fields := bytes.Fields(buf)
	if len(fields) < 2 {
		return fmt.Errorf("can't parse first line: " + string(buf))
	}

	r.Protocol = string(fields[0])

	status, err := strconv.Atoi(string(fields[1]))
	if err != nil {
		return fmt.Errorf("can't parse status: " + string(fields[1]))
	}

	r.Status = status
	r.Message = string(bytes.Join(fields[2:], []byte(" ")))

	return nil
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
	buff.WriteString(fmt.Sprintf("%s %d %s\r\n", r.Protocol, r.Status, r.Message))
	for h, val := range r.Headers {
		buff.WriteString(fmt.Sprintf("%s: %s \r\n", h, val))
	}
	buff.WriteString("\r\n")
	buff.Write(r.Body)

	return buff.Bytes()
}
