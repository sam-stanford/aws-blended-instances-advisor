package main

import (
	"flag"

	"go.uber.org/zap"
)

const DEFAULT_CONFIG_FILEPATH = "../../config.json"

type commandLineFlags struct {
	configFilepath string
	debugMode      bool
}

func parseCommandLineFlags() commandLineFlags {
	configFilepath := flag.String("c", DEFAULT_CONFIG_FILEPATH, "the relative path to the config file")
	debugMode := flag.Bool("d", false, "sets the program to debug mode")

	flag.Parse()

	clf := commandLineFlags{
		configFilepath: *configFilepath,
		debugMode:      *debugMode,
	}

	return clf
}

func logClf(clf commandLineFlags, logger *zap.Logger) {
	logger.Info("command line flags parsed", zap.Any("flags", clf))
}
