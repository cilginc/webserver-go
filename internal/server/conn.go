package server

import (
	"bufio"
	"io"
	"strings"

	"github.com/cilginc/webserver-go/internal/http"
	"github.com/cilginc/webserver-go/internal/router"
)

func HandleConnection(ctx *http.ConnContext, r *router.Router) error {
	if ctx.Reader == nil {
		ctx.Reader = bufio.NewReader(ctx.Conn)
	}
	if ctx.Writer == nil {
		ctx.Writer = bufio.NewWriter(ctx.Conn)
	}

	for {
		req, err := http.ReadRequest(ctx.Reader)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		h := r.Match(req.Method, req.Path)
		if h == nil {
			http.WriteSimpleResponse(ctx.Writer, 404, "Not Found")
			ctx.Writer.Flush()
			continue
		}

		if err := h.ServeHTTP(ctx, req); err != nil {
			ctx.Logger.Errorf("handler error: %v", err)
		}
		ctx.Writer.Flush()

		if strings.ToLower(req.Headers.Get("connection")) == "close" {
			return nil
		}
	}
}
