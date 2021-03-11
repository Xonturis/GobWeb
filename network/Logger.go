package network

import (
	"log"
	"os"
)

const INFO = 0
const WARNING = 1
const ERROR = 2
const NONE = 3

var logLevel = 0
var infoLogger = log.New(os.Stdout, "Cobweb INFO: ", log.Ldate|log.Ltime)
var warningLogger = log.New(os.Stdout, "Cobweb WARNING: ", log.Ldate|log.Ltime)
var errorLogger = log.New(os.Stdout, "Cobweb ERROR: ", log.Ldate|log.Ltime)

func Info(v ...interface{})  {
	if logLevel == INFO {
		infoLogger.Println(v...)
	}
}

func Warning(v ...interface{})  {
	if logLevel <= WARNING {
		warningLogger.Println(v...)
	}
}

func Error(v ...interface{})  {
	if logLevel <= ERROR {
		errorLogger.Println(v...)
	}
}

func SetLogLevel(level int)  {
	if level <= NONE && level >= INFO {
		logLevel = level
	}
}