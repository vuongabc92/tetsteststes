package log

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

const ErrorLog string = "error"
const DebugLog string = "debug"
const InfoLog string = "info"

type Log struct {
	requestID string
	error     *log.Logger
	info      *log.Logger
	debug     *log.Logger
	disable   bool
}

func NewLog(logFiles map[string]string) *Log {
	var logger Log
	if f, has := logFiles[ErrorLog]; has {
		f, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Can not open log file: %v", err)
		}

		logger.error = log.New(f, fmt.Sprintf("%s: ", ErrorLog), log.LstdFlags)
	}

	if f, has := logFiles[DebugLog]; has {
		f, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Can not open log file: %v", err)
		}

		logger.debug = log.New(f, fmt.Sprintf("%s: ", DebugLog), log.LstdFlags)
	}

	if f, has := logFiles[InfoLog]; has {
		f, err := os.OpenFile(f, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Can not open log file: %v", err)
		}

		logger.info = log.New(f, fmt.Sprintf("%s: ", InfoLog), log.LstdFlags)
	}

	return &logger
}

func (l *Log) Error(v ...interface{}) {
	if !l.disable && l.error != nil {
		debugStack := fmt.Sprintf("\n%s", debug.Stack())
		l.error.Println(v, l.getRequestId(), debugStack)
		log.Println(v, l.getRequestId(), debugStack)
	}
}

func (l *Log) Errorf(format string, agrs ...interface{}) {
	if !l.disable && l.error != nil {
		debugStack := fmt.Sprintf("\n%s", debug.Stack())
		l.error.Println(fmt.Sprintf(format, agrs...), l.getRequestId(), debugStack)
		log.Println(fmt.Sprintf(format, agrs...), l.getRequestId(), debugStack)
	}
}

func (l *Log) Info(v ...interface{}) {
	if !l.disable && l.info != nil {
		l.info.Println(v, l.getRequestId())
		log.Println(v, l.getRequestId())
	}
}

func (l *Log) Infof(format string, agrs ...interface{}) {
	if !l.disable && l.info != nil {
		l.info.Println(fmt.Sprintf(format, agrs...), l.getRequestId())
		log.Println(fmt.Sprintf(format, agrs...), l.getRequestId())
	}
}

func (l *Log) Debug(v ...interface{}) {
	if !l.disable && l.debug != nil {
		debugStack := fmt.Sprintf("\n%s", debug.Stack())
		l.debug.Println(v, l.getRequestId(), debugStack)
		log.Println(v, l.getRequestId(), debugStack)
	}
}

func (l *Log) Disable() {
	l.disable = true
}

func (l *Log) SetRequestID(requestId string) {
	l.requestID = requestId
}

func (l *Log) getRequestId() string {
	if len(l.requestID) > 0 {
		return fmt.Sprintf("[Request ID: %s]", l.requestID)
	}

	return ""
}
