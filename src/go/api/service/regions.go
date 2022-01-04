package service

import (
	"aws-blended-instances-advisor/api/schema"
	awsTypes "aws-blended-instances-advisor/aws/types"
	"aws-blended-instances-advisor/config"
	"aws-blended-instances-advisor/utils"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func getRegionsEndpointHandler(cfg *config.ApiConfig, logger *zap.Logger) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		reqId := utils.GenerateUuid()
		logger.Info(
			"request received",
			zap.String("url", r.Host),
			zap.String("requestId", reqId),
		)

		utils.AddCorsHeader(w, r, cfg.AllowedDomains)

		regions := []string{}
		for _, r := range awsTypes.GetAllRegions() {
			regions = append(regions, r.CodeString())
		}

		resp := schema.RegionsResponse{
			Regions: regions,
		}

		err := writeRegionsResponse(w, reqId, resp, logger)
		if err != nil {
			writeErrorResponse(w, reqId, err, http.StatusInternalServerError, logger)
			return
		}
	}
}

func writeRegionsResponse(
	w http.ResponseWriter,
	requestId string,
	resp schema.RegionsResponse,
	logger *zap.Logger,
) error {

	respBody, err := json.Marshal(resp)
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
