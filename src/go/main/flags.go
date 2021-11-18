package main

import (
	"flag"
)

const DEFAULT_CONFIG_FILEPATH = "../../config.json"

type commandLineFlags struct {
	configFilepath string
	debugMode      bool
	productionMode bool
}

func parseCommandLineFlags() commandLineFlags {
	configFilepath := flag.String("c", DEFAULT_CONFIG_FILEPATH, "the relative path to the config file")
	debugMode := flag.Bool("debug", false, "sets the program to debug mode")
	prodMode := flag.Bool("prod", false, "sets the program to production mode")

	flag.Parse()

	clf := commandLineFlags{
		configFilepath: *configFilepath,
		debugMode:      *debugMode,
		productionMode: *prodMode,
	}

	return clf
}
