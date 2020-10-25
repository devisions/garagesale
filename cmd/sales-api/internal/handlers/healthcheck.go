package handlers

import (
	"net/http"

	"github.com/devisions/garagesale/internal/platform/database"
	"github.com/devisions/garagesale/internal/platform/web"
	"github.com/jmoiron/sqlx"
)

// HealthCheck has handlers to implement the service orchestration.
type HealthCheck struct {
	DB *sqlx.DB
}

// Health responds with 200 (OK), if the service is healthy and ready for traffic.
func (hc *HealthCheck) Health(w http.ResponseWriter, r *http.Request) error {

	var health struct {
		Status string `json:"status"`
	}
	if err := database.StatusCheck(r.Context(), hc.DB); err != nil {
		health.Status = "db not ready"
		return web.Respond(w, health, http.StatusInternalServerError)
	}
	health.Status = "OK"
	return web.Respond(w, health, http.StatusOK)
}
