package service

import (
	"aws-blended-instances-advisor/api/schema"
	"aws-blended-instances-advisor/config"
	"aws-blended-instances-advisor/utils"
	"encoding/json"
	"io"
	"net/http"
	"sort"

	"go.uber.org/zap"
)

func getAdviseEndpointHandler(
	advise func(
		advisor schema.Advisor,
		services []schema.Service,
		options schema.Options,
	) (*schema.Advice, error),
	cfg *config.ApiConfig,
	logger *zap.Logger,
) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		reqId := utils.GenerateUuid()
		logger.Info(
			"request received",
			zap.String("url", r.Host),
			zap.String("method", r.Method),
			zap.String("requestId", reqId),
		)

		err := utils.AddCorsHeader(w, r, cfg.AllowedDomains)
		if err != nil {
			writeErrorResponse(w, reqId, err, http.StatusForbidden, logger)
			return
		}
		logger.Info("added CORS header", zap.String("requestId", reqId))

		switch r.Method {
		case "OPTIONS":
			allowedMethods := "OPTIONS, POST"
			allowedHeaders := getAllowedHeaders()
			w.Header().Set("Access-Control-Allow-Methods", allowedMethods)
			w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)

			logger.Info(
				"added headers to response",
				zap.String("requestId", reqId),
				zap.String("Access-Control-Allow-Methods", allowedMethods),
				zap.String("Access-Control-Allow-Headers", allowedHeaders),
			)

			w.WriteHeader(http.StatusOK)
			logger.Info("responded to request", zap.String("requestId", reqId))
			return

		case "POST":
			adviseEndpointPostHandler(w, r, reqId, advise, logger)
			return

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
}

func adviseEndpointPostHandler(
	w http.ResponseWriter,
	r *http.Request,
	reqId string,
	advise func(
		advisor schema.Advisor,
		services []schema.Service,
		options schema.Options,
	) (*schema.Advice, error),
	logger *zap.Logger,
) {
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

	orderServicesByDecreasingMemory(req.Services)
	logger.Info(
		"services sorted",
		zap.String("requestId", reqId),
		zap.Any("sortedServices", req.Services),
	)

	advice, err := advise(req.Advisor, req.Services, req.Options)
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

	err = req.Validate() // TODO

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

	utils.AddJsonContentTypeHeader(w)

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
