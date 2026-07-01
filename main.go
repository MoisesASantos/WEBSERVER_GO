package main

import _ "github.com/lib/pq"
import (
	"fmt"
	"net/http"
	"os"
	"github.com/MoisesASantos/WEBSERVER_GO/admin/config"
	"github.com/MoisesASantos/WEBSERVER_GO/api/healthz"
	"github.com/MoisesASantos/WEBSERVER_GO/api/chirp"
	"github.com/MoisesASantos/WEBSERVER_GO/internal/database"
)

func main() {

	apiconfig := config.ApiConfig{}

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	apiconfig.Db := database.New(db)

	mux := http.NewServeMux()

	const filepathRoot = "."
	mux.Handle("/app/", apiconfig.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", healthz.ReadinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiconfig.MetricRequestHandler)
	mux.HandleFunc("POST /admin/reset", apiconfig.ResetHitRequest)
	mux.HandleFunc("POST /api/validate_chirp", chirp.ChirpRequestHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Printf("Error on server: %v", err)
		os.Exit(1)
	}
}
