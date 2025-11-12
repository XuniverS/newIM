package logger

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

// Init 初始化日志
func Init() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info 记录信息日志
func Info(message string) {
	if infoLogger == nil {
		Init()
	}
	infoLogger.Println(message)
}

// Error 记录错误日志
func Error(message string) {
	if errorLogger == nil {
		Init()
	}
	errorLogger.Println(message)
}
