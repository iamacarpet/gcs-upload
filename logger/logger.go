package logger

import (
    "os"
    "fmt"
    "log/syslog"
)

var logWriter *syslog.Writer

func Debugf(format string, v ...interface{}){
    if initWriter() {
        logWriter.Debug(fmt.Sprintf(format, v...))
    }
}

func Infof(format string, v ...interface{}){
    if initWriter() {
        logWriter.Info(fmt.Sprintf(format, v...))
    }
}

func Warnf(format string, v ...interface{}){
    if initWriter() {
        logWriter.Warning(fmt.Sprintf(format, v...))
    }
}

func Errorf(format string, v ...interface{}){
    if initWriter() {
        logWriter.Err(fmt.Sprintf(format, v...))
    }
}

func Criticalf(format string, v ...interface{}){
    if initWriter() {
        logWriter.Crit(fmt.Sprintf(format, v...))
    }
}

func Fatalf(format string, v ...interface{}){
    if initWriter() {
        logWriter.Crit(fmt.Sprintf(format, v...))
    }
    os.Exit(1)
}

func initWriter() (bool) {
    if logWriter == nil {
        var err error
        logWriter, err = syslog.New(syslog.LOG_LOCAL0, "go-app")
        if err != nil {
            return false
        }
    }
    return true
}

func Init(name string) (bool) {
    if logWriter == nil {
        var err error
        logWriter, err = syslog.New(syslog.LOG_LOCAL0, name)
        if err != nil {
            return false
        }
    }
    return true
}
