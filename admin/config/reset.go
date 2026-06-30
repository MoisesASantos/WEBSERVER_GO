package config

import (
	"net/http"
)

func (cfg *ApiConfig) ResetHitRequest(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits.Store(0)
}
