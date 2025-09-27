package http

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/cilginc/webserver-go/internal/logging"
)

type ConnContext struct {
	Conn   net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer
	Logger logging.Logger
}

type Header map[string][]string

func (h Header) Get(k string) string {
	v := h[httpCanonical(k)]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

func httpCanonical(s string) string { return textproto.CanonicalMIMEHeaderKey(s) }

type Request struct {
	Method  string
	Path    string
	Proto   string
	Headers Header
	Body    []byte
}

func ReadRequest(r *bufio.Reader) (*Request, error) {
	line, err := readLine(r)
	if err != nil {
		return nil, err
	}
	parts := strings.SplitN(line, " ", 3)
	if len(parts) != 3 {
		return nil, errors.New("malformed request line")
	}
	method, path, proto := parts[0], parts[1], parts[2]

	h := make(Header)
	tp := textproto.NewReader(r)
	m, err := tp.ReadMIMEHeader()
	if err != nil {
		return nil, err
	}
	for k, vals := range m {
		h[httpCanonical(k)] = vals
	}

	var body []byte
	if cl := h.Get("Content-Length"); cl != "" {
		n, _ := strconv.Atoi(cl)
		if n > 0 {
			body = make([]byte, n)
			if _, err := io.ReadFull(r, body); err != nil {
				return nil, err
			}
		}
	}

	return &Request{Method: method, Path: path, Proto: proto, Headers: h, Body: body}, nil
}

func readLine(r *bufio.Reader) (string, error) {
	var buf bytes.Buffer
	for {
		b, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		if b == '\r' {
			next, err := r.ReadByte()
			if err != nil {
				return "", err
			}
			if next != '\n' {
				return "", errors.New("expected LF")
			}
			break
		}
		if b == '\n' {
			break
		}
		buf.WriteByte(b)
	}
	return buf.String(), nil
}
