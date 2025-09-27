package main

import (
	"flag"
	"os"
	"time"

	"github.com/cilginc/webserver-go/internal/config"
	"github.com/cilginc/webserver-go/internal/logging"
	"github.com/cilginc/webserver-go/internal/server"
)

func main() {
	cfg := config.Load()
	addr := flag.String("listen", cfg.Port, "listen address")
	root := flag.String("root", cfg.StaticDir, "static root dir")
	workers := flag.Int("workers", cfg.WorkerCount, "number of worker goroutines")
	flag.Parse()

	logger := logging.NewStdLogger()
	logger.Infof("starting server on %s, root=%s, workers=%d", *addr, *root, *workers)

	srv := server.New(&server.Config{
		Addr:         *addr,
		Root:         *root,
		Workers:      *workers,
		Logger:       logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatalf("server failed: %v", err)
		os.Exit(1)
	}
}
