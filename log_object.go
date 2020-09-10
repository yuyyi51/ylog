package ylog

import "time"

type logObject struct {
	file    string
	logTime time.Time
	level   LogLevel
	line    int
	format  string
	args    []interface{}
}

var objQueue chan *logObject

func init() {
	objQueue = make(chan *logObject, 10000)
}
