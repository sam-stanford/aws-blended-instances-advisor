package main

import (
	"ec2-test/aws"
	"ec2-test/config"
	"ec2-test/utils"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

const DEFAULT_CONFIG_FILEPATH = "../../config.json"

type commandLineFlags struct {
	configFilepath string
}

func main() {
	logger, deferCallback := createLogger()
	defer deferCallback()

	clf := parseCommandLineFlags(logger)
	config := parseConfig(clf, logger)

	aws.FetchInstanceInfo(config, logger)
}

func createLogger() (*zap.Logger, func() error) {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to start logger, %v\n", err)
		os.Exit(1)
	}
	logger.Info("logger started")
	return logger, logger.Sync
}

func parseConfig(commandLineFlags commandLineFlags, logger *zap.Logger) *config.Config {

	cwd, err := utils.GetCallerPath()
	if err != nil {
		logger.Fatal("failed to fetch current working directory", zap.Error(err))
	}

	configFilepath, err := filepath.Abs(cwd + string(os.PathSeparator) + commandLineFlags.configFilepath)
	if err != nil {
		logger.Fatal("failed to generate config filepath", zap.Error(err))
	}

	config, err := config.GetConfig(configFilepath)
	if err != nil {
		logger.Fatal("failed to parse config", zap.String("configFilepath", configFilepath), zap.Error(err))
	}

	logger.Info("config parsed", zap.String("configFilepath", configFilepath), zap.Any("config", config))

	return config
}

func parseCommandLineFlags(logger *zap.Logger) commandLineFlags {
	configFilepath := flag.String("c", DEFAULT_CONFIG_FILEPATH, "the relative path to the config file")
	flag.Parse()

	clf := commandLineFlags{
		configFilepath: *configFilepath,
	}

	logger.Info("command line flags parsed", zap.Any("flags", clf))

	return clf
}
