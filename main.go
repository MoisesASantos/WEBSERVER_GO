package main

import _ "github.com/lib/pq"
import (
	"fmt"
	"net/http"
	"os"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/MoisesASantos/WEBSERVER_GO/admin/config"
	"github.com/MoisesASantos/WEBSERVER_GO/api/healthz"
	"github.com/MoisesASantos/WEBSERVER_GO/internal/database"
)

func main() {

	apiconfig := config.ApiConfig{}

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	authDev := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	apiconfig.Db = database.New(db)
	apiconfig.AuthDev = authDev

	mux := http.NewServeMux()

	const filepathRoot = "."
	mux.Handle("/app/", apiconfig.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", healthz.ReadinessHandler)
	mux.HandleFunc("GET /admin/metrics", apiconfig.MetricRequestHandler)
	mux.HandleFunc("POST /admin/reset", apiconfig.ResetHitRequest)
	mux.HandleFunc("POST /api/chirps", apiconfig.ChirpRequestHandler)
	mux.HandleFunc("POST /api/users", apiconfig.UsersRequestHandler)
	mux.HandleFunc("GET /api/chirps", apiconfig.ChirpsGetRequestHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiconfig.GetOneChirpRequestHandler)
	mux.HandleFunc("POST /api/login", apiconfig.LoginRequestHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Printf("Error on server: %v", err)
		os.Exit(1)
	}
}
