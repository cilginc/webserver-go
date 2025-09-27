package http

import (
	"bufio"
	"strings"
	"testing"
)

func mockReader(s string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(s))
}

func TestReadRequest_SimpleGet(t *testing.T) {
	raw := "GET /index.html HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"User-Agent: TestClient\r\n" +
		"\r\n"

	req, err := ReadRequest(mockReader(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Method != "GET" {
		t.Errorf("expected GET, got %q", req.Method)
	}
	if req.Path != "/index.html" {
		t.Errorf("expected /index.html, got %q", req.Path)
	}
	if req.Proto != "HTTP/1.1" {
		t.Errorf("expected HTTP/1.1, got %q", req.Proto)
	}

	if got := req.Headers.Get("Host"); got != "example.com" {
		t.Errorf("expected Host=example.com, got %q", got)
	}
	if got := req.Headers.Get("User-Agent"); got != "TestClient" {
		t.Errorf("expected User-Agent=TestClient, got %q", got)
	}
	if len(req.Body) != 0 {
		t.Errorf("expected empty body, got %q", string(req.Body))
	}
}

func TestReadRequest_PostWithBody(t *testing.T) {
	raw := "POST /submit HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"Content-Length: 11\r\n" +
		"\r\n" +
		"hello=world"

	req, err := ReadRequest(mockReader(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Method != "POST" {
		t.Errorf("expected POST, got %q", req.Method)
	}
	if req.Path != "/submit" {
		t.Errorf("expected /submit, got %q", req.Path)
	}

	if got := req.Headers.Get("Content-Length"); got != "11" {
		t.Errorf("expected Content-Length=11, got %q", got)
	}
	if string(req.Body) != "hello=world" {
		t.Errorf("expected body 'hello=world', got %q", string(req.Body))
	}
}

func TestReadRequest_BadRequestLine(t *testing.T) {
	raw := "INVALIDREQUEST\r\n\r\n"
	_, err := ReadRequest(mockReader(raw))
	if err == nil {
		t.Fatal("expected error for malformed request line, got nil")
	}
}

func TestReadRequest_BadHeader(t *testing.T) {
	raw := "GET / HTTP/1.1\r\n" +
		"BadHeaderWithoutColon\r\n" +
		"\r\n"
	_, err := ReadRequest(mockReader(raw))
	if err == nil {
		t.Fatal("expected error for malformed header, got nil")
	}
}
