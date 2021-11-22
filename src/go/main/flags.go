package main

import (
	"flag"
)

const DEFAULT_CONFIG_FILEPATH = "../../config.json"

type commandLineFlags struct {
	configFilepath string
	debugMode      bool
	productionMode bool
	clearCache     bool
}

func parseCommandLineFlags() commandLineFlags {
	configFilepath := flag.String("c", DEFAULT_CONFIG_FILEPATH, "the relative path to the config file")
	debugMode := flag.Bool("debug", false, "sets the program to debug mode")
	prodMode := flag.Bool("prod", false, "sets the program to production mode")
	clearCache := flag.Bool("clear-cache", false, "clears cached files and requests")

	flag.Parse()

	clf := commandLineFlags{
		configFilepath: *configFilepath,
		debugMode:      *debugMode,
		productionMode: *prodMode,
		clearCache:     *clearCache,
	}

	return clf
}
