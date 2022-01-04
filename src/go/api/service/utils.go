package service

import (
	"net/http"
	"strings"

	"go.uber.org/zap"
)

var ALLOWED_HEADERS = [...]string{
	"Access-Control-Allow-Headers",
	"Origin", "Accept",
	"X-Requested-With",
	"Content-Type",
	"Access-Control-Request-Method",
	"Access-Control-Request-Headers",
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

func getAllowedHeaders() string {
	return strings.Join(ALLOWED_HEADERS[:], ", ")
}
