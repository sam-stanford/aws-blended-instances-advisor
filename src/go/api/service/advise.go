package service

import (
	"aws-blended-instances-advisor/api/schema"
	"aws-blended-instances-advisor/utils"
	"encoding/json"
	"io"
	"net/http"
	"sort"

	"go.uber.org/zap"
)

func getAdviseEndpointHandler(
	advise func(advisor schema.Advisor, services []schema.Service, options schema.Options) (*schema.Advice, error),
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
			utils.WriteHttpErrorResponse(w, reqId, err, http.StatusBadRequest, logger)
			return
		}
		logger.Info(
			"services parsed from request",
			zap.String("requestId", reqId),
			zap.Any("parsedServices", req.Services),
		)

		orderServicesByDecreasingMemory(req.Services)
		logger.Info(
			"services sorted",
			zap.String("requestId", reqId),
			zap.Any("sortedServices", req.Services),
		)

		advice, err := advise(req.Advisor, req.Services, req.Options)
		if err != nil {
			utils.WriteHttpErrorResponse(w, reqId, err, http.StatusInternalServerError, logger)
			return
		}
		logger.Info(
			"advice generated for request",
			zap.String("requestId", reqId),
			zap.Any("advice", advice),
		)

		err = writeAdviceResponse(w, reqId, advice, logger)
		if err != nil {
			utils.WriteHttpErrorResponse(w, reqId, err, http.StatusInternalServerError, logger)
			return
		}
	}
}

func parseRequest(r *http.Request, reqId string, logger *zap.Logger) (*schema.AdviseRequest, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, utils.PrependToError(err, "could not read request body") // TODO: Good for log, unhelpful to client
	}

	logger.Info(
		"request body read",
		zap.String("requestId", reqId),
		zap.ByteString("requestBody", body),
	)

	var req schema.AdviseRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		return nil, utils.PrependToError(err, "could not parse body JSON") // TODO: GOod for log, Unhelpful to client
	}

	err = req.Validate()

	return &req, nil
}

func orderServicesByDecreasingMemory(services []schema.Service) {
	sort.Slice(services, func(i, j int) bool {
		return services[i].MinMemory > services[j].MinMemory
	})
}

func writeAdviceResponse(
	w http.ResponseWriter,
	requestId string,
	advice *schema.Advice,
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
