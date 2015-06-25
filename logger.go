package main

import (
    "fmt"
)


typedef Logger struct {
    logFp *os.File
    level int
}

const (
    ERR = 0 
    INFO = 1
    DEBUG = 2
)

func (l *Logger) Init(file String, level int)  {
    if fp, err := os.OpenFile(file, O_RDWR, ModeAppend); err != nil {
        log.Errorf("Unable to open log file %v", err)
    } else {
        l.logFp = fp
        l.level = level
    }
}

func (l *Logger) Close() {
    l.logFp.Close()
}

func (l *Logger) printf(prefix string, format string, v ...interface{}) {
    msg := fmt.Sprintf(prefix+format, v...)
    fmt.Println(msg);
}

func (l *Logger) Info(format string, v ...interface{}) {
    if (level <= INFO) {
        l.printf("INFO:", format, v...)
    }
}

func (l *Logger) Debug(format string, v ...interface{}) {
    if (level <= DEBUG ) {
        l.printf("DEBUG:", format, v...)
    }
}

func (l *logger) Error(format string, v ...interface{}) {
    l.printf("ERROR:", format, v...)
}
