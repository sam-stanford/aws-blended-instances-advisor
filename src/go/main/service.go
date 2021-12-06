package main

import (
	"ec2-test/api"
	"ec2-test/config"
	"ec2-test/utils"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

// TODO: Add multiple end points - one for each advisor
// TODO: 		- Allows for advisor-specific info to be passed
// TODO:		- Check provided advisor type against used path

// TODO: Test & maybe move to API package

// TODO: Calc response time & log

// TODO: Edit
// StartAdviseService initialises HTTP endpoints to provide an API for
// advice.
//
// Blocks the current thread, until failure, at which point the error is logged
// and the service is stopped.
func StartAdviceService(
	cfg *config.ApiConfig,
	logger *zap.Logger,
	advise func(advisor api.Advisor, services []api.Service, regions []api.Region) (*api.Advice, error),
) {
	http.HandleFunc("/advise", getAdviseEndpointHandler(advise, logger))
	logger.Info("registered API endpoint", zap.String("path", "/advise"))

	logger.Info("starting API for advice service", zap.Int("port", cfg.Port))
	err := http.ListenAndServe(formatPort(cfg.Port), nil)

	logger.Fatal("API stopped listening to requests", zap.Error(err))
}

func formatPort(port int) string {
	return "127.0.0.1:" + strconv.Itoa(port)
}

func getAdviseEndpointHandler(
	advise func(advisor api.Advisor, services []api.Service, regions []api.Region) (*api.Advice, error),
	logger *zap.Logger,
) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		reqId := utils.GenerateUuid()
		logger.Info(
			"request received",
			zap.String("url", r.Host),
			zap.String("requestId", reqId),
		)

		req, err := parseRequest(r, reqId, logger)
		if err != nil {
			writeErrorResponse(w, reqId, err, http.StatusBadRequest, logger)
			return
		}
		logger.Info(
			"services parsed from request",
			zap.String("requestId", reqId),
			zap.Any("parsedServices", req.Services),
		)

		advice, err := advise(req.Advisor, req.Services, req.Regions)
		if err != nil {
			writeErrorResponse(w, reqId, err, http.StatusInternalServerError, logger)
			return
		}
		logger.Info(
			"advice generated for request",
			zap.String("requestId", reqId),
			zap.Any("advice", advice),
		)

		err = writeAdviceResponse(w, reqId, advice, logger)
		if err != nil {
			writeErrorResponse(w, reqId, err, http.StatusInternalServerError, logger)
			return
		}
	}
}

func parseRequest(r *http.Request, reqId string, logger *zap.Logger) (*api.Request, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, utils.PrependToError(err, "could not read request body") // TODO: Good for log, unhelpful to client
	}

	logger.Info(
		"request body read",
		zap.String("requestId", reqId),
		zap.ByteString("requestBody", body),
	)

	var req api.Request
	err = json.Unmarshal(body, &req)
	if err != nil {
		return nil, utils.PrependToError(err, "could not parse body JSON") // TODO: GOod for log, Unhelpful to client
	}

	// TODO: Validation (at least one region, etc.)

	return &req, nil
}

func writeErrorResponse(
	w http.ResponseWriter,
	requestId string,
	err error,
	errCode int,
	logger *zap.Logger,
) {
	http.Error(w, err.Error(), errCode)

	logger.Error(
		"responded to request with error",
		zap.String("reqId", requestId),
		zap.Int("responseCode", errCode),
		zap.Error(err),
	)
}

func writeAdviceResponse(
	w http.ResponseWriter,
	requestId string,
	advice *api.Advice,
	logger *zap.Logger,
) error {
	respBody, err := json.Marshal(advice)
	if err != nil {
		return utils.PrependToError(err, "could not marshal advice into JSON")
	}

	_, err = w.Write(respBody)
	if err != nil {
		return utils.PrependToError(err, "could not write body of HTTP response")
	}

	logger.Info(
		"responded to request",
		zap.String("requestId", requestId),
		zap.Int("responseCode", http.StatusOK),
		zap.ByteString("response", respBody),
	)

	return nil
}
