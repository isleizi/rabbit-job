package logg

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	traces  *log.Logger // 记录所有日志
	info    *log.Logger // 重要的信息
	warning *log.Logger // 需要注意的信息
	error   *log.Logger // 非常严重的问题
)

func Trace(v ...interface{}) {
	traces = log.New(ioutil.Discard, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	traces.Println(v)
}

func Info(v ...interface{}) {
	info = log.New(ioutil.Discard, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	info.Println(v)
}

func Warning(v ...interface{}) {
	info = log.New(ioutil.Discard, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	info.Println(v)
}

func Error(v ...interface{}) {
	file, err := os.OpenFile("errors.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}
	error = log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	error.Println(v)
}
