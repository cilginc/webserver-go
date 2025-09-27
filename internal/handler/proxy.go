package handler

import (
	"bufio"
	"fmt"
	"io"
	"time"

	"github.com/cilginc/webserver-go/internal/http"
	"github.com/cilginc/webserver-go/internal/router"
)

type DialFunc func(network, address string, timeout time.Duration) (io.ReadWriteCloser, error)

func ReverseProxy(targetHost string, dial DialFunc) router.Handler {
	return router.HandlerFunc(func(ctx *http.ConnContext, req *http.Request) error {
		upConn, err := dial("tcp", targetHost, 3*time.Second)
		if err != nil {
			return fmt.Errorf("proxy dial failed: %w", err)
		}
		defer upConn.Close()

		if _, err := upConn.Write([]byte(req.Method + " " + req.Path + " " + req.Proto + "\r\n")); err != nil {
			return err
		}

		if req.Headers != nil {
			for k, values := range req.Headers {
				for _, v := range values {
					line := fmt.Sprintf("%s: %s\r\n", k, v)
					if _, err := upConn.Write([]byte(line)); err != nil {
						return err
					}
				}
			}
		}
		_, _ = upConn.Write([]byte("\r\n"))

		if len(req.Body) > 0 {
			if _, err := upConn.Write(req.Body); err != nil {
				return err
			}
		}

		upReader := bufio.NewReader(upConn)
		_, err = io.Copy(ctx.Writer, upReader)
		if err != nil {
			return fmt.Errorf("proxy response copy failed: %w", err)
		}

		return ctx.Writer.Flush()
	})
}
