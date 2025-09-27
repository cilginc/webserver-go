package main

import (
	"flag"
	"os"
	"time"

	"github.com/cilginc/webserver-go/internal/logging"
	"github.com/cilginc/webserver-go/internal/server"
)

func main() {
	addr := flag.String("listen", ":8080", "listen address")
	root := flag.String("root", "./www", "static root dir")
	workers := flag.Int("workers", 8, "number of worker goroutines")
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
