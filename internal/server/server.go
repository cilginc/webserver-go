package server

import (
	"errors"
	"net"
	"path/filepath"
	"sync"
	"time"

	"github.com/cilginc/webserver-go/internal/handler"
	"github.com/cilginc/webserver-go/internal/http"
	"github.com/cilginc/webserver-go/internal/logging"
	"github.com/cilginc/webserver-go/internal/router"
)

type Config struct {
	Addr         string
	Root         string
	Workers      int
	Logger       logging.Logger
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Server struct {
	cfg  *Config
	ln   net.Listener
	wg   sync.WaitGroup
	quit chan struct{}
}

func New(cfg *Config) *Server {
	return &Server{cfg: cfg, quit: make(chan struct{})}
}

func (s *Server) ListenAndServe() error {
	if s.cfg == nil {
		return errors.New("nil config")
	}

	ln, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		return err
	}
	s.ln = ln

	r := router.New()
	absRoot, _ := filepath.Abs(s.cfg.Root)
	r.Handle("GET", "/", handler.Static(absRoot))
	r.Handle("GET", "/*", handler.Static(absRoot))

	conns := make(chan net.Conn, s.cfg.Workers*4)
	for i := 0; i < s.cfg.Workers; i++ {
		s.wg.Add(1)
		go s.worker(conns, r)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				select {
				case <-s.quit:
					return
				default:
					s.cfg.Logger.Errorf("accept error: %v", err)
					continue
				}
			}
			conns <- conn
		}
	}()

	select {}
}

func (s *Server) worker(conns <-chan net.Conn, r *router.Router) {
	defer s.wg.Done()
	for conn := range conns {
		conn.SetReadDeadline(time.Now().Add(s.cfg.ReadTimeout))
		conn.SetWriteDeadline(time.Now().Add(s.cfg.WriteTimeout))
		go func(c net.Conn) {
			defer c.Close()
			ctx := &http.ConnContext{Conn: c, Logger: s.cfg.Logger}
			if err := HandleConnection(ctx, r); err != nil {
				s.cfg.Logger.Errorf("connection error: %v", err)
			}
		}(conn)
	}
}
