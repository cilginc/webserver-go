package http

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/cilginc/webserver-go/internal/logging"
)

type ConnContext struct {
	Conn   io.ReadWriteCloser
	Reader *bufio.Reader
	Writer *bufio.Writer
	Logger logging.Logger
}

type Header map[string][]string

func (h Header) Get(key string) string {
	values := h[httpCanonical(key)]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

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

	headers, err := readHeaders(r)
	if err != nil {
		return nil, err
	}

	var body []byte
	if cl := headers.Get("Content-Length"); cl != "" {
		n, _ := strconv.Atoi(cl)
		if n > 0 {
			body = make([]byte, n)
			if _, err := io.ReadFull(r, body); err != nil {
				return nil, err
			}
		}
	}

	return &Request{
		Method:  method,
		Path:    path,
		Proto:   proto,
		Headers: headers,
		Body:    body,
	}, nil
}
