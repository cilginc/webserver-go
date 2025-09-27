package logging

import (
	"log"
	"os"
)

type Logger interface {
	Infof(format string, v ...any)
	Warnf(format string, v ...any)
	Errorf(format string, v ...any)
	Fatalf(format string, v ...any)
}

type stdLogger struct {
	l *log.Logger
}

func NewStdLogger() Logger {
	return &stdLogger{
		l: log.New(os.Stdout, "[go-nginx] ", log.LstdFlags),
	}
}

func (s *stdLogger) Infof(format string, v ...any) {
	s.l.Printf("[INFO] "+format, v...)
}

func (s *stdLogger) Warnf(format string, v ...any) {
	s.l.Printf("[WARN] "+format, v...)
}

func (s *stdLogger) Errorf(format string, v ...any) {
	s.l.Printf("[ERROR] "+format, v...)
}

func (s *stdLogger) Fatalf(format string, v ...any) {
	s.l.Fatalf("[FATAL] "+format, v...)
}
