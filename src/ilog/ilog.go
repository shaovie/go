package ilog

import (
	"log"
	"os"
	"path"
	"sync"
)

var (
	iLog *log.Logger
	lmtx sync.Mutex
	lf   *os.File
)

func Open(logFile string) error {
	lmtx.Lock()
	defer lmtx.Unlock()

	if err := os.MkdirAll(path.Dir(logFile), 0755); err != nil {
		return err
	}
	lf, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	iLog = log.New(lf, "", log.Ldate|log.Lmicroseconds)
	return nil
}

func ReOpen(logFile string) error {
	if lf != nil {
		lf.Close()
	}
	lf, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	iLog = log.New(lf, "", log.Ldate|log.Lmicroseconds)
	return nil
}

func Rinfo(format string, v ...interface{}) {
	lmtx.Lock()
	defer lmtx.Unlock()

	format = "Rinfo - " + format
	if iLog != nil {
		iLog.Printf(format, v...)
	}
}

func Error(format string, v ...interface{}) {
	lmtx.Lock()
	defer lmtx.Unlock()

	format = "Error - " + format
	if iLog != nil {
		iLog.Printf(format, v...)
	}
}
