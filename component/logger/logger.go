package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

var logger *log.Logger

var errorLogger *log.Logger

func init() {
	errFile, err := os.OpenFile("/Users/tongsiqi/go/src/github.com/tongsq/gin-example/logs/errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("打开日志文件失败：", err)
	}
	//errorLogger = log.New(io.MultiWriter(os.Stderr, errFile), "Error:", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(errFile, "Error:", log.Ldate|log.Ltime)
	logger = log.New(os.Stdout, "Info:", log.Ldate|log.Ltime)
}

func Info(v ...interface{}) {
	logger.Println(getFileInfo(), v)
}

func Success(v ...interface{}) {
	logger.Println("\x1b[0;32m", getFileInfo(), v, "\x1b[0m")
}

func Warning(v ...interface{}) {
	logger.Println("\x1b[0;33m", getFileInfo(), v, "\x1b[0m")
}

func Error(v ...interface{}) {
	logger.Println("\x1b[0;31m", getFileInfo(), v, "\x1b[0m")
	errorLogger.Println(v)
}

func getFileInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	return fmt.Sprintf("%s:%d:", file, line)
}
