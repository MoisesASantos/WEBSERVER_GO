package reset

import (
	"net/http"
	"github.com/MoisesASantos/WEBSERVER_GO/admin/config"
)

func (cfg *config.ApiConfig) ResetHitRequest(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits.Store(0)
}
