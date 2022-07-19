package logging

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Logger :
type Logger struct {
	// UUID string `json:"uuid,omitempty"`
	// E    echo.Context
}

// Debug :
func (l *Logger) Debug(v ...interface{}) {
	var audit AuditLog
	l.setUserLogPrefix(&audit, DEBUG)
	log.Println(v...)
	logger.Println(v...)
	audit.Message = fmt.Sprintf("%v", v)
	audit.saveAudit()
}

// Info :
func (l *Logger) Info(v ...interface{}) {
	var audit AuditLog
	l.setUserLogPrefix(&audit, INFO)
	log.Println(v...)
	logger.Println(v...)
	audit.Message = fmt.Sprintf("%v", v)
	audit.saveAudit()
}

// Query :
func (l *Logger) Query(v ...interface{}) {
	var audit AuditLog
	l.setUserLogPrefix(&audit, QUERY)
	log.Println(v...)
	logger.Println(v...)
	audit.Message = fmt.Sprintf("%v", v)
	audit.saveAudit()
}

// Warn :
func (l *Logger) Warn(v ...interface{}) {
	var audit AuditLog
	l.setUserLogPrefix(&audit, WARNING)
	log.Println(v...)
	logger.Println(v...)
	audit.Message = fmt.Sprintf("%v", v)
	audit.saveAudit()
}

// Error :
func (l *Logger) Error(v ...interface{}) {
	var audit AuditLog
	l.setUserLogPrefix(&audit, ERROR)
	log.Println(v...)
	logger.Println(v...)
	audit.Message = fmt.Sprintf("%v", v)
	audit.saveAudit()
}

// Fatal :
func (l *Logger) Fatal(v ...interface{}) {
	var audit AuditLog
	l.setUserLogPrefix(&audit, FATAL)
	log.Println(v...)
	logger.Fatalln(v...)
}

func (l *Logger) setUserLogPrefix(audit *AuditLog, level Level) {
	// loc, err := time.LoadLocation("Asia/Jakarta")
	// if err != nil {
	// 	log.Print(err)
	// }
	t := time.Now()
	function, file, line, ok := runtime.Caller(DefaultCallerDepth)
	audit.Level = levelFlags[level]
	// audit.UUID = l.UUID
	audit.FuncName = ""
	audit.FileName = filepath.Base(file)
	audit.Line = line
	audit.Time = fmt.Sprintf("%s", t.Format("2006-01-02 15:04:05"))

	if ok {
		s := strings.Split(runtime.FuncForPC(function).Name(), ".")
		_, fn := s[0], s[1]
		// logPrefix = fmt.Sprintf("[%s][%s][%s][%s:%d]", levelFlags[level], l.UUID, fn, filepath.Base(file), line)
		logPrefix = fmt.Sprintf("[%s][%s][%s:%d]", levelFlags[level], fn, filepath.Base(file), line)
		eFlag = levelFlags[level]
		eFunc = fn
		eFile = filepath.Base(file)
		eLine = line
		audit.FuncName = fn
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}
	logger.SetPrefix(logPrefix)
}
