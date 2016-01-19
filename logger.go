package main

import (
	"fmt"
	"log"
	"os"
)

type Logger struct {
	logFp *os.File
	level int
}

const (
	ERR   = 0
	INFO  = 1
	DEBUG = 2
)
const (
	SDKDERROR = "[SdkdError]"
	SDKDINFO  = "[SdkdInfo]"
	SDKDDEBUG = "[SdkdDebug]"
)

func (l *Logger) Init(file string, level int) {
	var fp *os.File
	if file != "" {
		var err error
		fp, err = os.Create(file)
		if err != nil {
			log.Fatalf("Unable to open log file %v", err)
		}
	} else {
		fp = os.Stdout
	}
	l.logFp = fp
	l.level = level

}

func (l *Logger) Close() {
	l.logFp.Close()
}

func (l *Logger) printf(prefix string, format string, v ...interface{}) {
	msg := fmt.Sprintf(prefix+format+"\n", v...)
	l.logFp.WriteString(msg)
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.level >= INFO {
		l.printf(SDKDINFO, format, v...)
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level >= DEBUG {
		l.printf(SDKDDEBUG, format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l.level >= ERR {
		l.printf(SDKDERROR, format, v...)
	}
}

func (l *Logger) Output(message string) error {
	l.printf("", message)
	return nil
}
