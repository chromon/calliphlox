package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	// error 日志实例
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	// info 日志实例
	infoLog = log.New(os.Stdout, "\033[32m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers = []*log.Logger{
		errorLog,
		infoLog,
	}
	mutex sync.Mutex
)

// 日志方法
var (
	Error = errorLog.Println
	Errorf = errorLog.Printf
	Info = infoLog.Println
	Infof = infoLog.Printf
)

// 日志层级
const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

// 设置日志层级
func SetLevel(level int) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}

	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}
}