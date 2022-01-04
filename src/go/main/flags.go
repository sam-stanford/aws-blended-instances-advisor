package main

import (
	"flag"
)

const DEFAULT_CONFIG_FILEPATH = "../../config.json"

type commandLineFlags struct {
	ConfigFilepath string `json:"configFilepath"`
	DebugMode      bool   `json:"debugMode"`
	ProductionMode bool   `json:"productionMode"`
	ClearCache     bool   `json:"clearCache"`
}

func parseCommandLineFlags() commandLineFlags {
	configFilepath := flag.String("c", DEFAULT_CONFIG_FILEPATH, "the relative path to the config file")
	debugMode := flag.Bool("debug", false, "sets the program to debug mode")
	prodMode := flag.Bool("prod", false, "sets the program to production mode")
	clearCache := flag.Bool("clear-cache", false, "clears cached files and requests")

	flag.Parse()

	return commandLineFlags{
		ConfigFilepath: *configFilepath,
		DebugMode:      *debugMode,
		ProductionMode: *prodMode,
		ClearCache:     *clearCache,
	}
}
