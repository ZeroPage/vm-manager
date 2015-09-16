package main

import (
	"io/ioutil"
	"log"
	"os"
)

type Logger struct {
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

var L *Logger

func init() {
	L = &Logger{
		log.New(ioutil.Discard, "[DEBUG]", log.Ldate|log.Ltime),
		log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime),
		log.New(os.Stderr, "[WARN] ", log.Ldate|log.Ltime),
		log.New(os.Stderr, "[ERR]  ", log.Ldate|log.Ltime),
	}
}

func (logger *Logger) Verbose(active bool) {
	if active {
		//TODO replace in golang 1.5
		//logger.Debug.SetOutput(os.Stdout)
		logger.Debug = log.New(os.Stdout, "[DEBUG]", log.Ldate|log.Ltime)
	} else {
		logger.Debug = log.New(ioutil.Discard, "[DEBUG]", log.Ldate|log.Ltime)
	}
}
func (logger *Logger) DEBUG(v ...interface{}) {
	logger.Debug.Println(v...)
}
func (logger *Logger) ERR(v ...interface{}) {
	logger.Error.Println(v...)
}
func (logger *Logger) WARN(v ...interface{}) {
	logger.Warning.Println(v...)
}
func (logger *Logger) INFO(v ...interface{}) {
	logger.Info.Println(v...)
}
