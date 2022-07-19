package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"app/pkg/file"
)

// Level :
type Level int

// F :
var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2

	logger     *log.Logger
	logPrefix  = ""
	levelFlags = []string{"DEBUG", "INFO", "QUERY", "WARN", "ERROR", "FATAL"}
	eFlag      = ""
	eFunc      = ""
	eFile      = ""
	eLine      = -1
)

// DEBUG :
const (
	DEBUG Level = iota
	INFO
	QUERY
	WARNING
	ERROR
	FATAL
)

// Setup :
func Setup() {
	now := time.Now()
	var err error
	filePath := getLogFilePath()
	fileName := getLogFileName()
	F, err = file.MustOpen(fileName, filePath)

	if err != nil {
		log.Fatalf("logging.Setup err: %v", err)
	}

	logger = log.New(F, DefaultPrefix, log.LstdFlags)
	timeSpent := time.Since(now)
	log.Printf("Config logging is ready in %v", timeSpent)
}

// Debug :
func Debug(user string, v ...interface{}) {
	var audit AuditLog
	setPrefix(&audit, DEBUG)
	log.Println(v...)
	logger.Println(v...)
	audit.Message = fmt.Sprintf("%v", v)
	go audit.saveAudit()
}

// Info :
func Info(user string, v ...interface{}) {
	var audit AuditLog
	setPrefix(&audit, INFO)
	log.Println(v...)
	logger.Println(v...)
	audit.Message = fmt.Sprintf("%v", v)
	go audit.saveAudit()
}

// Query :
func Query(user string, v ...interface{}) {
	var audit AuditLog
	setPrefix(&audit, QUERY)
	log.Println(v...)
	logger.Println(v...)
	audit.Message = fmt.Sprintf("%v", v)
	go audit.saveAudit()
}

// Warn :
func Warn(user string, v ...interface{}) {
	var audit AuditLog
	setPrefix(&audit, WARNING)
	log.Println(v...)
	logger.Println(v...)
	audit.Message = fmt.Sprintf("%v", v)
	go audit.saveAudit()
}

// Error :
func Error(user string, v ...interface{}) {
	var audit AuditLog
	setPrefix(&audit, ERROR)
	log.Println(v...)
	logger.Println(v...)
	audit.Message = fmt.Sprintf("%v", v)
	go audit.saveAudit()
}

// Fatal :
func Fatal(user string, v ...interface{}) {
	var audit AuditLog
	setPrefix(&audit, FATAL)
	log.Println(v...)
	logger.Println(v...)
}

func setPrefix(audit *AuditLog, level Level) {
	// loc, err := time.LoadLocation("Asia/Jakarta")
	// if err != nil {
	// 	log.Print(err)
	// }
	t := time.Now()
	function, file, line, ok := runtime.Caller(DefaultCallerDepth)
	audit.Level = levelFlags[level]
	audit.UUID = "SYS"
	audit.FuncName = ""
	audit.FileName = filepath.Base(file)
	audit.Line = line
	audit.Time = fmt.Sprintf("%s", t.Format("2006-01-02 15:04:05"))
	if ok {
		s := strings.Split(runtime.FuncForPC(function).Name(), ".")
		_, fn := s[0], s[1]
		logPrefix = fmt.Sprintf("[%s][SYS][%s][%s:%d]", levelFlags[level], fn, filepath.Base(file), line)
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
