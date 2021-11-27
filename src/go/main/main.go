package main

import (
	"ec2-test/advisor"
	awsApi "ec2-test/aws/api"
	"ec2-test/cache"
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
	cache := createCache(config.CacheConfig.Dirpath, clf.clearCache, logger)

	creds := getCredentialsForMode(clf.productionMode, config)

	regionInstancesMap, err := awsApi.GetInstances(
		&config.ApiConfig,
		config.GetRegions(),
		&creds,
		cache,
		logger,
	)
	if err != nil {
		logger.Error("Error fetching instances", zap.Error(err))
	}

	for region, instances := range regionInstancesMap {
		logger.Info("Instance count for region", zap.String("region", region.ToCodeString()), zap.Int("instanceCount", len(instances)))
	}

	advisor := advisor.NewNaiveReliabilityAdvisor()
	advisor.AdviseForRegions(regionInstancesMap, &config.Constraints)
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
		zap.Any("config", config.ToString()), // TODO: Rename to String for Go convention
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
	// TODO: Not printing parsed flags
	logger.Info("command line flags parsed", zap.Any("flags", clf))
}

func createCache(cacheFilepath string, useNewCache bool, logger *zap.Logger) *cache.Cache {
	var c *cache.Cache
	var err error

	if useNewCache {
		c, err = cache.New(cacheFilepath)
	} else {
		c, err = cache.ParseIfExistsElseNew(cacheFilepath)
	}

	if err != nil {
		err = utils.PrependToError(
			err,
			fmt.Sprintf("failed to create cache from %s", cacheFilepath),
		)
		utils.StopProgramExecution(err, 1)
	}

	return c
}
