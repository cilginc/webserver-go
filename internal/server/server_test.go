package server

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/cilginc/webserver-go/internal/http"
)

func TestWriteSimpleResponse(t *testing.T) {
	buf := new(bytes.Buffer)
	w := bufio.NewWriter(buf)
	http.WriteSimpleResponse(w, 200, "hello")
	w.Flush()
	out := buf.String()
	if !bytes.Contains([]byte(out), []byte("200 OK")) {
		t.Errorf("expected status line with 200 OK, got %q", out)
	}
	if !bytes.Contains([]byte(out), []byte("hello")) {
		t.Errorf("expected body 'hello', got %q", out)
	}
}

func TestReadRequest(t *testing.T) {
	raw := "GET /test HTTP/1.1\r\nHost: example.com\r\nContent-Length: 5\r\n\r\nhello"
	r := bufio.NewReader(bytes.NewBufferString(raw))
	req, err := http.ReadRequest(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Method != "GET" || req.Path != "/test" {
		t.Errorf("unexpected request parsed: %+v", req)
	}
	if string(req.Body) != "hello" {
		t.Errorf("expected body 'hello', got %q", string(req.Body))
	}
}
