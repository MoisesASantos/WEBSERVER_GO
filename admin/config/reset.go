package config

import (
	"net/http"
)

func (cfg *ApiConfig) ResetHitRequest(w http.ResponseWriter, r *http.Request) {

	if cfg.AuthDev != "dev" {
		_ = RespondWithError(w, 403, "403 Forbidden")
		return 
	}

	cfg.Db.DeleteAllUser(r.Context())
	cfg.FileserverHits.Store(0)
}
