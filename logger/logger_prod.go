//go:build prod 
// +build prod

package logger

import (
	"fmt"
	"log"
)

type logger struct {}

func New() Logger {
    return &logger{}
}

func (l *logger) Debug(msg string, args ...interface{}) {
	log.Println("DEBUG: " + fmt.Sprintf(msg, args...))
}

func (l *logger) Info(msg string, args ...interface{}) {
	log.Println("INFO: " + fmt.Sprintf(msg, args...))
}

func (l *logger) Warn(msg string, args ...interface{}) {
	log.Println("WARN: " + fmt.Sprintf(msg, args...))
}

func (l *logger) Error(msg string, args ...interface{}) {
	log.Println("ERROR: " + fmt.Sprintf(msg, args...))
}
