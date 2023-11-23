package handler

import (
	"go.uber.org/zap"
	"net/http"
)

type Generic struct {
	LoggerSugar *zap.SugaredLogger
}

func (h *Generic) HealthCheck(w http.ResponseWriter, r *http.Request) {
	h.LoggerSugar.Warnw("health check")
	responseReturn(w, http.StatusOK, nil)
}

func (h *Generic) NotFound(w http.ResponseWriter, r *http.Request) {
	h.LoggerSugar.Warnw("resource not found")
	responseReturn(w, http.StatusNotFound, nil)
}
