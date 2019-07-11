package main

import (
	"github.com/go-godin/log"
)

func main() {
	logger := log.NewLoggerFromEnv()
	//logger := log.NewLogger("")

	logger.Debug("test", "foo", "bar")
	logger.Info("test", "foo", "bar")
	logger.Warning("test", "foo", "bar")
	logger.Error("test", "foo", "bar")
}
