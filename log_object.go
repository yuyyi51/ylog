package ylog

import (
	"time"
)

type logObject struct {
	file    string
	logTime time.Time
	level   LogLevel
	line    int
	format  string
	args    []interface{}
}
