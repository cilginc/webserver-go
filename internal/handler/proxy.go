package handler

import (
	"io"
	"net"
	"net/url"
	"time"

	"github.com/cilginc/webserver-go/internal/http"
	"github.com/cilginc/webserver-go/internal/router"
)

func ReverseProxy(target string) router.Handler {
	u, _ := url.Parse(target)
	return router.HandlerFunc(func(ctx *http.ConnContext, req *http.Request) error {
		upConn, err := net.DialTimeout("tcp", u.Host, 3*time.Second)
		if err != nil {
			return err
		}
		defer upConn.Close()
		upConn.Write([]byte(req.Method + " " + req.Path + " " + req.Proto + "\r\n"))
		io.Copy(upConn, ctx.Reader)
		io.Copy(ctx.Conn, upConn)
		return nil
	})
}
