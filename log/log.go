package log

import "gopkg.in/jucardi/go-logger-lib.v1/log"

// For more information about how the logger works in the github.com/jucardi/go-logger-lib/log
// package, please refer to https://github.com/jucardi/go-logger-lib/blob/master/README.md

// LoggerMgo defines the name for the logger used for the mgo package
const LoggerMgo = "mgo"

var instance ILogger = log.Get(LoggerMgo)

// Get returns the current logger instance
func Get() ILogger {
	return instance
}

// Set assigns a new ILogger instance to be used as the logger for the mgo package
func Set(logger ILogger) {
	instance = logger
}

// Disable disables logging for the mgo package by assigning a the nil implementation of ILogger
func Disable() {
	Set(log.NewNil())
}

type ILogger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}
