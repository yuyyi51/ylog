package main

import (
	"log"
	"time"

	"code.int-2.me/yuyyi51/ylog"
)

func main() {
	logger, err := ylog.NewFileLogger("log/", "example", ylog.StringToLogLevel("debug"), 0)
	if err != nil {
		log.Fatalf("%v", err)
	}
	consoleLogger, err := ylog.NewConsoleLogger(ylog.StringToLogLevel("info"), 0)
	if err != nil {
		log.Fatalf("%v", err)
	}
	logger.AddLogChain(consoleLogger)
	i := 0
	for i < 100 {
		logger.Info("simple information %d", i)
		logger.Debug("much more complex logs %d, %d", i, 2*i)
		logger.Trace("trace should not printed in debug level %d", i)
		i++
		time.Sleep(time.Millisecond * 200)
	}
}
