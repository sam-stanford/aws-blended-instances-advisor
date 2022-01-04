package service

import (
	"aws-blended-instances-advisor/api/schema"
	"aws-blended-instances-advisor/config"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

// StartAdviseService initialises HTTP endpoints for the API.
//
// Blocks the current thread, until failure, at which point the error is logged
// and the service is stopped.
func StartService(
	cfg *config.ApiConfig,
	logger *zap.Logger,
	advise func(advisor schema.Advisor, services []schema.Service, options schema.Options) (*schema.Advice, error),
) {
	http.HandleFunc("/regions", getRegionsEndpointHandler(cfg, logger))
	http.HandleFunc("/advise", getAdviseEndpointHandler(advise, cfg, logger))
	logger.Info("registered API endpoint", zap.String("path", "/advise"))

	logger.Info("starting API for advice service", zap.Int("port", cfg.Port))
	err := http.ListenAndServe(formatPort(cfg.Port), nil)

	logger.Fatal("API stopped listening to requests", zap.Error(err))
}

func formatPort(port int) string {
	return "127.0.0.1:" + strconv.Itoa(port)
}
