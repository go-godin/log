package main

import (
	"github.com/go-godin/log"
)

func main() {
	var logger = log.NewLoggerFromEnv()

	logger.Debug("test", "foo", "bar")
	logger.Info("test", "foo", "bar")
	logger.Warning("test", "foo", "bar")
	logger.Error("test", "foo", "bar")
}
