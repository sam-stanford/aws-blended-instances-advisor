package main

import (
	awsApi "ec2-test/aws/api"
	"ec2-test/config"
	"ec2-test/utils"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

func main() {
	clf := parseCommandLineFlags()

	logger, deferCallback := createLogger(clf.debugMode)
	defer deferCallback()

	logCommandLineFlags(&clf, logger)

	config := parseAndLogConfig(clf.configFilepath, logger)
	creds := getCredentialsForMode(clf.productionMode, config)

	// TODO
	os.Exit(0)

	regionInstancesMap, err := awsApi.GetInstances(
		&config.ApiConfig,
		&creds,
		config.GetRegions(),
		logger,
	)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(regionInstancesMap)
	}
}

func createLogger(debugMode bool) (logger *zap.Logger, deferCallback func() error) {
	logger, err := instantiateLogger(debugMode)
	if err != nil {
		err = utils.PrependToError(err, "failed to start logger")
		utils.StopProgramExecution(err, 1)
	}
	logger.Info("logger started")
	return logger, logger.Sync
}

func instantiateLogger(debugMode bool) (*zap.Logger, error) {
	if debugMode {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}

func parseAndLogConfig(configFilepath string, logger *zap.Logger) *config.Config {
	config, absFilepath, err := parseConfig(configFilepath)
	if err != nil {
		err = utils.PrependToError(
			err,
			fmt.Sprintf("failed to parse config from %s", absFilepath),
		)
		utils.StopProgramExecution(err, 1)
	}
	logConfig(config, absFilepath, logger)
	return config
}

func parseConfig(configFilepath string) (cfg *config.Config, absFilepath string, err error) {
	cwd, err := utils.GetCallerPath()
	if err != nil {
		err = utils.PrependToError(err, "failed to fetch current working directory")
		return
	}

	absFilepath, err = filepath.Abs(cwd + string(os.PathSeparator) + configFilepath)
	if err != nil {
		err = utils.PrependToError(err, "failed to generate config filepath")
		return
	}

	cfg, err = config.ParseConfig(absFilepath)
	if err != nil {
		err = utils.PrependToError(err, "failed to parse config")
		return
	}

	return
}

func logConfig(config *config.Config, configFilepath string, logger *zap.Logger) {
	logger.Info(
		"config parsed",
		zap.String("configFilepath", configFilepath),
		zap.String("config", config.ToString()),
	)
	// TODO: Stop escaping quotes
}

func getCredentialsForMode(isProductionMode bool, c *config.Config) config.Credentials {
	if isProductionMode {
		return c.Credentials.Production
	}
	return c.Credentials.Development
}

func logCommandLineFlags(clf *commandLineFlags, logger *zap.Logger) {
	fmt.Printf("FLAGS: %+v\n", clf)
	// TODO: Not printing parsed flags
	logger.Info("command line flags parsed", zap.Any("flags", clf))
}
