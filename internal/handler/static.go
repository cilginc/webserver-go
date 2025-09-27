package handler

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/cilginc/webserver-go/internal/http"
	"github.com/cilginc/webserver-go/internal/router"
)

type staticHandler struct {
	root string
}

func Static(root string) router.Handler {
	return &staticHandler{root: root}
}

func (s *staticHandler) ServeHTTP(ctx *http.ConnContext, req *http.Request) error {
	p := req.Path
	if i := strings.Index(p, "?"); i >= 0 {
		p = p[:i]
	}
	if p == "/" {
		p = "/index.html"
	}
	fp := filepath.Join(s.root, filepath.Clean(p))
	if !strings.HasPrefix(fp, s.root) {
		http.WriteSimpleResponse(ctx.Writer, 403, "Forbidden")
		return nil
	}
	fi, err := os.Stat(fp)
	if err != nil {
		if os.IsNotExist(err) {
			http.WriteSimpleResponse(ctx.Writer, 404, "Not Found")
			return nil
		}
		http.WriteSimpleResponse(ctx.Writer, 500, "Internal")
		return err
	}
	if fi.IsDir() {
		idx := filepath.Join(fp, "index.html")
		if st, e := os.Stat(idx); e == nil && !st.IsDir() {
			fp = idx
		} else {
			http.WriteSimpleResponse(ctx.Writer, 404, "Not Found")
			return nil
		}

	}

	f, err := os.Open(fp)
	if err != nil {
		http.WriteSimpleResponse(ctx.Writer, 500, "Internal")
		return err
	}
	defer f.Close()

	ct := mime.TypeByExtension(filepath.Ext(fp))
	if ct == "" {
		ct = "application/octet-stream"
	}

	fmt.Fprintf(ctx.Writer, "HTTP/1.1 200 OK\r\n")
	fmt.Fprintf(ctx.Writer, "Content-Length: %d\r\n", fi.Size())
	fmt.Fprintf(ctx.Writer, "Content-Type: %s\r\n", ct)
	fmt.Fprintf(ctx.Writer, "Connection: keep-alive\r\n")
	fmt.Fprintf(ctx.Writer, "\r\n")
	ctx.Writer.Flush()

	buf := make([]byte, 32*1024)
	for {
		n, err := f.Read(buf)
		if n > 0 {
			if _, werr := ctx.Conn.Write(buf[:n]); werr != nil {
				return werr
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}
