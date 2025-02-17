package main

import (
	"os"

	"github.com/skerkour/rz"
	"github.com/skerkour/rz/log"
)

func main() {

	env := os.Getenv("GO_ENV")
	hostname, _ := os.Hostname()

	// update global logger's context fields
	log.SetLogger(log.With(rz.Formatter(rz.FormatterCLI()), rz.Fields(rz.String("hostname", hostname), rz.String("environment", env))))

	log.Info("hello from logger", rz.String("hello", "world"), rz.Caller(true))

	if env == "production" {
		log.SetLogger(log.With(rz.Level(rz.InfoLevel)))
	}

	subLogger := log.With(rz.Level(rz.DebugLevel), rz.Formatter(rz.FormatterConsole()))
	SubsubLogger := rz.New(rz.Formatter(rz.FormatterConsole()))

	log.Info("info from logger", rz.String("hello", "world"), rz.Caller(true))
	// {"level":"info","hostname":"","environment":"","hello":"world","timestamp":"2019-02-07T09:30:07Z","message":"info from logger"}

	log.Error("error from logger", rz.String("hello", "world"), rz.Caller(true))

	subLogger.Debug("hello world", rz.Caller(true))
	SubsubLogger.Debug("", rz.Caller(true))

	SubsubLogger.Warn("warnin from sub logger", rz.Caller(true))
}
