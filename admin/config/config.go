package config

import (
	"sync/atomic"
	"net/http"
	"github.com/MoisesASantos/WEBSERVER_GO/internal/database"
)

type ApiConfig struct {
	FileserverHits 	atomic.Int32
	Db  			*database.Queries
	AuthDev 		string
	JwtKey			string
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
