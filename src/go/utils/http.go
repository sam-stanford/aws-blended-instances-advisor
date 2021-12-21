package utils

import (
	"net/http"

	"go.uber.org/zap"
)

// TODO: Doc & test

func WriteHttpErrorResponse(
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
