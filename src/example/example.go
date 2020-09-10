package main

import (
	"log"
	"time"

	"code.int-2.me/yuyyi51/ylog"
)

func main() {
	logger, err := ylog.NewLogger("log/", "example", "debug", 0)
	if err != nil {
		log.Fatalf("%v", err)
	}
	i := 0
	for {
		logger.Fatal("%d", i)
		i++
		time.Sleep(time.Second)
	}
}
