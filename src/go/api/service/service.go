package service

import (
	"aws-blended-instances-advisor/api/schema"
	"aws-blended-instances-advisor/config"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

// TODO: Add multiple end points  - one for each advisor
// TODO: 		- Allows for advisor-specific info to be passed
// TODO:		- Check provided advisor type against used path

// TODO: Test

// TODO: Calc response time & log

// TODO: Sort services

// TODO: Edit
// StartAdviseService initialises HTTP endpoints to provide an API for
// advice.
//
// Blocks the current thread, until failure, at which point the error is logged
// and the service is stopped.
func StartService(
	cfg *config.ApiConfig,
	logger *zap.Logger,
	advise func(advisor schema.Advisor, services []schema.Service, options schema.Options) (*schema.Advice, error),
) {
	http.HandleFunc("/regions", getRegionsEndpointHandler(logger))
	http.HandleFunc("/advise", getAdviseEndpointHandler(advise, logger))
	logger.Info("registered API endpoint", zap.String("path", "/advise"))

	logger.Info("starting API for advice service", zap.Int("port", cfg.Port))
	err := http.ListenAndServe(formatPort(cfg.Port), nil)

	logger.Fatal("API stopped listening to requests", zap.Error(err))
}

func formatPort(port int) string {
	return "127.0.0.1:" + strconv.Itoa(port)
}
