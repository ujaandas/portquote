package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Logger struct {
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
}

func NewLogger() *Logger {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	origin := filepath.Base(file) + ":" + fmt.Sprint(line)

	return &Logger{
		debug: log.New(os.Stdout, fmt.Sprintf("\033[37m[%s] DEBUG: ", origin), log.LstdFlags),
		info:  log.New(os.Stdout, fmt.Sprintf("\033[34m[%s] INFO: ", origin), log.LstdFlags),
		warn:  log.New(os.Stdout, fmt.Sprintf("\033[33m[%s] WARN: ", origin), log.LstdFlags),
		err:   log.New(os.Stdout, fmt.Sprintf("\033[31m[%s] ERROR:: ", origin), log.LstdFlags),
	}
}

func (l *Logger) Debug(v ...interface{}) {
	l.info.Println(v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.info.Println(v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.info.Println(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.info.Println(v...)
}
