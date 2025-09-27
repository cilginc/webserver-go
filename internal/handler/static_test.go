package handler

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/cilginc/webserver-go/internal/http"
)

func TestStaticHandler(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "index.html")
	content := []byte("<h1>Hello</h1>")
	if err := os.WriteFile(file, content, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	h := Static(dir)
	buf := new(bytes.Buffer)
	ctx := &http.ConnContext{
		Conn:   nopConn{buf},
		Writer: bufio.NewWriter(buf),
	}
	req := &http.Request{Method: "GET", Path: "/"}
	if err := h.ServeHTTP(ctx, req); err != nil {
		t.Fatalf("ServeHTTP error: %v", err)
	}
	ctx.Writer.Flush()
	out := buf.String()
	if !bytes.Contains([]byte(out), []byte("200 OK")) {
		t.Errorf("expected 200 OK, got %q", out)
	}
	if !bytes.Contains([]byte(out), content) {
		t.Errorf("expected body to contain %q, got %q", content, out)
	}
}

type nopConn struct{ *bytes.Buffer }

func (n nopConn) Read(b []byte) (int, error)         { return 0, io.EOF }
func (n nopConn) Write(b []byte) (int, error)        { return n.Buffer.Write(b) }
func (n nopConn) Close() error                       { return nil }
func (n nopConn) LocalAddr() net.Addr                { return nil }
func (n nopConn) RemoteAddr() net.Addr               { return nil }
func (n nopConn) SetDeadline(t time.Time) error      { return nil }
func (n nopConn) SetReadDeadline(t time.Time) error  { return nil }
func (n nopConn) SetWriteDeadline(t time.Time) error { return nil }
